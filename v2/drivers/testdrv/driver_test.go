package testdrv

import (
	"testing"

	"gitlab.com/gomidi/midi/v2/drivers"
	"gitlab.com/gomidi/midi/v2/drivers/drivertest"
)

func runTest(t *testing.T, fn func(*testing.T, drivers.In, drivers.Out)) func(*testing.T) {
	return func(*testing.T) {
		drv := New("testdrv")
		ins, _ := drv.Ins()
		outs, _ := drv.Outs()
		fn(t, ins[0], outs[0])
		drv.Close()
	}
}

func TestSpec(t *testing.T) {
	drv := New("testdrv")

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
