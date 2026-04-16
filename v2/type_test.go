package midi

import "testing"

func TestIsType(t *testing.T) {

	tests := []struct {
		input    []byte
		expected []Type
	}{
		{ // 0
			AfterTouch(1, 120),
			[]Type{AfterTouchMsg, ChannelMsg},
		},
		{ // 1
			ControlChange(8, 7, 110),
			[]Type{ControlChangeMsg, ChannelMsg},
		},
		{ // 2
			NoteOn(2, 100, 80),
			[]Type{NoteOnMsg, ChannelMsg},
		},
		{ // 3
			NoteOff(3, 80),
			[]Type{NoteOffMsg, ChannelMsg},
		},
		{
			NoteOffVelocity(4, 80, 20),
			[]Type{NoteOffMsg, ChannelMsg},
		},
		{
			Pitchbend(4, 300),
			[]Type{PitchBendMsg, ChannelMsg},
		},
		{
			PolyAfterTouch(4, 86, 109),
			[]Type{PolyAfterTouchMsg, ChannelMsg},
		},
		{
			ProgramChange(4, 83),
			[]Type{ProgramChangeMsg, ChannelMsg},
		},
		//
		{
			TimingClock(),
			[]Type{TimingClockMsg, RealTimeMsg},
		},
		{
			Tick(),
			[]Type{TickMsg, RealTimeMsg},
		},
		{
			Start(),
			[]Type{StartMsg, RealTimeMsg},
		},
		{
			Continue(),
			[]Type{ContinueMsg, RealTimeMsg},
		},
		{
			Stop(),
			[]Type{StopMsg, RealTimeMsg},
		},
		/*
			{
				Unknown(),
				[]Type{UnknownMsg, RealTimeMsg},
			},
		*/
		{
			Activesense(),
			[]Type{ActiveSenseMsg, RealTimeMsg},
		},
		{
			Reset(),
			[]Type{ResetMsg, RealTimeMsg},
		},
		{
			MTC(4),
			[]Type{MTCMsg, SysCommonMsg},
		},
		{
			SPP(4),
			[]Type{SPPMsg, SysCommonMsg},
		},
		{
			SongSelect(4),
			[]Type{SongSelectMsg, SysCommonMsg},
		},
		{
			Tune(),
			[]Type{TuneMsg, SysCommonMsg},
		},
		{
			SysEx([]byte{2, 3, 4}),
			[]Type{SysExMsg},
		},
	}

	for i, test := range tests {

		msg := Message(test.input)

		for _, ty := range test.expected {
			if !msg.Is(ty) {
				t.Errorf("[%v] expected %q to be of type %q but is not", i, msg, ty)
			}
		}
	}

}

func TestIsNotType(t *testing.T) {

	tests := []struct {
		input        []byte
		not_expected []Type
	}{
		{ // 0
			AfterTouch(1, 120),
			[]Type{metaMsg, RealTimeMsg, UnknownMsg, SysCommonMsg, SysExMsg},
		},
		{ // 1
			ControlChange(8, 7, 110),
			[]Type{metaMsg, RealTimeMsg, UnknownMsg, SysCommonMsg, SysExMsg},
		},
		{ // 2
			NoteOn(2, 100, 80),
			[]Type{metaMsg, RealTimeMsg, UnknownMsg, SysCommonMsg, SysExMsg},
		},
		{ // 3
			NoteOff(3, 80),
			[]Type{metaMsg, RealTimeMsg, UnknownMsg, SysCommonMsg, SysExMsg},
		},
		{
			NoteOffVelocity(4, 80, 20),
			[]Type{metaMsg, RealTimeMsg, UnknownMsg, SysCommonMsg, SysExMsg},
		},
		{
			Pitchbend(4, 300),
			[]Type{metaMsg, RealTimeMsg, UnknownMsg, SysCommonMsg, SysExMsg},
		},
		{
			PolyAfterTouch(4, 86, 109),
			[]Type{metaMsg, RealTimeMsg, UnknownMsg, SysCommonMsg, SysExMsg},
		},
		{
			ProgramChange(4, 83),
			[]Type{metaMsg, RealTimeMsg, UnknownMsg, SysCommonMsg, SysExMsg},
		},
		//
		{
			TimingClock(),
			[]Type{metaMsg, ChannelMsg, UnknownMsg, SysCommonMsg, SysExMsg},
		},
		{
			Tick(),
			[]Type{metaMsg, ChannelMsg, UnknownMsg, SysCommonMsg, SysExMsg},
		},
		{
			Start(),
			[]Type{metaMsg, ChannelMsg, UnknownMsg, SysCommonMsg, SysExMsg},
		},
		{
			Continue(),
			[]Type{metaMsg, ChannelMsg, UnknownMsg, SysCommonMsg, SysExMsg},
		},
		{
			Stop(),
			[]Type{metaMsg, ChannelMsg, UnknownMsg, SysCommonMsg, SysExMsg},
		},
		/*
			{
				Unknown(),
				[]Type{metaMsg, ChannelMsg, UnknownMsg, SysCommonMsg, SysExMsg},
			},
		*/
		{
			Activesense(),
			[]Type{metaMsg, ChannelMsg, UnknownMsg, SysCommonMsg, SysExMsg},
		},
		{
			Reset(),
			[]Type{metaMsg, ChannelMsg, UnknownMsg, SysCommonMsg, SysExMsg},
		},
		//
		{
			MTC(4),
			[]Type{metaMsg, ChannelMsg, UnknownMsg, RealTimeMsg, SysExMsg},
		},
		{
			SPP(4),
			[]Type{metaMsg, ChannelMsg, UnknownMsg, RealTimeMsg, SysExMsg},
		},
		{
			SongSelect(4),
			[]Type{metaMsg, ChannelMsg, UnknownMsg, RealTimeMsg, SysExMsg},
		},
		{
			Tune(),
			[]Type{metaMsg, ChannelMsg, UnknownMsg, RealTimeMsg, SysExMsg},
		},
		{
			SysEx([]byte{2, 3, 4}),
			[]Type{metaMsg, ChannelMsg, UnknownMsg, RealTimeMsg, SysCommonMsg},
		},
	}

	for i, test := range tests {

		msg := Message(test.input)

		for _, ty := range test.not_expected {
			if msg.Is(ty) {
				t.Errorf("[%v] not expected %q to be of type %q but is", i, msg, ty)
			}
		}
	}

}
