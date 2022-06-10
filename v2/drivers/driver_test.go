package drivers_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"

	"gitlab.com/gomidi/midi/v2/drivers/midicatdrv"
	"gitlab.com/gomidi/midi/v2/drivers/portmididrv"
	"gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
	"gitlab.com/gomidi/midi/v2/drivers/testdrv"

	"gitlab.com/gomidi/midi/v2/internal/runningstatus"
)

func mkWithRunningStatus(t *testing.T) []byte {

	var bf bytes.Buffer
	wr := runningstatus.NewLiveWriter(&bf)
	wr.Write(midi.NoteOn(2, 65, 120))
	wr.Write(midi.NoteOn(2, 55, 120))
	wr.Write(midi.NoteOn(2, 65, 0))
	wr.Write(midi.NoteOn(1, 65, 20))
	wr.Write(midi.NoteOn(2, 55, 0))
	wr.Write(midi.NoteOn(1, 65, 0))

	bts := bf.Bytes()
	got := fmt.Sprintf("% X", bts)

	expected := `92 41 78 37 78 41 00 91 41 14 92 37 00 91 41 00`

	if got != expected {
		t.Fatalf("\nexpected: %q\n     got: %q", expected, got)
	}

	return bts

}

func RunningStatusTest(t *testing.T, in drivers.In, out drivers.Out) {
	//midi.CloseDriver()
	var got bytes.Buffer

	src := mkWithRunningStatus(t)

	in.Open()
	out.Open()

	var conf drivers.ListenConfig
	stop, _ := in.Listen(func(msg []byte, ms int32) {
		got.WriteString(fmt.Sprintf("% X\n", msg))
	}, conf)

	time.Sleep(30 * time.Millisecond)
	for _, s := range src {
		out.Send([]byte{s})
	}
	//out.Send(src)
	time.Sleep(30 * time.Millisecond)
	stop()
	in.Close()
	out.Close()

	expected := `92 41 78
92 37 78
92 41 00
91 41 14
92 37 00
91 41 00
`

	if got.String() != expected {
		t.Errorf("\nexpected: \n%s\n     got: \n%s\n", expected, got.String())
	}
	//fmt.Printf("%s\n", got.String())
}

func FullStatusTest(t *testing.T, in drivers.In, out drivers.Out) {
	//midi.CloseDriver()
	var got bytes.Buffer

	in.Open()
	out.Open()

	var conf drivers.ListenConfig

	stop, err := in.Listen(func(msg []byte, ms int32) {
		got.WriteString(fmt.Sprintf("% X\n", msg))
	}, conf)

	if err != nil {
		t.Fatalf("ERROR: %s", err.Error())
	}

	time.Sleep(30 * time.Millisecond)

	err = out.Send(midi.NoteOn(2, 65, 120))
	if err != nil {
		t.Fatalf("ERROR: %s", err.Error())
	}
	err = out.Send(midi.NoteOn(2, 55, 120))
	if err != nil {
		t.Fatalf("ERROR: %s", err.Error())
	}
	err = out.Send(midi.NoteOn(2, 65, 0))
	if err != nil {
		t.Fatalf("ERROR: %s", err.Error())
	}
	err = out.Send(midi.NoteOn(1, 65, 20))
	if err != nil {
		t.Fatalf("ERROR: %s", err.Error())
	}
	err = out.Send(midi.NoteOn(2, 55, 0))
	if err != nil {
		t.Fatalf("ERROR: %s", err.Error())
	}
	err = out.Send(midi.NoteOn(1, 65, 0))
	if err != nil {
		t.Fatalf("ERROR: %s", err.Error())
	}

	time.Sleep(30 * time.Millisecond)
	stop()
	err = in.Close()
	if err != nil {
		t.Fatalf("ERROR: %s", err.Error())
	}
	err = out.Close()
	if err != nil {
		t.Fatalf("ERROR: %s", err.Error())
	}

	expected := `92 41 78
92 37 78
92 41 00
91 41 14
92 37 00
91 41 00
`

	if got.String() != expected {
		t.Errorf("\nexpected: \n%s\n     got: \n%s\n", expected, got.String())
	}

}

func TestRunningStatusForTestDrv(t *testing.T) {
	drv := testdrv.New("runningstatus")
	ins, _ := drv.Ins()
	outs, _ := drv.Outs()

	in := ins[0]
	out := outs[0]
	RunningStatusTest(t, in, out)
	drv.Close()
}

