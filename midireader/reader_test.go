package midireader

import (
	"bytes"
	"io"
	"testing"

	"github.com/gomidi/midi/messages/channel"
	"github.com/gomidi/midi/messages/realtime"
	"github.com/gomidi/midi/messages/syscommon"
	"github.com/gomidi/midi/messages/sysex"
	"github.com/gomidi/midi/midiwriter"
)

func mkMIDI() io.Reader {
	var bf bytes.Buffer

	wr := midiwriter.New(&bf)

	wr.Write(channel.New(1).NoteOn(65, 100))
	wr.Write(realtime.Start)
	wr.Write(sysex.Start([]byte{0x50}))
	wr.Write(channel.New(1).NoteOffVelocity(65, 64))
	wr.Write(syscommon.TuneRequest)
	wr.Write(channel.New(2).NoteOn(62, 30))
	bf.Write([]byte{0xF5, 0x51, 0x52})
	wr.Write(sysex.SysEx([]byte{0x50, 0x51}))
	wr.Write(channel.New(2).NoteOn(62, 0))

	return bytes.NewReader(bf.Bytes())
}

func TestRead(t *testing.T) {

	var bf bytes.Buffer

	bf.WriteString("\n")

	rtCallBack := func(m realtime.Message) {
		bf.WriteString("Realtime: " + m.String() + "\n")
	}
	rd := New(mkMIDI(), rtCallBack)

	for {
		ev, err := rd.Read()
		if err != nil {
			break
		}
		bf.WriteString(ev.String() + "\n")
	}

	expected := `
channel.NoteOn channel 1 key 65 velocity 100
Realtime: Start
sysex.SysEx len: 1
channel.NoteOff channel 1 key 65
syscommon.tuneRequest
channel.NoteOn channel 2 key 62 velocity 30
sysex.SysEx len: 2
channel.NoteOff channel 2 key 62
`
	if got, wanted := bf.String(), expected; got != wanted {
		t.Errorf("got:\n%s\n\nwanted:\n%s\n\n", got, wanted)
	}

}

func TestReadNoteOffVelocity(t *testing.T) {

	var bf bytes.Buffer

	bf.WriteString("\n")

	rtCallBack := func(m realtime.Message) {
		bf.WriteString("Realtime: " + m.String() + "\n")
	}
	rd := New(mkMIDI(), rtCallBack, NoteOffPedantic())

	for {
		ev, err := rd.Read()
		if err != nil {
			break
		}
		bf.WriteString(ev.String() + "\n")
	}

	expected := `
channel.NoteOn channel 1 key 65 velocity 100
Realtime: Start
sysex.SysEx len: 1
channel.NoteOffVelocity channel 1 key 65 velocity 64
syscommon.tuneRequest
channel.NoteOn channel 2 key 62 velocity 30
sysex.SysEx len: 2
channel.NoteOff channel 2 key 62
`
	if got, wanted := bf.String(), expected; got != wanted {
		t.Errorf("got:\n%s\n\nwanted:\n%s\n\n", got, wanted)
	}

}
