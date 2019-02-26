package mid

import (
	"sync"
	"time"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/midireader"
	"gitlab.com/gomidi/midi/smf"
)

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
	SMFHeader func(smf.Header)

	// Msg provides callbacks for MIDI messages
	Msg struct {

		// Each is called for every MIDI message in addition to the other callbacks.
		Each func(*Position, midi.Message)

		// Unknown is called for undefined or unknown messages
		Unknown func(p *Position, msg midi.Message)

		// Meta provides callbacks for meta messages (only in SMF files)
		Meta struct {

			// Copyright is called for the copyright message
			Copyright func(p Position, text string)

			// TempoBPM is called for the tempo (change) message, BPM is fractional
			TempoBPM func(p Position, bpm float64)

			// TimeSigis called for the time signature (change) message
			TimeSig func(p Position, num, denom uint8)

			// Key is called for the key signature (change) message
			Key func(p Position, key uint8, ismajor bool, num_accidentals uint8, accidentals_are_flat bool)

			// Instrument is called for the instrument (name) message
			Instrument func(p Position, name string)

			// TrackSequenceName is called for the sequence / track name message
			// If in a format 0 track, or the first track in a format 1 file, the name of the sequence. Otherwise, the name of the track.
			TrackSequenceName func(p Position, name string)

			// SequenceNo is called for the sequence number message
			SequenceNo func(p Position, number uint16)

			// Marker is called for the marker message
			Marker func(p Position, text string)

			// Cuepoint is called for the cuepoint message
			Cuepoint func(p Position, text string)

			// Text is called for the text message
			Text func(p Position, text string)

			// Lyric is called for the lyric message
			Lyric func(p Position, text string)

			// EndOfTrack is called for the end of a track message
			EndOfTrack func(p Position)

			// Device is called for the device port message
			Device func(p Position, name string)

			// Program is called for the program name message
			Program func(p Position, text string)

			// SMPTE is called for the smpte offset message
			SMPTE func(p Position, hour, minute, second, frame, fractionalFrame byte)

			// SequencerData is called for the sequencer specific message
			SequencerData func(p Position, data []byte)

			Deprecated struct {
				// Channel is called for the deprecated MIDI channel message
				Channel func(p Position, channel uint8)

				// Port is called for the deprecated MIDI port message
				Port func(p Position, port uint8)
			}
		}

		// Channel provides callbacks for channel messages
		// They may occur in SMF files and in live MIDI.
		// For live MIDI *Position is nil.
		Channel struct {

			// NoteOn is just called for noteon messages with a velocity > 0.
			// Noteon messages with velocity == 0 will trigger NoteOff with a velocity of 0.
			NoteOn func(p *Position, channel, key, velocity uint8)

			// NoteOff is called for noteoff messages (then the given velocity is passed)
			// and for noteon messages of velocity 0 (then velocity is 0).
			NoteOff func(p *Position, channel, key, velocity uint8)

			// Pitchbend is called for pitch bend messages
			Pitchbend func(p *Position, channel uint8, value int16)

			// ProgramChange is called for program change messages. Program numbers start with 0.
			ProgramChange func(p *Position, channel, program uint8)

			// Aftertouch is called for aftertouch messages  (aka "channel pressure")
			Aftertouch func(p *Position, channel, pressure uint8)

			// PolyAftertouch is called for polyphonic aftertouch messages (aka "key pressure").
			PolyAftertouch func(p *Position, channel, key, pressure uint8)

			// ControlChange deals with control change messages
			ControlChange struct {

				// Each is called for every control change message
				// If RPN or NRPN callbacks are defined, the corresponding control change messages will not
				// be passed to each and the corrsponding RPN/NRPN callback are called.
				Each func(p *Position, channel, controller, value uint8)

				// RPN deals with Registered Program Numbers (RPN) and their values.
				// If the callbacks are set, the corresponding control change messages will not be passed of ControlChange.Each.
				RPN struct {

					// MSB is called, when the MSB of a RPN arrives
					MSB func(p *Position, channel, typ1, typ2, msbVal uint8)

					// LSB is called, when the MSB of a RPN arrives
					LSB func(p *Position, channel, typ1, typ2, lsbVal uint8)

					// Increment is called, when the increment of a RPN arrives
					Increment func(p *Position, channel, typ1, typ2 uint8)

					// Decrement is called, when the decrement of a RPN arrives
					Decrement func(p *Position, channel, typ1, typ2 uint8)

					// Reset is called, when the reset or null RPN arrives
					Reset func(p *Position, channel uint8)
				}

				// NRPN deals with Non-Registered Program Numbers (NRPN) and their values.
				// If the callbacks are set, the corresponding control change messages will not be passed of ControlChange.Each.
				NRPN struct {

					// MSB is called, when the MSB of a NRPN arrives
					MSB func(p *Position, channel uint8, typ1, typ2, msbVal uint8)

					// LSB is called, when the LSB of a NRPN arrives
					LSB func(p *Position, channel uint8, typ1, typ2, msbVal uint8)

					// Increment is called, when the increment of a NRPN arrives
					Increment func(p *Position, channel, typ1, typ2 uint8)

					// Decrement is called, when the decrement of a NRPN arrives
					Decrement func(p *Position, channel, typ1, typ2 uint8)

					// Reset is called, when the reset or null NRPN arrives
					Reset func(p *Position, channel uint8)
				}
			}
		}

		// Realtime provides callbacks for realtime messages.
		// They are only used with "live" MIDI
		Realtime struct {

			// Clock is called for a clock message
			Clock func()

			// Tick is called for a tick message
			Tick func()

			// Activesense is called for a active sense message
			Activesense func()

			// Start is called for a start message
			Start func()

			// Stop is called for a stop message
			Stop func()

			// Continue is called for a continue message
			Continue func()

			// Reset is called for a reset message
			Reset func()
		}

		// SysCommon provides callbacks for system common messages.
		// They are only used with "live" MIDI
		SysCommon struct {

			// Tune is called for a tune request message
			Tune func()

			// SongSelect is called for a song select message
			SongSelect func(num uint8)

			// SPP is called for a song position pointer message
			SPP func(pos uint16)

			// MTC is called for a MIDI timing code message
			MTC func(frame uint8)
		}

		// SysEx provides callbacks for system exclusive messages.
		// They may occur in SMF files and in live MIDI.
		// For live MIDI *Position is nil.
		SysEx struct {

			// Complete is called for a complete system exclusive message
			Complete func(p *Position, data []byte)

			// Start is called for a starting system exclusive message
			Start func(p *Position, data []byte)

			// Continue is called for a continuing system exclusive message
			Continue func(p *Position, data []byte)

			// End is called for an ending system exclusive message
			End func(p *Position, data []byte)

			// Escape is called for an escaping system exclusive message
			Escape func(p *Position, data []byte)
		}
	}
}

// NewReader returns a new Reader
func NewReader(opts ...ReaderOption) *Reader {
	h := &Reader{logger: logfunc(printf)}

	for _, opt := range opts {
		opt(h)
	}

	return h
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
