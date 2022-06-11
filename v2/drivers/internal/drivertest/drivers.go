package drivertest

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	"gitlab.com/gomidi/midi/v2/internal/runningstatus"
	"gitlab.com/gomidi/midi/v2/sysex"
)

func mkWithRunningStatus() (messages [][]byte) {

	var bf bytes.Buffer
	wr := runningstatus.NewLiveWriter(&bf)
	wr.Write(midi.NoteOn(2, 65, 120))
	messages = append(messages, bf.Bytes())
	bf = bytes.Buffer{}
	wr.Write(midi.NoteOn(2, 55, 120))
	messages = append(messages, bf.Bytes())
	bf = bytes.Buffer{}
	wr.Write(midi.NoteOn(2, 65, 0))
	messages = append(messages, bf.Bytes())
	bf = bytes.Buffer{}
	wr.Write(midi.NoteOn(1, 65, 20))
	messages = append(messages, bf.Bytes())
	bf = bytes.Buffer{}
	wr.Write(midi.NoteOn(2, 55, 0))
	messages = append(messages, bf.Bytes())
	bf = bytes.Buffer{}
	wr.Write(midi.NoteOn(1, 65, 0))
	messages = append(messages, bf.Bytes())
	bf = bytes.Buffer{}

	//bts := bytes.Join(messages, []byte{255})
	//fmt.Printf("% X\n", bts)

	/*
			bts := bf.Bytes()
			got := fmt.Sprintf("% X", bts)

			expected := `92 41 78 37 78 41 00 91 41 14 92 37 00 91 41 00`

			if got != expected {
				t.Fatalf("\nexpected: %q\n     got: %q", expected, got)
			}
		return bts
	*/

	return
}

// 5. The midi.In ports respects any running status bytes when converting to midi.Message(s)
// 6. The midi.Out ports accept running status bytes.
func RunningStatusTest(t *testing.T, in drivers.In, out drivers.Out) {
	//midi.CloseDriver()
	var got bytes.Buffer

	src := mkWithRunningStatus()

	in.Open()
	out.Open()

	var conf drivers.ListenConfig
	stop, _ := in.Listen(func(msg []byte, ms int32) {
		got.WriteString(fmt.Sprintf("% X\n", msg))
	}, conf)

	time.Sleep(30 * time.Millisecond)
	for _, s := range src {
		//out.Send([]byte{s})
		out.Send(s)
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
	conf.ActiveSense = true
	conf.TimeCode = true

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
	err = out.Send(midi.Activesense())
	if err != nil {
		t.Fatalf("ERROR: %s", err.Error())
	}
	err = out.Send(midi.TimingClock())
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
FE
F8
`

	if got.String() != expected {
		t.Errorf("\nexpected: \n%s\n     got: \n%s\n", expected, got.String())
	}

}

/*
[x] 1. Implement the midi.Driver, midi.In and midi.Out interfaces
[x] 2. Autoregister via midi.RegisterDriver(drv) within init function
[x] 3. The String() method must return a unique and indentifying string
4. The driver expects multiple goroutines to use its writing and reading methods and locks accordingly (to be threadsafe).
[x] 5. The midi.In ports respects any running status bytes when converting to midi.Message(s)
[x] 6. The midi.Out ports accept running status bytes.
7. The New function may take optional driver specific arguments. These are all organized as functional optional arguments.
8. The midi.In port is responsible for buffering of sysex data. Only complete sysex data is passed to the listener.
9. The midi In port must check the ListenConfig and act accordingly.
10. incomplete sysex data must be cached inside the sender and flushed, if the data is complete.

*/

// 1. Implement the midi.Driver, midi.In and midi.Out interfaces
func DriverInterfaceImplementationTest(t *testing.T, drv interface{}) {
	d := drv.(drivers.Driver)
	_ = d
}

// 2. Autoregister via midi.RegisterDriver(drv) within init function
// 3. The String() method must return a unique and indentifying string
func AutoregisterTest(t *testing.T, drv drivers.Driver) {
	_, has := drivers.REGISTRY[drv.String()]

	if !has {
		t.Errorf("driver %T did not autoregister", drv)
	}
}

// 9. The midi In port must check the ListenConfig and act accordingly.
func NoActiveSenseTest(t *testing.T, in drivers.In, out drivers.Out) {
	//midi.CloseDriver()
	var got bytes.Buffer

	in.Open()
	out.Open()

	var conf drivers.ListenConfig
	conf.ActiveSense = false
	conf.TimeCode = true

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
	err = out.Send(midi.Activesense())
	if err != nil {
		t.Fatalf("ERROR: %s", err.Error())
	}
	err = out.Send(midi.TimingClock())
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
F8
`

	if got.String() != expected {
		t.Errorf("\nexpected: \n%s\n     got: \n%s\n", expected, got.String())
	}

}

// 9. The midi In port must check the ListenConfig and act accordingly.
func NoTimeCodeTest(t *testing.T, in drivers.In, out drivers.Out) {
	//midi.CloseDriver()
	var got bytes.Buffer

	in.Open()
	out.Open()

	var conf drivers.ListenConfig
	conf.ActiveSense = true
	conf.TimeCode = false

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
	err = out.Send(midi.Activesense())
	if err != nil {
		t.Fatalf("ERROR: %s", err.Error())
	}
	err = out.Send(midi.TimingClock())
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
FE
`

	if got.String() != expected {
		t.Errorf("\nexpected: \n%s\n     got: \n%s\n", expected, got.String())
	}

}

func mkSysex() []byte {
	return sysex.GMSystem(2, true)
}

func SysexTest(t *testing.T, in drivers.In, out drivers.Out) {
	var got bytes.Buffer

	in.Open()
	out.Open()

	var conf drivers.ListenConfig
	conf.SysEx = true

	sys := mkSysex()

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
	err = out.Send(sys)
	if err != nil {
		t.Fatalf("ERROR: %s", err.Error())
	}

	time.Sleep(100 * time.Millisecond)
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
F0 7E 02 09 01 F7
`

	if got.String() != expected {
		t.Errorf("\nexpected: \n%s\n     got: \n%s\n", expected, got.String())
	}

}

func NoSysexTest(t *testing.T, in drivers.In, out drivers.Out) {
	var got bytes.Buffer

	in.Open()
	out.Open()

	var conf drivers.ListenConfig
	conf.SysEx = false

	sys := mkSysex()
	//didsend := fmt.Sprintf("% X\n", sys)
	//fmt.Println(didsend)

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
	err = out.Send(sys)
	if err != nil {
		t.Fatalf("ERROR: %s", err.Error())
	}

	time.Sleep(100 * time.Millisecond)
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
`

	if got.String() != expected {
		t.Errorf("\nexpected: \n%s\n     got: \n%s\n", expected, got.String())
	}

}
