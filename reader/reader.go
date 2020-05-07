package reader

import (
	"fmt"
	"io"
	"sync"
	"time"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/midimessage/realtime"
	"gitlab.com/gomidi/midi/midireader"
	"gitlab.com/gomidi/midi/smf"
)

func ReadFrom(r *Reader, rd midi.Reader) {
	r.pos = nil
	r.reset()
	r.reader = rd
}

// ReadAll reads midi messages until an error happens
func ReadAll(r *Reader) error {
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
func ReadAllFrom(r *Reader, src io.Reader) (err error) {
	ReadFrom(r, midireader.New(src, r.dispatchRealTime, r.midiReaderOptions...))
	return ReadAll(r)
}

func (r *Reader) reset() {
	r.tempoChanges = []tempoChange{tempoChange{0, 120}}

	for c := 0; c < 16; c++ {
		r.channelRPN_NRPN[c] = [4]uint8{0, 0, 0, 0}
	}
}

func (r *Reader) saveTempoChange(pos Position, bpm float64) {
	r.tempoChanges = append(r.tempoChanges, tempoChange{pos.AbsoluteTicks, bpm})
}

// Reader allows the reading of either "over the wire" MIDI
// data (via Read) or SMF MIDI data (via ReadSMF or ReadSMFFile).
//
// Before any of the Read* methods are called, callbacks for the MIDI messages of interest
// need to be attached to the Reader. These callbacks are then invoked when the corresponding
// MIDI message arrives. They must not be changed while any of the Read* methods is running.
//
// It is possible to share the same Reader for reading of the wire MIDI ("live")
// and SMF Midi data as long as not more than one Read* method is running at a point in time.
// However, only channel messages and system exclusive message may be used in both cases.
// To enable this, the corresponding callbacks receive a pointer to the Position of the
// MIDI message. This pointer is always nil for "live" MIDI data and never nil when
// reading from a SMF.
//
// The SMF header callback and the meta message callbacks are only called, when reading data
// from an SMF. Therefore the Position is passed directly and can't be nil.
//
// System common and realtime message callbacks will only be called when reading "live" MIDI,
// so they get no Position.
type Reader struct {
	tempoChanges      []tempoChange       // track tempo changes
	header            smf.Header          // store the SMF header
	logger            Logger              // optional logger
	pos               *Position           // the current SMFPosition
	errSMF            error               // error when reading SMF
	midiReaderOptions []midireader.Option // options for the midireader
	reader            midi.Reader
	midiClocks        [3]*time.Time
	clockmx           sync.Mutex // protect the midiClocks
	ignoreMIDIClock   bool

	channelRPN_NRPN [16][4]uint8 // channel -> [cc0,cc1,valcc0,valcc1], initial value [-1,-1,-1,-1]

	// ticks per quarternote
	resolution smf.MetricTicks

	// SMFHeader is the callback that gets SMF header data
	smfheader func(smf.Header)

	message
}

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

func (r *Reader) dispatchRealTime(m realtime.Message) {

	// ticks (most important, must be sent every 10 milliseconds) comes first
	if m == realtime.Tick {
		if r.message.Realtime.Tick != nil {
			r.message.Realtime.Tick()
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

		if r.message.Realtime.Clock != nil {
			r.message.Realtime.Clock()
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
		if r.message.Meta.TempoBPM != nil {
			r.message.Meta.TempoBPM(*r.pos, bpm)
		}

		return
	}

	// starting should not take too long
	if m == realtime.Start {
		if r.message.Realtime.Start != nil {
			r.message.Realtime.Start()
		}
		return
	}

	// continuing should not take too long
	if m == realtime.Continue {
		if r.message.Realtime.Continue != nil {
			r.message.Realtime.Continue()
		}
		return
	}

	// Active Sense must come every 300 milliseconds
	// (but is seldom implemented)
	if m == realtime.Activesense {
		if r.message.Realtime.Activesense != nil {
			r.message.Realtime.Activesense()
		}
		return
	}

	// put any user defined realtime message here
	if m == realtime.Undefined4 {
		if r.message.Unknown != nil {
			r.message.Unknown(r.pos, m)
		}
		return
	}

	// stopping is not so urgent
	if m == realtime.Stop {
		if r.message.Realtime.Stop != nil {
			r.message.Realtime.Stop()
		}
		return
	}

	// reset may take some time anyway
	if m == realtime.Reset {
		if r.message.Realtime.Reset != nil {
			r.message.Realtime.Reset()
		}
		return
	}
}

// New returns a new Reader
func New(callbacksAndOptions ...func(r *Reader)) *Reader {
	r := &Reader{logger: logfunc(printf)}

	for _, c := range callbacksAndOptions {
		if c != nil {
			c(r)
		}
	}

	return r
}

func (r *Reader) Callback(cb func(*Reader)) {
	cb(r)
}

// Position is the position of the event inside a standard midi file (SMF) or since
// start listening on a connection.
type Position struct {

	// Track is the number of the track, starting with 0
	Track int16

	// DeltaTicks is number of ticks that passed since the previous message in the same track
	DeltaTicks uint32

	// AbsoluteTicks is the number of ticks that passed since the beginning of the track
	AbsoluteTicks uint64
}

// log does the logging
func (r *Reader) log(m midi.Message) {
	if r.pos != nil {
		r.logger.Printf("#%v [%v d:%v] %s\n", r.pos.Track, r.pos.AbsoluteTicks, r.pos.DeltaTicks, m)
	} else {
		r.logger.Printf("%s\n", m)
	}
}

// dispatch dispatches the messages from the midi.Reader (which might be an smf reader)
// for realtime reading, the passed *SMFPosition is nil
func (r *Reader) dispatch() (err error) {
	for {
		err = r.dispatchMessageFromReader()
		if err != nil {
			return
		}
	}
}

func (r *Reader) _RPN_NRPN_Reset(ch uint8, isRPN bool) {
	// reset tracking on this channel
	r.channelRPN_NRPN[ch] = [4]uint8{0, 0, 0, 0}

	if isRPN {
		if r.message.Channel.ControlChange.RPN.Reset != nil {
			r.message.Channel.ControlChange.RPN.Reset(r.pos, ch)
			return
		}
		if r.message.Channel.ControlChange.RPN.MSB != nil {
			r.message.Channel.ControlChange.RPN.MSB(r.pos, ch, 127, 127, 0)
		}

		return
	}

	if r.message.Channel.ControlChange.NRPN.Reset != nil {
		r.message.Channel.ControlChange.NRPN.Reset(r.pos, ch)
		return
	}
	if r.message.Channel.ControlChange.NRPN.MSB != nil {
		r.message.Channel.ControlChange.NRPN.MSB(r.pos, ch, 127, 127, 0)
	}

}

func (r *Reader) sendAsCC(ch, cc, val uint8) error {
	if r.message.Channel.ControlChange.Each != nil {
		r.message.Channel.ControlChange.Each(r.pos, ch, cc, val)
	}
	return nil
}

func (r *Reader) hasRPNCallback() bool {
	return !(r.message.Channel.ControlChange.RPN.MSB == nil && r.message.Channel.ControlChange.RPN.LSB == nil)
}

func (r *Reader) hasNRPNCallback() bool {
	return !(r.message.Channel.ControlChange.NRPN.MSB == nil && r.message.Channel.ControlChange.NRPN.LSB == nil)
}

func (r *Reader) hasNoRPNorNRPNCallback() bool {
	return !r.hasRPNCallback() && !r.hasNRPNCallback()
}

// TimeAt returns the time.Duration at the given absolute position counted
// from the beginning of the file, respecting all the tempo changes in between.
// If the time format is not of type smf.MetricTicks, nil is returned.
func TimeAt(r *Reader, absTicks uint64) *time.Duration {
	if r.resolution == 0 {
		return nil
	}

	var tc = tempoChange{0, 120}
	var lastTick uint64
	var lastDur time.Duration
	for _, t := range r.tempoChanges {
		if t.absTicks >= absTicks {
			// println("stopping")
			break
		}
		// println("pre", "lastDur", lastDur, "lastTick", lastTick, "bpm", tc.bpm)
		lastDur += calcDeltaTime(r.resolution, uint32(t.absTicks-lastTick), tc.bpm)
		tc = t
		lastTick = t.absTicks
	}
	result := lastDur + calcDeltaTime(r.resolution, uint32(absTicks-lastTick), tc.bpm)
	return &result
}

// dispatchMessageFromReader dispatches a single message from the midi.Reader (which might be an smf reader)
// for realtime reading, the passed *SMFPosition is nil
func (r *Reader) dispatchMessageFromReader() (err error) {
	var m midi.Message
	m, err = r.reader.Read()
	if err != nil {
		return
	}

	return r.dispatchMessage(m)
}

type tempoChange struct {
	absTicks uint64
	bpm      float64
}

func calcDeltaTime(mt smf.MetricTicks, deltaTicks uint32, bpm float64) time.Duration {
	return mt.FractionalDuration(bpm, deltaTicks)
}
