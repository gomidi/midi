package hyperarp_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	hyperarp "gitlab.com/gomidi/midi/tools/hyperarp/lib"
	"gitlab.com/gomidi/midi/v2"

	"gitlab.com/gomidi/midi/v2/drivers"
	"gitlab.com/gomidi/midi/v2/drivers/testdrv"
)

type cable struct {
	drivers.Driver
	in  drivers.In
	out drivers.Out
}

func newCable(name string) *cable {
	var c cable
	drv := testdrv.New("fake cable: " + name)
	c.Driver = drv
	ins, _ := c.Driver.Ins()
	outs, _ := c.Driver.Outs()
	c.in, c.out = ins[0], outs[0]
	c.in.Open()
	c.out.Open()
	return &c
}

type arpTester struct {
	arp *hyperarp.Arp
	//rd  *reader.Reader
	readFunc func(msg midi.Message, timestampms int32)
	readStop func()
	//*writer.Writer
	send     func(msg midi.Message) error
	bf       bytes.Buffer
	cable1   *cable
	cable2   *cable
	lastTime time.Time
}

func newArpTester() *arpTester {
	/*
		at.rd = reader.New(
			reader.NoLogger(),
			reader.Each(func(p *reader.Position, msg midi.Message) {
				now := time.Now()
				if at.lastTime.Unix() == 0 {
					at.bf.WriteString(msg.String() + "\n")
				} else {
					dur := now.Sub(at.lastTime)
					at.bf.WriteString(fmt.Sprintf("[%v] %s\n", dur.Milliseconds(), msg.String()))
				}
				at.lastTime = time.Now()
			}),
		)
	*/
	var at arpTester

	at.lastTime = time.Unix(0, 0)

	at.readFunc = func(msg midi.Message, timestampms int32) {
		now := time.Now()
		if at.lastTime.Unix() == 0 {
			at.bf.WriteString(msg.String() + "\n")
		} else {
			dur := now.Sub(at.lastTime)
			at.bf.WriteString(fmt.Sprintf("[%v] %s\n", dur.Milliseconds(), msg.String()))
		}
		at.lastTime = time.Now()
	}

	at.cable1 = newCable("write to arp")
	at.cable2 = newCable("read from arp")

	at.arp = hyperarp.New(at.cable1.in, at.cable2.out)
	//at.Writer = writer.New(at.cable1.out)
	at.send, _ = midi.SendTo(at.cable1.out)
	return &at
}

func (at *arpTester) Run() {
	//go at.rd.ListenTo(at.cable2.in)
	at.readStop, _ = midi.ListenTo(at.cable2.in, at.readFunc)
	//at.cable2.Driver.()
	//at.cable1.Driver
	go at.arp.Run()
}

func (at *arpTester) Close() {
	at.arp.Close()
	at.cable1.Close()
	at.cable2.Close()
}

func (at *arpTester) Result() string {
	return at.bf.String()
}

