package mid

import (
	"fmt"
	"io"

	"time"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/midimessage/realtime"
	"gitlab.com/gomidi/midi/midireader"
	"gitlab.com/gomidi/midi/smf"
)

var _ midi.Reader = &Reader{}
var _ smf.Reader = &Reader{}

func (r *Reader) Track() int16 {
	return r.pos.Track
}

func (r *Reader) Delta() uint32 {
	return r.pos.DeltaTicks
}

func (r *Reader) Header() smf.Header {
	return r.header
}

func (r *Reader) ReadHeader() error {
	rd, ok := r.reader.(smf.Reader)
	if !ok {
		return fmt.Errorf("header could only be read from SMF files")
	}
	err := rd.ReadHeader()
	if err != nil {
		return err
	}
	r.setHeader(rd.Header())
	return nil
}

// Read reads a midi.Message and dispatches it.
func (r *Reader) Read() (m midi.Message, err error) {
	m, err = r.reader.Read()
	if err != nil {
		return nil, err
	}

	err = r.dispatchMessage(m)
	return
}

func (r *Reader) ReadFrom(rd midi.Reader) {
	r.pos = nil
	r.reset()
	r.reader = rd
}

// ReadAll reads midi messages until an error happens
func (r *Reader) ReadAll() error {
	return r.dispatch()
}

// ReadAllFrom reads midi messages from src until an error happens (for "live" MIDI data "over the wire").
// io.EOF is the expected error that is returned when reading should stop.
//
// ReadAllFrom does not close the src.
//
// The messages are dispatched to the corresponding attached functions of the Reader.
//
// They must be attached before Reader.ReadAllFrom is called
// and they must not be unset or replaced until Read returns.
// For more infomation about dealing with the MIDI messages, see Reader.
func (r *Reader) ReadAllFrom(src io.Reader) (err error) {
	r.ReadFrom(midireader.New(src, r.dispatchRealTime, r.midiReaderOptions...))
	return r.ReadAll()
}

func (r *Reader) dispatchRealTime(m realtime.Message) {

	// ticks (most important, must be sent every 10 milliseconds) comes first
	if m == realtime.Tick {
		if r.Msg.Realtime.Tick != nil {
			r.Msg.Realtime.Tick()
		}
		return
	}

	// clock a bit slower synchronization method (24 MIDI Clocks in every quarter note) comes next
	// we can use this to calculate the tempo.
	if m == realtime.TimingClock {
		var gotClock time.Time
		if !r.ignoreMIDIClock {
			gotClock = time.Now()
		}

		if r.Msg.Realtime.Clock != nil {
			r.Msg.Realtime.Clock()
		}

		if r.ignoreMIDIClock {
			return
		}

		r.clockmx.Lock()

		if r.midiClocks[0] == nil {
			r.midiClocks[0] = &gotClock
			return
		}

		if r.midiClocks[1] == nil {
			r.midiClocks[1] = &gotClock
			return
		}

		if r.midiClocks[2] == nil {
			r.midiClocks[2] = &gotClock
			return
		}

		bpm := tempoBasedOnMIDIClocks(r.midiClocks[0], r.midiClocks[1], r.midiClocks[2], &gotClock)

		// move them over
		r.midiClocks[0] = r.midiClocks[1]
		r.midiClocks[1] = r.midiClocks[2]
		r.midiClocks[2] = &gotClock

		r.clockmx.Unlock()

		r.saveTempoChange(*r.pos, bpm)
		if r.Msg.Meta.TempoBPM != nil {
			r.Msg.Meta.TempoBPM(*r.pos, bpm)
		}

		return
	}

	// starting should not take too long
	if m == realtime.Start {
		if r.Msg.Realtime.Start != nil {
			r.Msg.Realtime.Start()
		}
		return
	}

	// continuing should not take too long
	if m == realtime.Continue {
		if r.Msg.Realtime.Continue != nil {
			r.Msg.Realtime.Continue()
		}
		return
	}

	// Active Sense must come every 300 milliseconds
	// (but is seldom implemented)
	if m == realtime.Activesense {
		if r.Msg.Realtime.Activesense != nil {
			r.Msg.Realtime.Activesense()
		}
		return
	}

	// put any user defined realtime message here
	if m == realtime.Undefined4 {
		if r.Msg.Unknown != nil {
			r.Msg.Unknown(r.pos, m)
		}
		return
	}

	// stopping is not so urgent
	if m == realtime.Stop {
		if r.Msg.Realtime.Stop != nil {
			r.Msg.Realtime.Stop()
		}
		return
	}

	// reset may take some time anyway
	if m == realtime.Reset {
		if r.Msg.Realtime.Reset != nil {
			r.Msg.Realtime.Reset()
		}
		return
	}
}
