package smf

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/gm"
)

var smfJSONKeys = struct {
	format,
	timeformat,
	metricticks,
	timecode,
	framespersecond,
	tracks,
	delta,
	typ,
	data,
	channel,
	pressure,
	controller,
	value,
	key,
	velocity,
	relative,
	absolute,
	program,
	num,
	isMajor,
	isFlat,
	hour,
	minute,
	second,
	frame,
	fractframe,
	denom,
	clockspertick,
	demisemiquaverperquarter,
	controllername,
	keyname,
	programname,
	subframes string

	aftertouchType,
	controlchangeType,
	noteonType,
	noteoffType,
	pitchbendType,
	polyaftertouchType,
	programchangeType,
	sysexType,
	channelType,
	copyrightType,
	cuepointType,
	deviceType,
	instrumentType,
	keysignatureType,
	lyricType,
	markerType,
	portType,
	programnameType,
	smpteoffsetType,
	seqdataType,
	seqnumberType,
	tempoType,
	textType,
	timesignatureType,
	tracknameType string
}{
	format:                   "format",
	timeformat:               "timeformat",
	metricticks:              "metricticks",
	timecode:                 "timecode",
	framespersecond:          "framespersecond",
	tracks:                   "tracks",
	delta:                    "delta",
	typ:                      "type",
	data:                     "data",
	channel:                  "channel",
	pressure:                 "pressure",
	controller:               "controller",
	value:                    "value",
	key:                      "key",
	velocity:                 "velocity",
	relative:                 "relative",
	absolute:                 "absolute",
	program:                  "program",
	num:                      "num",
	isMajor:                  "isMajor",
	isFlat:                   "isFlat",
	hour:                     "hour",
	minute:                   "minute",
	second:                   "second",
	frame:                    "frame",
	fractframe:               "fractframe",
	denom:                    "denom",
	clockspertick:            "clockspertick",
	demisemiquaverperquarter: "demisemiquaverperquarter",
	controllername:           "controllername",
	keyname:                  "keyname",
	programname:              "programname",
	subframes:                "subframes",

	aftertouchType:     "aftertouch",
	controlchangeType:  "controlchange",
	noteonType:         "noteon",
	noteoffType:        "noteoff",
	pitchbendType:      "pitchbend",
	polyaftertouchType: "polyaftertouch",
	programchangeType:  "programchange",
	sysexType:          "sysex",
	channelType:        "channel",
	copyrightType:      "copyright",
	cuepointType:       "cuepoint",
	deviceType:         "device",
	instrumentType:     "instrument",
	keysignatureType:   "keysignature",
	lyricType:          "lyric",
	markerType:         "marker",
	portType:           "port",
	programnameType:    "programname",
	smpteoffsetType:    "smpteoffset",
	seqdataType:        "seqdata",
	seqnumberType:      "seqnumber",
	tempoType:          "tempo",
	textType:           "text",
	timesignatureType:  "timesignature",
	tracknameType:      "trackname",
}

