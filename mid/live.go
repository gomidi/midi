package mid

import (
	"io"

	"github.com/gomidi/midi/midimessage/realtime"
	"github.com/gomidi/midi/midireader"
)

// ReadLive reads midi messages from src until an error or io.EOF happens.
//
// If io.EOF happened the returned error is nil.
//
// ReadLive does not close the src.
//
// The messages are dispatched to the corresponding attached functions of the handler.
//
// They must be attached before Handler.ReadLive is called
// and they must not be unset or replaced until ReadLive returns.
func (h *Handler) ReadLive(src io.Reader, options ...midireader.Option) (err error) {
	h.pos = nil
	rthandler := func(m realtime.Message) {
		switch m {
		// ticks (most important, must be sent every 10 milliseconds) comes first
		case realtime.Tick:
			if h.Message.Realtime.Tick != nil {
				h.Message.Realtime.Tick()
			}

		// clock a bit slower synchronization method (24 MIDI Clocks in every quarter note) comes next
		case realtime.TimingClock:
			if h.Message.Realtime.Clock != nil {
				h.Message.Realtime.Clock()
			}

		// ok starting and continuing should not take too lpng
		case realtime.Start:
			if h.Message.Realtime.Start != nil {
				h.Message.Realtime.Start()
			}

		case realtime.Continue:
			if h.Message.Realtime.Continue != nil {
				h.Message.Realtime.Continue()
			}

		// Active Sense must come every 300 milliseconds (but is seldom implemented)
		case realtime.ActiveSensing:
			if h.Message.Realtime.ActiveSense != nil {
				h.Message.Realtime.ActiveSense()
			}

		// put any user defined realtime message here
		case realtime.Undefined4:
			if h.Message.Unknown != nil {
				h.Message.Unknown(h.pos, m)
			}

		// ok, stopping is not so urgent
		case realtime.Stop:
			if h.Message.Realtime.Stop != nil {
				h.Message.Realtime.Stop()
			}

		// reset may take some time
		case realtime.Reset:
			if h.Message.Realtime.Reset != nil {
				h.Message.Realtime.Reset()
			}

		}
	}

	rd := midireader.New(src, rthandler, options...)
	err = h.read(rd)

	if err == io.EOF {
		return nil
	}

	return
}
