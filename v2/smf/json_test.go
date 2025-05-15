package smf

import (
	"strings"
	"testing"

	"gitlab.com/gomidi/midi/v2"
)

func TestJSONUnMarshal(t *testing.T) {
	var smf SMF
	//	t.Skip()
	err := smf.UnmarshalJSON([]byte(jsonStr))
	if err != nil {
		t.Fatalf("error: %s\n", err.Error())
		return
	}

	if smf.Format() != 1 {
		t.Fatalf("wrong format: %v // expected 1", smf.Format())
		return
	}

	mt, ok := smf.TimeFormat.(MetricTicks)
	if !ok {
		t.Fatalf("wrong timeformat type: %s (expected MetricTicks)\n", smf.TimeFormat.String())
		return
	}

	if mt.Resolution() != 480 {
		t.Fatalf("wrong metric ticks resolution: %s // expected: %v\n", smf.TimeFormat.String(), 480)
		return
	}

	if smf.NumTracks() != 2 {
		t.Fatalf("wrong number of tracks: %v // expected: %v\n", smf.NumTracks(), 2)
	}

	if len(smf.Tracks) != 2 {
		t.Fatalf("wrong number of tracks: %v // expected: %v\n", smf.NumTracks(), 2)
	}

	t1 := smf.Tracks[0]
	if len(t1) != 25 {
		t.Fatalf("wrong number of events in track1: %v // expected: %v\n", len(t1), 25)
	}

	ev1 := t1[0]
	//ev2 := t1[1]

	if ev1.Delta != 24 {
		t.Fatalf("wrong delta of event 1 in track1: %v // expected: %v\n", ev1.Delta, 24)
	}

	msgTests := []string{
		"AfterTouch channel: 3 pressure: 123",
		"ControlChange channel: 3 controller: 40 value: 100",
		"NoteOn channel: 3 key: 12 velocity: 120",
		"NoteOff channel: 3 key: 12",
		"PitchBend channel: 3 pitch: 200 (8392)",
		"PolyAfterTouch channel: 3 key: 15 pressure: 80",
		"ProgramChange channel: 3 program: 33",
		"SysExType data: DD A1 BB",
		"MetaChannel channel: 16",
		"MetaCopyright text: \"copyright\"",
		"MetaCuepoint text: \"Cuepoint\"",
		"MetaDevice text: \"Device\"",
		"MetaInstrument text: \"Instrument\"",
		"MetaKeySig key: FMaj",
		"MetaLyric text: \"lyric\"",
		"MetaMarker text: \"marker\"",
		"MetaPort port: 100",
		"MetaProgramName text: \"program\"",
		"MetaSMPTEOffset hour: 1 minute: 2 second: 3 frame: 4 fractframe: 5",
		"MetaSeqData bytes: DD A1 BC",
		"MetaSeqNumber number: 123",
		"MetaTempo bpm: 60.00",
		"MetaText text: \"text\"",
		"MetaTimeSig meter: 3/4",
		"MetaTrackName text: \"trackname\"",
	}

	for i, expectedMsg := range msgTests {
		evx := t1[i]
		gotmsg := evx.Message.String()
		if gotmsg != expectedMsg {
			t.Fatalf("wrong message of event %v in track 1: %q // expected: %q\n", i, gotmsg, expectedMsg)
		}
	}

	t2 := smf.Tracks[1]
	if len(t2) != 1 {
		t.Fatalf("wrong number of events in track2: %v // expected: %v\n", len(t2), 1)
	}

	ev3 := t2[0]
	if ev3.Delta != 20 {
		t.Fatalf("wrong delta of event 1 in track2: %v // expected: %v\n", ev3.Delta, 20)
	}

	msg3 := ev3.Message.String()
	expectedMsg3 := "ControlChange channel: 4 controller: 10 value: 110"
	if msg3 != expectedMsg3 {
		t.Fatalf("wrong message of event 3 in track1: %q // expected: %q\n", msg3, expectedMsg3)
	}
}

