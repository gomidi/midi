package handler

import (
	"io"

	"github.com/gomidi/midi/live/midireader"
	"github.com/gomidi/midi/messages/realtime"
)

// ReadLive reads midi messages from src until an error or io.EOF happens.
//
// If io.EOF happend the returned error is nil.
//
// ReadLive does not close the src.
//
// The messages are dispatched to the corresponding attached functions of the handler.
//
// They must be attached before Handler.ReadLive is called
// and they must not be unset or replaced until ReadLive returns.
//
// The *Pos parameter that is passed to the functions is nil, because we are in a live setting.
func (h *Handler) ReadLive(src io.Reader, options ...midireader.Option) (err error) {
	rthandler := func(m realtime.Message) {
		switch m {
		// ticks (most important, must be sent every 10 milliseconds) comes first
		case realtime.Tick:
			if h.Tick != nil {
				h.Tick()
			}

		// clock a bit slower synchronization method (24 MIDI Clocks in every quarter note) comes next
		case realtime.TimingClock:
			if h.Clock != nil {
				h.Clock()
			}

		// ok starting and continuing should not take too lpng
		case realtime.Start:
			if h.Start != nil {
				h.Start()
			}

		case realtime.Continue:
			if h.Continue != nil {
				h.Continue()
			}

		// Active Sense must come every 300 milliseconds (but is seldom implemented)
		case realtime.ActiveSensing:
			if h.ActiveSense != nil {
				h.ActiveSense()
			}

		// put any user defined realtime message here
		case realtime.Undefined4:
			if h.UndefinedRealtime4 != nil {
				h.UndefinedRealtime4()
			}

		// ok, stopping is not so urgent
		case realtime.Stop:
			if h.Stop != nil {
				h.Stop()
			}

		// reset may take some time
		case realtime.Reset:
			if h.Reset != nil {
				h.Reset()
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
