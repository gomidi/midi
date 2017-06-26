package runningstatus

import "io"

type Reader interface {
	Read(canary byte) (status byte, changed bool)
}

type reader struct {
	status byte
}

func (r *reader) read(canary byte) (status byte, changed bool) {

	// channel/Voice Category Status
	if canary >= 0x80 && canary <= 0xEF {
		r.status = canary
		changed = true
	}

	return r.status, changed
}

type livereader struct {
	reader
}

/*
   his (http://midi.teragonaudio.com/tech/midispec.htm) take on running status buffer
   A recommended approach for a receiving device is to maintain its "running status buffer" as so:

       Buffer is cleared (ie, set to 0) at power up.
       Buffer stores the status when a Voice Category Status (ie, 0x80 to 0xEF) is received.
       Buffer is cleared when a System Common Category Status (ie, 0xF0 to 0xF7) is received.
       Nothing is done to the buffer when a RealTime Category message is received.
       Any data bytes are ignored when the buffer is 0. (I think that only holds for realtime midi)
*/
func (r *livereader) Read(canary byte) (status byte, changed bool) {

	// here we clear for System Common Category messages
	if canary >= 0xF0 && canary <= 0xF7 {
		r.status = 0
		return r.status, true
	}

	return r.read(canary)
}

type smfreader struct {
	reader
}

func (r *smfreader) Read(canary byte) (status byte, changed bool) {

	// here we clear for meta messages
	if canary == 0xFF {
		r.status = 0
		return r.status, true
	}

	return r.read(canary)
}

func NewLiveReader() Reader {
	return &livereader{}
}

func NewSMFReader() Reader {
	return &smfreader{}
}

func hasBitU8(n uint8, pos uint8) bool {
	val := n & (1 << pos)
	return (val > 0)
}

func IsStatusByte(b uint8) bool {
	return hasBitU8(b, 7)
}

type Writer interface {
	io.Writer
	runningstatus()
}

func NewSMFWriter(output io.Writer) Writer {
	return &smfwriter{output, 0}
}

func NewLiveWriter(output io.Writer) Writer {
	return &liveWriter{&smfwriter{output, 0}}
}

type smfwriter struct {
	output io.Writer
	status byte
}

func (w *smfwriter) runningstatus() {}

func (w *smfwriter) write(b []byte) (n int, err error) {
	return w.output.Write(b)
}

func (w *smfwriter) Write(msg []byte) (int, error) {

	// for non channel messages, reset status and write whole message
	if !IsStatusByte(msg[0]) {
		w.status = 0
		return w.write(msg)
	}

	// for a different status, store runningStatus and write whole message
	if msg[0] != w.status {
		w.status = msg[0]
		return w.write(msg)
	}

	// we got the same status as runningStatus, so omit the status byte when writing
	return w.write(msg[1:])
}

type liveWriter struct {
	*smfwriter
}

func (w *liveWriter) Write(msg []byte) (int, error) {
	// for realtime system messages, don't affect status and write the whole message
	if msg[0] > 0xF7 {
		return w.smfwriter.write(msg)
	}

	return w.smfwriter.Write(msg)
}