func (s *SMF) UnmarshalJSON(data []byte) error {

	var all = map[string]any{}

	err := json.Unmarshal(data, &all)

	if err != nil {
		return err
	}

	s.Tracks = nil
	s.format = uint16(all[smfJSONKeys.format].(float64))

	if tf, has := all[smfJSONKeys.timeformat]; has {
		ttf := tf.(map[string]any)
		if mc, has := ttf[smfJSONKeys.metricticks]; has {
			s.TimeFormat = MetricTicks(uint16(mc.(float64)))
		}

		if tc, has := ttf[smfJSONKeys.timecode]; has {
			var t TimeCode
			ttc := tc.(map[string]any)
			t.FramesPerSecond = uint8(ttc[smfJSONKeys.framespersecond].(float64))
			t.SubFrames = uint8(ttc[smfJSONKeys.subframes].(float64))
			s.TimeFormat = t
		}
	} else {
		s.TimeFormat = defaultMetric
	}

	tracks := all[smfJSONKeys.tracks].([]any) // ([][]map[string]any)

	for _, _track := range tracks {
		track := _track.([]any)

		var t Track

		var deltaoffset uint32

		for _, _ev := range track {
			ev := _ev.(map[string]any)
			var e Event
			e.Delta = uint32(ev[smfJSONKeys.delta].(float64)) + deltaoffset

			switch ev[smfJSONKeys.typ] {

			case smfJSONKeys.aftertouchType:
				d := ev[smfJSONKeys.data].(map[string]any)
				channel := d[smfJSONKeys.channel].(float64)
				pressure := d[smfJSONKeys.pressure].(float64)
				e.Message = midi.AfterTouch(uint8(channel), uint8(pressure)).Bytes()

			case smfJSONKeys.controlchangeType:
				d := ev[smfJSONKeys.data].(map[string]any)
				channel := d[smfJSONKeys.channel].(float64)
				controller := d[smfJSONKeys.controller].(float64)
				value := d[smfJSONKeys.value].(float64)
				e.Message = midi.ControlChange(uint8(channel), uint8(controller), uint8(value)).Bytes()

			case smfJSONKeys.noteonType:
				d := ev[smfJSONKeys.data].(map[string]any)
				channel := d[smfJSONKeys.channel].(float64)
				key := d[smfJSONKeys.key].(float64)
				velocity := d[smfJSONKeys.velocity].(float64)
				e.Message = midi.NoteOn(uint8(channel), uint8(key), uint8(velocity)).Bytes()

			case smfJSONKeys.noteoffType:
				d := ev[smfJSONKeys.data].(map[string]any)
				channel := d[smfJSONKeys.channel].(float64)
				key := d[smfJSONKeys.key].(float64)
				// velocity := d[smfJSONKeys.velocity].(float64)
				e.Message = midi.NoteOff(uint8(channel), uint8(key)).Bytes()

			case smfJSONKeys.pitchbendType:
				d := ev[smfJSONKeys.data].(map[string]any)
				channel := d[smfJSONKeys.channel].(float64)
				relative := d[smfJSONKeys.relative].(float64)
				e.Message = midi.Pitchbend(uint8(channel), int16(relative)).Bytes()

			case smfJSONKeys.polyaftertouchType:
				d := ev[smfJSONKeys.data].(map[string]any)
				channel := d[smfJSONKeys.channel].(float64)
				key := d[smfJSONKeys.key].(float64)
				pressure := d[smfJSONKeys.pressure].(float64)
				e.Message = midi.PolyAfterTouch(uint8(channel), uint8(key), uint8(pressure)).Bytes()

			case smfJSONKeys.programchangeType:
				d := ev[smfJSONKeys.data].(map[string]any)
				channel := d[smfJSONKeys.channel].(float64)
				program := d[smfJSONKeys.program].(float64)
				e.Message = midi.ProgramChange(uint8(channel), uint8(program)).Bytes()

			case smfJSONKeys.sysexType:
				bt, err := hex.DecodeString(ev[smfJSONKeys.data].(string))
				if err != nil {
					return fmt.Errorf("can't decode sysex %q", ev[smfJSONKeys.data].(string))
				}
				e.Message = midi.SysEx(bt).Bytes()

			case smfJSONKeys.channelType:
				e.Message = MetaChannel(uint8(ev[smfJSONKeys.data].(float64)))

			case smfJSONKeys.copyrightType:
				e.Message = MetaCopyright(ev[smfJSONKeys.data].(string))

			case smfJSONKeys.cuepointType:
				e.Message = MetaCuepoint(ev[smfJSONKeys.data].(string))

			case smfJSONKeys.deviceType:
				e.Message = MetaDevice(ev[smfJSONKeys.data].(string))

			case smfJSONKeys.instrumentType:
				e.Message = MetaInstrument(ev[smfJSONKeys.data].(string))

			case smfJSONKeys.keysignatureType:
				d := ev[smfJSONKeys.data].(map[string]any)
				key := d[smfJSONKeys.key].(float64)
				num := d[smfJSONKeys.num].(float64)
				isMajor := d[smfJSONKeys.isMajor].(bool)
				isFlat := d[smfJSONKeys.isFlat].(bool)
				e.Message = MetaKey(uint8(key), isMajor, uint8(num), isFlat)

			case smfJSONKeys.lyricType:
				e.Message = MetaLyric(ev[smfJSONKeys.data].(string))

			case smfJSONKeys.markerType:
				e.Message = MetaMarker(ev[smfJSONKeys.data].(string))

			case smfJSONKeys.portType:
				e.Message = MetaPort(uint8(ev[smfJSONKeys.data].(float64)))

			case smfJSONKeys.programnameType:
				e.Message = MetaProgram(ev[smfJSONKeys.data].(string))

			case smfJSONKeys.smpteoffsetType:
				d := ev[smfJSONKeys.data].(map[string]any)
				hour := d[smfJSONKeys.hour].(float64)
				minute := d[smfJSONKeys.minute].(float64)
				second := d[smfJSONKeys.second].(float64)
				frame := d[smfJSONKeys.frame].(float64)
				fractframe := d[smfJSONKeys.fractframe].(float64)
				e.Message = MetaSMPTE(uint8(hour), uint8(minute), uint8(second), uint8(frame), uint8(fractframe)).Bytes()

			case smfJSONKeys.seqdataType:
				bt, err := hex.DecodeString(ev[smfJSONKeys.data].(string))
				if err != nil {
					return fmt.Errorf("can't decode seqdata %q", ev[smfJSONKeys.data].(string))
				}
				e.Message = MetaSequencerData(bt).Bytes()

			case smfJSONKeys.seqnumberType:
				e.Message = MetaSequenceNo(uint16(ev[smfJSONKeys.data].(float64)))

			case smfJSONKeys.tempoType:
				e.Message = MetaTempo(ev[smfJSONKeys.data].(float64))

			case smfJSONKeys.textType:
				e.Message = MetaText(ev[smfJSONKeys.data].(string))

			case smfJSONKeys.timesignatureType:
				d := ev[smfJSONKeys.data].(map[string]any)
				num := d[smfJSONKeys.num].(float64)
				denom := d[smfJSONKeys.denom].(float64)
				clockspertick := d[smfJSONKeys.clockspertick].(float64)
				demisemiquaverperquarter := d[smfJSONKeys.demisemiquaverperquarter].(float64)
				e.Message = MetaTimeSig(uint8(num), uint8(denom), uint8(clockspertick), uint8(demisemiquaverperquarter))

			case smfJSONKeys.tracknameType:
				e.Message = MetaTrackSequenceName(ev[smfJSONKeys.data].(string))

			default:
				deltaoffset = e.Delta + deltaoffset

				// ignore the (invalid) message but not the delta movement, to keep everything in place
				continue
			}

			t = append(t, e)
		}

		s.Tracks = append(s.Tracks, t)
	}

	return nil
}

