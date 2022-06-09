package hyperarp

import (
	"gitlab.com/gomidi/midi/v2"
)

type Option func(a *Arp)

func Tempo(bpm float64) Option {
	return func(a *Arp) {
		a.tempoBPM = bpm
	}
}

// NotePoolOctave sets the octave that defines the note pool, instead of the starting note
func NotePoolOctave(oct uint8) Option {
	return func(a *Arp) {
		a.notePoolOctave = oct
	}
}

// CCDirectionSwitch sets the controller for the direction switch
func CCDirectionSwitch(controller uint8) Option {
	return func(a *Arp) {
		a.directionSwitchHandler = func(msg midi.Message) (down bool, ok bool) {
			var channel, cc, value uint8
			if msg.GetControlChange(&channel, &controller, &value) {
				if cc != controller {
					return
				}

				ch := a.ControlChannel()
				if ch >= 0 && uint8(ch) != channel {
					return
				}

				return value > 0, true
			}
			return
		}
	}
}

// NoteDirectionSwitch sets the key for the direction switch
func NoteDirectionSwitch(key uint8) Option {
	return func(a *Arp) {
		a.directionSwitchHandler = func(msg midi.Message) (down bool, ok bool) {
			ch := a.ControlChannel()
			var _ch, _key, _vel uint8
			switch {
			case msg.GetNoteOn(&_ch, &_key, &_vel):
				if ch >= 0 && uint8(ch) != _ch {
					return
				}
				if _key != key {
					return
				}
				ok = true
				down = _vel > 0
			case msg.GetNoteOn(&_ch, &_key, &_vel):
				if ch >= 0 && uint8(ch) != _ch {
					return
				}
				if _key != key {
					return
				}
				ok = true
			}
			return
		}
	}
}

// CCTimeInterval sets the controller for the time interval
func CCTimeInterval(controller uint8) Option {
	return func(a *Arp) {
		a.noteDistanceHandler = func(msg midi.Message) (dist float64, ok bool) {
			dist = -1
			var _cc, _val, _ch uint8
			//cc, is := msg.(channel.ControlChange)

			if msg.GetControlChange(&_ch, &_cc, &_val) {

				if _cc != controller {
					return
				}

				ch := a.ControlChannel()
				if ch >= 0 && uint8(ch) != _ch {
					return
				}

				ok = true
				if _val > 0 {
					dist = noteDistanceMap[_val%16]
				}
			}
			return
		}
	}
}

// NoteTimeInterval sets the key for the time interval
func NoteTimeInterval(key uint8) Option {
	return func(a *Arp) {
		a.noteDistanceHandler = func(msg midi.Message) (dist float64, ok bool) {
			dist = -1
			ch := a.ControlChannel()
			var _ch, _key, _vel uint8
			switch {
			case msg.GetNoteOn(&_ch, &_key, &_vel):
				if ch >= 0 && uint8(ch) != _ch {
					return
				}
				if _key != key {
					return
				}
				ok = true
				if _vel > 0 {
					dist = noteDistanceMap[_vel%12]
				}
			case msg.GetNoteOff(&_ch, &_key, &_vel):
				if ch >= 0 && uint8(ch) != _ch {
					return
				}
				if _key != key {
					return
				}
				ok = true
			}
			return
		}
	}
}

// CCStyle sets the controller for the playing style (staccato, legato, non-legato)
func CCStyle(controller uint8) Option {
	return func(a *Arp) {
		a.styleHandler = func(msg midi.Message) (val uint8, ok bool) {
			var _ch, _cc, _val uint8

			if msg.GetControlChange(&_ch, &_cc, &_val) {

				if _cc != controller {
					return
				}

				ch := a.ControlChannel()
				if ch >= 0 && uint8(ch) != _ch {
					return
				}

				ok = true
				if _val > 0 {
					val = _val
				}
			}
			return
		}
	}
}

// NoteStyle sets the key for the playing style (staccato, legato, non-legato)
func NoteStyle(key uint8) Option {
	return func(a *Arp) {
		a.styleHandler = func(msg midi.Message) (val uint8, ok bool) {
			ch := a.ControlChannel()
			var _ch, _key, _vel uint8

			switch {
			case msg.GetNoteOn(&_ch, &_key, &_vel):
				if ch >= 0 && uint8(ch) != _ch {
					return
				}
				if _key != key {
					return
				}
				ok = true
				if _vel > 0 {
					val = _vel
				}
			case msg.GetNoteOff(&_ch, &_key, &_vel):
				if ch >= 0 && uint8(ch) != _ch {
					return
				}
				if _key != key {
					return
				}
				ok = true
			}
			return
		}
	}
}

// ControlChannel sets a separate MIDI channel for the control messages
func ControlChannel(ch uint8) Option {
	return func(a *Arp) {
		if ch < 16 {
			a.controlchannelIn = int8(ch)
		}
	}
}

// ChannelIn sets the midi channel to listen to (0-15)
func ChannelIn(ch uint8) Option {
	return func(a *Arp) {
		if ch < 16 {
			a.channelIn = int8(ch)
		}
	}
}

// Transpose sets the transposition for the midi
func Transpose(halfnotes int8) Option {
	return func(a *Arp) {
		a.transpose = halfnotes
	}
}

// ChannelOut sets the midi channel to write to
func ChannelOut(ch uint8) Option {
	return func(a *Arp) {
		if ch < 16 {
			a.channelOut = ch
		}
	}
}