var jsonStr = strings.ReplaceAll(strings.TrimSpace(`
{
	"format": 1,
	"timeformat": {
		"metricticks": 480
	},
	"tracks": [
		[
			{
				"data": {
					"channel": 3,
					"pressure": 123
				},
				"delta": 24,
				"type": "aftertouch"
			},
			{
				"data": {
					"channel": 3,
					"controller": 40,
					"controllername": "Balance (LSB)",
					"value": 100
				},
				"delta": 24,
				"type": "controlchange"
			},
			{
				"data": {
					"channel": 3,
					"key": 12,
					"keyname": "C1",
					"velocity": 120
				},
				"delta": 24,
				"type": "noteon"
			},
			{
				"data": {
					"channel": 3,
					"key": 12,
					"keyname": "C1"
				},
				"delta": 24,
				"type": "noteoff"
			},
			{
				"data": {
					"channel": 3,
					"relative": 200
				},
				"delta": 24,
				"type": "pitchbend"
			},
			{
				"data": {
					"channel": 3,
					"key": 15,
					"keyname": "Eb1",
					"pressure": 80
				},
				"delta": 24,
				"type": "polyaftertouch"
			},
			{
				"data": {
					"channel": 3,
					"program": 33,
					"programname": "ElectricBassFinger"
				},
				"delta": 24,
				"type": "programchange"
			},
			{
				"data": "dda1bb",
				"delta": 24,
				"type": "sysex"
			},
			{
				"data": 16,
				"delta": 24,
				"type": "channel"
			},
			{
				"data": "copyright",
				"delta": 24,
				"type": "copyright"
			},
			{
				"data": "Cuepoint",
				"delta": 24,
				"type": "cuepoint"
			},
			{
				"data": "Device",
				"delta": 24,
				"type": "device"
			},
			{
				"data": "Instrument",
				"delta": 24,
				"type": "instrument"
			},
			{
				"data": {
					"isFlat": true,
					"isMajor": true,
					"key": 5,
					"num": 1
				},
				"delta": 24,
				"type": "keysignature"
			},
			{
				"data": "lyric",
				"delta": 24,
				"type": "lyric"
			},
			{
				"data": "marker",
				"delta": 24,
				"type": "marker"
			},
			{
				"data": 100,
				"delta": 24,
				"type": "port"
			},
			{
				"data": "program",
				"delta": 24,
				"type": "programname"
			},
			{
				"data": {
					"fractframe": 5,
					"frame": 4,
					"hour": 1,
					"minute": 2,
					"second": 3
				},
				"delta": 24,
				"type": "smpteoffset"
			},
			{
				"data": "dda1bc",
				"delta": 24,
				"type": "seqdata"
			},
			{
				"data": 123,
				"delta": 24,
				"type": "seqnumber"
			},
			{
				"data": 60,
				"delta": 24,
				"type": "tempo"
			},
			{
				"data": "text",
				"delta": 24,
				"type": "text"
			},
			{
				"data": {
					"clockspertick": 12,
					"demisemiquaverperquarter": 14,
					"denom": 4,
					"num": 3
				},
				"delta": 24,
				"type": "timesignature"
			},
			{
				"data": "trackname",
				"delta": 24,
				"type": "trackname"
			}
		],
		[
			{
				"data": {
					"channel": 4,
					"controller": 10,
					"controllername": "Pan position (MSB)",
					"value": 110
				},
				"delta": 20,
				"type": "controlchange"
			}
		]
	]
}
`), "\t", "  ")

func TestJSONMarshal(t *testing.T) {
	smf1 := NewSMF1()
	smf1.TimeFormat = MetricTicks(480)

	var t1, t2 Track
	t1.Add(24, midi.AfterTouch(3, 123))
	t1.Add(24, midi.ControlChange(3, 40, 100))
	t1.Add(24, midi.NoteOn(3, 12, 120))
	t1.Add(24, midi.NoteOff(3, 12))
	t1.Add(24, midi.Pitchbend(3, 200))
	t1.Add(24, midi.PolyAfterTouch(3, 15, 80))
	t1.Add(24, midi.ProgramChange(3, 33))
	t1.Add(24, midi.SysEx([]byte{0xDD, 0xA1, 0xBB}))
	t1.Add(24, MetaChannel(16))
	t1.Add(24, MetaCopyright("copyright"))
	t1.Add(24, MetaCuepoint("Cuepoint"))
	t1.Add(24, MetaDevice("Device"))
	t1.Add(24, MetaInstrument("Instrument"))
	t1.Add(24, MetaKey(3, true, 1, true))
	t1.Add(24, MetaLyric("lyric"))
	t1.Add(24, MetaMarker("marker"))
	t1.Add(24, MetaPort(100))
	t1.Add(24, MetaProgram("program"))
	t1.Add(24, MetaSMPTE(1, 2, 3, 4, 5))
	t1.Add(24, MetaSequencerData([]byte{0xDD, 0xA1, 0xBC}))
	t1.Add(24, MetaSequenceNo(123))
	t1.Add(24, MetaTempo(60.0))
	t1.Add(24, MetaText("text"))
	t1.Add(24, MetaTimeSig(3, 4, 12, 14))
	t1.Add(24, MetaTrackSequenceName("trackname"))

	t2.Add(20, midi.ControlChange(4, 10, 110))

	smf1.Add(t1)
	smf1.Add(t2)

	res, err := smf1.MarshalJSONIndent()

	if err != nil {
		t.Fatalf("error: %s\n", err.Error())
	}

	got := strings.TrimSpace(string(res))
	expected := jsonStr
	if got != expected {
		t.Errorf("got:\n%s\nexpected:\n%s\n", got, expected)
	}
}