// This example reads from the first input and and writes to the first output port
func TestFirst(t *testing.T) {
	var a *arpTester

	var tests = []struct {
		fn       func()
		descr    string
		expected string
	}{
		{
			func() {
				a.send(midi.NoteOn(0, 70, 100))
				time.Sleep(200 * time.Millisecond)
				a.send(midi.NoteOff(0, 70))
			},
			"note 70",
			"NoteOn channel: 0 key: 70 velocity: 100\n[167] NoteOff channel: 0 key: 70\n",
		},
		{
			func() {
				a.send(midi.Pitchbend(0, 1000))
				time.Sleep(10 * time.Millisecond)
			},
			"pitchbend passthrough",
			"PitchBend channel: 0 pitch: 1000 (9192)\n",
		},
		{
			func() {
				a.send(midi.Pitchbend(0, 100))
				a.send(midi.AfterTouch(0, 100))
				time.Sleep(10 * time.Millisecond)
			},
			"pitchbend and aftertouch passthrough",
			"PitchBend channel: 0 pitch: 100 (8292)\n[0] AfterTouch channel: 0 pressure: 100\n",
		},
		{
			func() {
				a.send(midi.ControlChange(0, midi.GeneralPurposeSlider1, 3))
				a.send(midi.NoteOn(0, hyperarp.D, 100))
				a.send(midi.NoteOn(0, uint8(12+hyperarp.E), 120))
				time.Sleep(500 * time.Millisecond)
				a.send(midi.NoteOff(0, hyperarp.D))
				a.send(midi.NoteOff(0, uint8(12+hyperarp.E)))
				time.Sleep(10 * time.Millisecond)
			},
			"2 arp notes upward",
			`NoteOn channel: 0 key: 16 velocity: 120
[84] NoteOff channel: 0 key: 16
[41] NoteOn channel: 0 key: 26 velocity: 100
[84] NoteOff channel: 0 key: 26
[41] NoteOn channel: 0 key: 28 velocity: 120
[84] NoteOff channel: 0 key: 28
[41] NoteOn channel: 0 key: 38 velocity: 100
[84] NoteOff channel: 0 key: 38
`,
		},
		{
			func() {
				a.send(midi.ControlChange(0, midi.GeneralPurposeSlider1, 3))
				a.send(midi.NoteOn(0, hyperarp.D, 100))
				a.send(midi.NoteOn(0, hyperarp.G, 80))
				a.send(midi.NoteOn(0, uint8(12+hyperarp.E), 120))
				time.Sleep(530 * time.Millisecond)
				a.send(midi.NoteOff(0, hyperarp.D))
				a.send(midi.NoteOff(0, hyperarp.G))
				a.send(midi.NoteOff(0, uint8(12+hyperarp.E)))
			},
			"3 arp notes upward",
			`NoteOn channel: 0 key: 16 velocity: 120
[84] NoteOff channel: 0 key: 16
[41] NoteOn channel: 0 key: 19 velocity: 80
[84] NoteOff channel: 0 key: 19
[41] NoteOn channel: 0 key: 26 velocity: 100
[84] NoteOff channel: 0 key: 26
[41] NoteOn channel: 0 key: 28 velocity: 120
[84] NoteOff channel: 0 key: 28
`,
		},

		{
			func() {
				a.send(midi.ControlChange(0, midi.GeneralPurposeSlider1, 3))
				a.send(midi.ControlChange(0, midi.GeneralPurposeButton1Switch, midi.On))
				time.Sleep(time.Microsecond)
				a.send(midi.NoteOn(0, hyperarp.D, 100))
				a.send(midi.NoteOn(0, uint8(12+hyperarp.E), 120))
				time.Sleep(500 * time.Millisecond)
				a.send(midi.NoteOff(0, hyperarp.D))
				a.send(midi.NoteOff(0, uint8(12+hyperarp.E)))
				//time.Sleep(time.Microsecond)
				a.send(midi.ControlChange(0, midi.GeneralPurposeButton1Switch, midi.Off))
				time.Sleep(10 * time.Millisecond)
			},
			"2 arp notes downward",
			`NoteOn channel: 0 key: 16 velocity: 120
			[83] NoteOff channel: 0 key: 16
			[41] NoteOn channel: 0 key: 14 velocity: 100
			[83] NoteOff channel: 0 key: 14
			[41] NoteOn channel: 0 key: 4 velocity: 120
			[83] NoteOff channel: 0 key: 4
			[41] NoteOn channel: 0 key: 16 velocity: 100
			[83] NoteOff channel: 0 key: 16
			`,
		},
	}

	for i, test := range tests {

		if i > 3 {
			continue
		}

		//fmt.Printf("running test [%v]\n", i)
		a = newArpTester()
		a.Run()

		//fmt.Println("let it start (400ms)")
		time.Sleep(100 * time.Millisecond)

		test.fn()

		//fmt.Println("waiting for result (2400ms)")
		//fmt.Printf("after test %v\n", i)
		//time.Sleep(2400 * time.Millisecond)
		//time.Sleep(1000 * time.Millisecond)

		got := a.Result()
		//fmt.Printf("got result %q\n", got)

		//fmt.Println("before close")
		a.Close()
		//fmt.Println("after close")

		if got != test.expected {
			t.Errorf("[%v] %q\ngot:\n%q\n\nexpected:\n%q", i, test.descr, got, test.expected)
		}

		time.Sleep(100 * time.Millisecond)
	}

}
