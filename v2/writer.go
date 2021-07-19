package midi

import (
	"gitlab.com/gomidi/midi/v2/drivers"
)

func newWriter(out drivers.Out) *writer {
	wr := &writer{
		out: out,
	}

	if syss, canDo := out.(drivers.SysExSender); canDo {
		wr.sysexOut = syss
	}

	if rts, canDo := out.(drivers.RealtimeSender); canDo {
		wr.realtimeOut = rts
	}

	return wr
}

type writer struct {
	out         drivers.Out
	sysexOut    drivers.SysExSender
	realtimeOut drivers.RealtimeSender
	//mx          sync.RWMutex
	//inSysEx        bool
	//interruptSysEx chan bool
}

/*
func (w *writer) isInSysEx() bool {
	w.mx.RLock()
	defer w.mx.RUnlock()
	return w.inSysEx
}
*/

func (w *writer) sendSysExToDriver(bt []byte) {
	if w.sysexOut != nil {
		w.sysexOut.SendSysEx(bt)
	} else {
		w.out.Send(bt)
	}
}

/*
// TODO have a compile switch to enable / disable the async sysex writing, since e.g. tinygo does not
// support goroutines on all platforms
func (w *writer) listenForSysEx(ch chan byte, stop chan bool) {
	var bf bytes.Buffer
	for {
		select {
		case b := <-ch:
			if w.isInSysEx() {
				if b == 0xF0 {
					// starts a new sysex
					bf.Reset()
				}
				bf.WriteByte(b)
				if b == 0xF7 {
					w.sendSysExToDriver(bf.Bytes())
					stop <- true
					return
				}
			} else {
				stop <- true
				return
			}
		}
	}
}

// TODO have a compile switch to enable / disable the async sysex writing, since e.g. tinygo does not
// support goroutines on all platforms
func (w *writer) sendSysEx(bt []byte) {
	ch := make(chan byte, len(bt))
	stop := make(chan bool)
	go w.listenForSysEx(ch, stop)
	for _, b := range bt {
		ch <- b
	}
	<-stop
	w.mx.Lock()
	w.inSysEx = false
	w.mx.Unlock()
}
*/

func (w *writer) sendRT(bt []byte) error {
	if w.realtimeOut != nil {
		return w.realtimeOut.SendRealtime(bt[0])
	} else {
		return w.out.Send(bt)
	}
}

func (w *writer) Send(msg Message) error {
	switch {
	case msg.Is(RealTimeMsg):
		return w.sendRT(msg.Data)
	case msg.Is(SysExMsg):
		//w.mx.Lock()
		//w.inSysEx = true
		w.sendSysExToDriver(msg.Data)
		//w.mx.Unlock()
		//go w.sendSysEx(msg.Data)
		return nil
	default:
		//w.mx.Lock()
		err := w.out.Send(msg.Data)
		/*
			if w.inSysEx {
				w.inSysEx = false
			}
		*/
		//w.mx.Unlock()
		return err
	}
}

var _ Sender = &writer{}

func SenderToPort(portnumber int) (Sender, error) {
	out, err := drivers.OutByNumber(portnumber)
	if err != nil {
		return nil, err
	}
	return newWriter(out), nil
}