func TestRunningStatusForRtMIDIDrv(t *testing.T) {
	drv, err := rtmididrv.New()
	if err != nil {
		t.Fatalf("ERROR: %s", err.Error())
	}
	in, err := midi.FindInPort("Midi Through Port-0")
	if err != nil {
		//t.Fatalf("ERROR: %s", err.Error())
		t.Skipf("could not find in port Midi Through Port-0")
		return
	}
	out, err := midi.FindOutPort("Midi Through Port-0")
	if err != nil {
		t.Skipf("could not find out port Midi Through Port-0")
		return
		//t.Fatalf("ERROR: %s", err.Error())
	}

	RunningStatusTest(t, in, out)
	drv.Close()
}

func TestFullStatusForTestDrv(t *testing.T) {
	drv := testdrv.New("fullstatus")
	ins, _ := drv.Ins()
	outs, _ := drv.Outs()

	in := ins[0]
	out := outs[0]
	FullStatusTest(t, in, out)
	drv.Close()
}

func TestFullStatusForRtMIDIDrv(t *testing.T) {
	drv, err := rtmididrv.New()
	if err != nil {
		t.Fatalf("ERROR: %s", err.Error())
	}
	in, err := midi.FindInPort("Midi Through Port-0")
	if err != nil {
		//t.Fatalf("ERROR: %s", err.Error())
		t.Skipf("could not find in port Midi Through Port-0")
		return
	}
	out, err := midi.FindOutPort("Midi Through Port-0")
	if err != nil {
		t.Skipf("could not find out port Midi Through Port-0")
		return
		//t.Fatalf("ERROR: %s", err.Error())
	}

	FullStatusTest(t, in, out)
	drv.Close()
}

func TestFullStatusForPortMIDIDrv(t *testing.T) {
	drv, err := portmididrv.New()
	if err != nil {
		t.Fatalf("ERROR: %s", err.Error())
	}

	in, err := midi.FindInPort("Midi Through Port-0")
	if err != nil {
		//t.Fatalf("ERROR: %s", err.Error())
		t.Skipf("could not find in port Midi Through Port-0")
		return
	}
	out, err := midi.FindOutPort("Midi Through Port-0")
	if err != nil {
		t.Skipf("could not find out port Midi Through Port-0")
		return
		//t.Fatalf("ERROR: %s", err.Error())
	}

	FullStatusTest(t, in, out)
	drv.Close()
}

func TestRunningStatusForPortMIDIDrv(t *testing.T) {
	drv, err := portmididrv.New()
	if err != nil {
		t.Fatalf("ERROR: %s", err.Error())
	}

	in, err := midi.FindInPort("Midi Through Port-0")
	if err != nil {
		//t.Fatalf("ERROR: %s", err.Error())
		t.Skipf("could not find in port Midi Through Port-0")
		return
	}
	out, err := midi.FindOutPort("Midi Through Port-0")
	if err != nil {
		t.Skipf("could not find out port Midi Through Port-0")
		return
		//t.Fatalf("ERROR: %s", err.Error())
	}

	RunningStatusTest(t, in, out)
	drv.Close()
}

func TestFullStatusForMidicatDrv(t *testing.T) {
	drv, err := midicatdrv.New()
	if err != nil {
		t.Fatalf("ERROR: %s", err.Error())
	}

	in, err := midi.FindInPort("Midi Through Port-0")
	if err != nil {
		//t.Fatalf("ERROR: %s", err.Error())
		t.Skipf("could not find in port Midi Through Port-0")
		return
	}
	out, err := midi.FindOutPort("Midi Through Port-0")
	if err != nil {
		t.Skipf("could not find out port Midi Through Port-0")
		return
		//t.Fatalf("ERROR: %s", err.Error())
	}

	//RunningStatusTest(t, in, out)
	FullStatusTest(t, in, out)
	drv.Close()
}

func TestRunningStatusForMidicatDrv(t *testing.T) {
	drv, err := midicatdrv.New()
	if err != nil {
		t.Fatalf("ERROR: %s", err.Error())
	}

	in, err := midi.FindInPort("Midi Through Port-0")
	if err != nil {
		//t.Fatalf("ERROR: %s", err.Error())
		t.Skipf("could not find in port Midi Through Port-0")
		return
	}
	out, err := midi.FindOutPort("Midi Through Port-0")
	if err != nil {
		t.Skipf("could not find out port Midi Through Port-0")
		return
		//t.Fatalf("ERROR: %s", err.Error())
	}

	RunningStatusTest(t, in, out)
	drv.Close()
}