func (s *SMF) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.serializeMap())
}

func (s *SMF) MarshalJSONIndent() ([]byte, error) {
	return json.MarshalIndent(s.serializeMap(), "", "  ")
}

func (s *SMF) serializeMap() map[string]any {

	var all = map[string]any{}

	all[smfJSONKeys.format] = s.format

	timeformat := map[string]any{}

	if mt, ok := s.TimeFormat.(MetricTicks); ok {
		timeformat[smfJSONKeys.metricticks] = mt.Resolution()
	}

	if tc, ok := s.TimeFormat.(TimeCode); ok {
		timeformat[smfJSONKeys.timecode] = map[string]any{
			smfJSONKeys.framespersecond: tc.FramesPerSecond,
			smfJSONKeys.subframes:       tc.SubFrames,
		}
	}

	all[smfJSONKeys.timeformat] = timeformat

	var tracks [][]map[string]any

	for _, tr := range s.Tracks {

		var track []map[string]any

		var deltaoffset uint32

		for _, ev := range tr {

			var msg = map[string]any{}
			msg[smfJSONKeys.delta] = ev.Delta + deltaoffset

			var channel, byte1, byte2 uint8
			var pbrelative int16
			var pbabsolute, seqnumber uint16
			var data []byte
			var text string
			var key, num, denom uint8
			var isMajor, isFlat bool
			var hour, minute, second, frame, fractframe uint8
			var bpm float64

			switch {

			case ev.Message.GetAfterTouch(&channel, &byte1):
				msg[smfJSONKeys.typ] = smfJSONKeys.aftertouchType
				msg[smfJSONKeys.data] = map[string]any{
					smfJSONKeys.channel:  channel,
					smfJSONKeys.pressure: byte1,
				}

			case ev.Message.GetControlChange(&channel, &byte1, &byte2):
				msg[smfJSONKeys.typ] = smfJSONKeys.controlchangeType
				msg[smfJSONKeys.data] = map[string]any{
					smfJSONKeys.channel:        channel,
					smfJSONKeys.controller:     byte1,
					smfJSONKeys.value:          byte2,
					smfJSONKeys.controllername: midi.ControlChangeName[byte1],
				}

			case ev.Message.GetNoteOn(&channel, &byte1, &byte2):
				msg[smfJSONKeys.typ] = smfJSONKeys.noteonType
				msg[smfJSONKeys.data] = map[string]any{
					smfJSONKeys.channel:  channel,
					smfJSONKeys.key:      byte1,
					smfJSONKeys.keyname:  midi.Note(byte1).String(),
					smfJSONKeys.velocity: byte2,
				}

			case ev.Message.GetNoteOff(&channel, &byte1, &byte2):
				msg[smfJSONKeys.typ] = smfJSONKeys.noteoffType
				msg[smfJSONKeys.data] = map[string]any{
					smfJSONKeys.channel: channel,
					smfJSONKeys.key:     byte1,
					smfJSONKeys.keyname: midi.Note(byte1).String(),
				}

			case ev.Message.GetPitchBend(&channel, &pbrelative, &pbabsolute):
				msg[smfJSONKeys.typ] = smfJSONKeys.pitchbendType
				msg[smfJSONKeys.data] = map[string]any{
					smfJSONKeys.channel:  channel,
					smfJSONKeys.relative: pbrelative,
				}

			case ev.Message.GetPolyAfterTouch(&channel, &byte1, &byte2):
				msg[smfJSONKeys.typ] = smfJSONKeys.polyaftertouchType
				msg[smfJSONKeys.data] = map[string]any{
					smfJSONKeys.channel:  channel,
					smfJSONKeys.key:      byte1,
					smfJSONKeys.pressure: byte2,
					smfJSONKeys.keyname:  midi.Note(byte1).String(),
				}

			case ev.Message.GetProgramChange(&channel, &byte1):
				msg[smfJSONKeys.typ] = smfJSONKeys.programchangeType
				msg[smfJSONKeys.data] = map[string]any{
					smfJSONKeys.channel:     channel,
					smfJSONKeys.program:     byte1,
					smfJSONKeys.programname: gm.Instr(byte1).String(),
				}

			case ev.Message.GetSysEx(&data):
				msg[smfJSONKeys.typ] = smfJSONKeys.sysexType
				msg[smfJSONKeys.data] = hex.EncodeToString(data)

			case ev.Message.GetMetaChannel(&channel):
				msg[smfJSONKeys.typ] = smfJSONKeys.channelType
				msg[smfJSONKeys.data] = channel

			case ev.Message.GetMetaCopyright(&text):
				msg[smfJSONKeys.typ] = smfJSONKeys.copyrightType
				msg[smfJSONKeys.data] = text

			case ev.Message.GetMetaCuepoint(&text):
				msg[smfJSONKeys.typ] = smfJSONKeys.cuepointType
				msg[smfJSONKeys.data] = text

			case ev.Message.GetMetaDevice(&text):
				msg[smfJSONKeys.typ] = smfJSONKeys.deviceType
				msg[smfJSONKeys.data] = text

			case ev.Message.GetMetaInstrument(&text):
				msg[smfJSONKeys.typ] = smfJSONKeys.instrumentType
				msg[smfJSONKeys.data] = text

			case ev.Message.GetMetaKeySig(&key, &num, &isMajor, &isFlat):
				msg[smfJSONKeys.typ] = smfJSONKeys.keysignatureType
				msg[smfJSONKeys.data] = map[string]any{
					smfJSONKeys.key:     key,
					smfJSONKeys.num:     num,
					smfJSONKeys.isMajor: isMajor,
					smfJSONKeys.isFlat:  isFlat,
				}

			case ev.Message.GetMetaLyric(&text):
				msg[smfJSONKeys.typ] = smfJSONKeys.lyricType
				msg[smfJSONKeys.data] = text

			case ev.Message.GetMetaMarker(&text):
				msg[smfJSONKeys.typ] = smfJSONKeys.markerType
				msg[smfJSONKeys.data] = text

			case ev.Message.GetMetaPort(&byte1):
				msg[smfJSONKeys.typ] = smfJSONKeys.portType
				msg[smfJSONKeys.data] = byte1

			case ev.Message.GetMetaProgramName(&text):
				msg[smfJSONKeys.typ] = smfJSONKeys.programnameType
				msg[smfJSONKeys.data] = text

			case ev.Message.GetMetaSMPTEOffsetMsg(&hour, &minute, &second, &frame, &fractframe):
				msg[smfJSONKeys.typ] = smfJSONKeys.smpteoffsetType
				msg[smfJSONKeys.data] = map[string]any{
					smfJSONKeys.hour:       hour,
					smfJSONKeys.minute:     minute,
					smfJSONKeys.second:     second,
					smfJSONKeys.frame:      frame,
					smfJSONKeys.fractframe: fractframe,
				}

			case ev.Message.GetMetaSeqData(&data):
				msg[smfJSONKeys.typ] = smfJSONKeys.seqdataType
				msg[smfJSONKeys.data] = hex.EncodeToString(data)

			case ev.Message.GetMetaSeqNumber(&seqnumber):
				msg[smfJSONKeys.typ] = smfJSONKeys.seqnumberType
				msg[smfJSONKeys.data] = seqnumber

			case ev.Message.GetMetaTempo(&bpm):
				msg[smfJSONKeys.typ] = smfJSONKeys.tempoType
				msg[smfJSONKeys.data] = bpm

			case ev.Message.GetMetaText(&text):
				msg[smfJSONKeys.typ] = smfJSONKeys.textType
				msg[smfJSONKeys.data] = text

			case ev.Message.GetMetaTimeSig(&num, &denom, &byte1, &byte2):
				msg[smfJSONKeys.typ] = smfJSONKeys.timesignatureType
				msg[smfJSONKeys.data] = map[string]any{
					smfJSONKeys.num:                      num,
					smfJSONKeys.denom:                    denom,
					smfJSONKeys.clockspertick:            byte1,
					smfJSONKeys.demisemiquaverperquarter: byte2,
				}

			case ev.Message.GetMetaTrackName(&text):
				msg[smfJSONKeys.typ] = smfJSONKeys.tracknameType
				msg[smfJSONKeys.data] = text

			default:
				deltaoffset = ev.Delta + deltaoffset

				// ignore the (invalid) message but not the delta movement, to keep everything in place
				continue

			}

			track = append(track, msg)
		}

		tracks = append(tracks, track)
	}

	all[smfJSONKeys.tracks] = tracks
	return all
}

var _ json.Marshaler = &SMF{}
var _ json.Unmarshaler = &SMF{}
