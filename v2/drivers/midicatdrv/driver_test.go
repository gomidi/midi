package midicatdrv

import (
	"testing"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	"gitlab.com/gomidi/midi/v2/drivers/internal/drivertest"
)

func runTest(t *testing.T, fn func(*testing.T, drivers.In, drivers.Out)) func(*testing.T) {
	return func(*testing.T) {
		drv, err := New()
		if err != nil {
			t.Fatalf("ERROR: %s", err.Error())
		}

		in, err := midi.FindInPort("Midi Through Port-0")
		if err != nil {
			//t.Skipf("could not find in port Midi Through Port-0")
			return
		}
		out, err := midi.FindOutPort("Midi Through Port-0")
		if err != nil {
			//t.Skipf("could not find out port Midi Through Port-0")
			return
		}

		fn(t, in, out)
		drv.Close()
	}
}

func TestSpec(t *testing.T) {
	drv, err := New()
	if err != nil {
		t.Fatalf("ERROR: %s", err.Error())
	}

	drivertest.DriverInterfaceImplementationTest(t, drv)
	drivertest.AutoregisterTest(t, drv)

	tests := []struct {
		name string
		fn   func(*testing.T, drivers.In, drivers.Out)
	}{
		{
			"RunningStatus",
			drivertest.RunningStatusTest,
		},
		{
			"FullStatus",
			drivertest.FullStatusTest,
		},
		{
			"NoActiveSense",
			drivertest.NoActiveSenseTest,
		},
		{
			"NoTimeCode",
			drivertest.NoTimeCodeTest,
		},
		{
			"Sysex",
			drivertest.SysexTest,
		},
		{
			"NoSysex",
			drivertest.NoSysexTest,
		},
	}

	for _, test := range tests {
		f := runTest(t, test.fn)
		t.Run(test.name, f)
	}

}
