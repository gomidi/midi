package midi

import (
	"fmt"

	"gitlab.com/gomidi/midi/v2/drivers"
)

func newWriter(out drivers.Out) *writer {
	wr := &writer{
		out: out,
	}

	/*
		if syss, canDo := out.(drivers.SysExSender); canDo {
			wr.sysexOut = syss
		}

		if rts, canDo := out.(drivers.RealtimeSender); canDo {
			wr.realtimeOut = rts
		}
	*/
	return wr
}

type writer struct {
	out drivers.Out
	//sysexOut    drivers.SysExSender
	//realtimeOut drivers.RealtimeSender
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
	w.out.Send(bt)
	/*
		if w.sysexOut != nil {
			w.sysexOut.SendSysEx(bt)
		} else {
			var i int
			l := len(bt)
			for {
				if i >= l {
					return
				}
				packet := [3]byte{bt[i], 0, 0}
				i++
				if i >= l {
					w.out.Send(packet)
					return
				}
				packet[1] = bt[i]
				i++
				if i >= l {
					w.out.Send(packet)
					return
				}
				packet[2] = bt[i]
				w.out.Send(packet)
				i++
			}
		}
	*/
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

func (w *writer) sendRT(b byte) error {
	return w.out.Send([]byte{b})
	/*
		if w.realtimeOut != nil {
			return w.realtimeOut.SendRealtime(b)
		} else {
			return w.out.Send([3]byte{b, 0, 0})
		}
	*/
}

func (w *writer) SendSysEx(data []byte) error {
	w.sendSysExToDriver(data)
	return nil
}

func (w *writer) Send(msg Msg) error {
	switch {
	case msg.Kind() == RealTimeMsg:
		return w.sendRT(msg.Data[0])
		/*
			case msg.Is(SysExMsg):
				//w.mx.Lock()
				//w.inSysEx = true
				w.sendSysExToDriver(msg.Data)
				//w.mx.Unlock()
				//go w.sendSysEx(msg.Data)
				return nil
		*/
	default:
		/*
			//w.mx.Lock()
			var packet [3]byte
			packet[0] = msg.Data[0]
			l := len(msg.Data)
			if l > 1 {
				packet[1] = msg.Data[1]
			}
			if l > 2 {
				packet[2] = msg.Data[2]
			}
			err := w.out.Send(packet)
		*/
		err := w.out.Send(msg.Bytes())
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
		return nil, fmt.Errorf("can't get out port %v: %s\n", portnumber, err)
	}

	//fmt.Printf("returning writer to %q\n", out.String())
	return newWriter(out), nil
}
