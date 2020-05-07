package reader

import (
	"fmt"
	"io"
	"os"

	"gitlab.com/gomidi/midi/smf"
	"gitlab.com/gomidi/midi/smf/smfreader"
)

// ReadSMFFile open, reads and closes a complete SMF file.
// If the read content was a valid midi file, nil is returned.
//
// The messages are dispatched to the corresponding attached functions of the handler.
//
// They must be attached before Reader.ReadSMF is called
// and they must not be unset or replaced until ReadSMF returns.
// For more infomation about dealing with the SMF midi messages, see Reader and
// SMFPosition.
func ReadSMFFile(r *Reader, file string, options ...smfreader.Option) error {
	r.errSMF = nil
	r.pos = &Position{}
	r.reset()
	err := smfreader.ReadFile(file, r.readSMF2, options...)
	if err != nil && err != smf.ErrFinished {
		return err
	}
	if r.errSMF == smf.ErrFinished {
		return nil
	}
	return nil
}

// ReadSMFFileHeader reads just the header of a SMF file
func ReadSMFFileHeader(r *Reader, file string, options ...smfreader.Option) (smf.Header, error) {
	r.errSMF = nil
	r.pos = &Position{}
	r.reset()
	f, err := os.Open(file)
	if err != nil {
		fmt.Printf("can't open file: %v", err)
		return smf.Header{}, err
	}
	r.reader = smfreader.New(f, options...)

	err2 := r.ReadHeader()

	f.Close()

	if err2 != nil && err2 != smf.ErrFinished {
		return r.Header(), err2
	}
	return r.Header(), nil
}

func ReadSMFFrom(r *Reader, rd smf.Reader) error {
	r.errSMF = nil
	r.pos = &Position{}
	r.reset()
	r.reader = rd

	err := r.ReadHeader()
	if err != nil {
		return err
	}
	r.readSMF()

	if r.errSMF == smf.ErrFinished {
		return nil
	}
	return r.errSMF
}

// ReadAllSMF reads midi messages from src (which is supposed to be the content of a standard midi file (SMF))
// until an error or io.EOF happens.
//
// ReadAllSMF does not close the src.
//
// If the read content was a valid midi file, nil is returned.
//
// The messages are dispatched to the corresponding attached functions of the Reader.
//
// They must be attached before Reader.ReadSMF is called
// and they must not be unset or replaced until ReadSMF returns.
// For more infomation about dealing with the SMF midi messages, see Reader and
// SMFPosition.
func ReadSMF(r *Reader, src io.Reader, options ...smfreader.Option) error {
	return ReadSMFFrom(r, smfreader.New(src, options...))
}

func (r *Reader) setHeader(hd smf.Header) {
	r.header = hd

	if metric, isMetric := r.header.TimeFormat.(smf.MetricTicks); isMetric {
		r.resolution = metric
	}

	if r.smfheader != nil {
		r.smfheader(r.header)
	}
}

func (r *Reader) readSMF2(rd smf.Reader) {
	r.reader = rd
	err := r.ReadHeader()
	if err != nil {
		r.errSMF = err
	}
	r.readSMF()
}

func (r *Reader) readSMF() {
	err := r.dispatch()
	if err != io.EOF {
		r.errSMF = err
	}
}

/*
These are differences  some differences between the two:
  - SMF events have a message and a delta time while live MIDI just has messages
  - a SMF file has a header
  - a SMF file may have meta messages (e.g. track name, copyright etc) while live MIDI must not have them
  - live MIDI may have realtime and syscommon messages while a SMF must not have them

In sum the only MIDI message handling that could be shared in practice is the handling of channel
and sysex messages.

That is reflected in the type signature of the callbacks and results in three kinds of callbacks:
  - The ones that don't receive a SMFPosition are called for messages that may only appear live
  - The ones with that receive a SMFPosition are called for messages that may only appear
    within a SMF file
  - The onse with that receive a *SMFPosition are called for messages that may appear live
    and within a SMF file. In a SMF file the pointer is never nil while in a live situation it always is.

Due to the nature of a SMF file, the tracks have no "assigned numbers", and can in the worst case just be
distinguished by the order in which they appear inside the file.

In addition the delta times are just relative to the previous message inside the track.

This has some severe consequences:

  - Eeven if the user is just interested in some kind of messages (e.g. note on and note off),
    he still has to deal with the delta times of every message.
  - At the beginning of each track he has to reset the tracking time.

This procedure is error prone and therefor SMFPosition provides a helper that contains not just
the original delta (which shouldn't be needed most of the time), but also the absolute position (in ticks)
and order number of the track in which the message appeared.

This way the user can ignore the messages and tracks he is not interested in.

See the example for a handler that handles both live and SMF messages.

For writing there is a LiveWriter for "over the wire" writing and a SMFWriter to write SMF files.

*/
