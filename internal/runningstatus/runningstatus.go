package runningstatus

import (
	"github.com/gomidi/midi"
	"github.com/gomidi/midi/internal/midilib"

	"io"
)

// Reader is a running status reader
type Reader interface {
	// Read reads the status byte off the canary and returns
	// if it has changed compared to the previous read
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

// Read reads the status byte from the given canary, while respecting
// running status and returns whether the status has changed
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

// Read reads the status byte from the given canary, while respecting
// running status and returns whether the status has changed
func (r *smfreader) Read(canary byte) (status byte, changed bool) {

	// here we clear for meta messages
	if canary == 0xFF || canary == 0xF0 || canary == 0xF7 {
		r.status = 0
		return r.status, true
	}

	return r.read(canary)
}

// NewLiveReader returns a new Reader for reading of live MIDI data
func NewLiveReader() Reader {
	return &livereader{}
}

// NewSMFReader returns a new Reader for reading of SMF MIDI data
func NewSMFReader() Reader {
	return &smfreader{}
}

// Writer writes messages with running status byte
type Writer interface {
	io.Writer
	runningstatus()
}

// NewSMFWriter returns a new SMFWriter
func NewSMFWriter() SMFWriter {
	return &smfwriter{0}
}

// SMFWriter is a writer for writing messages with running status byte in SMF files
type SMFWriter interface {
	Write(midi.Message) []byte
	ResetStatus()
}

// NewLiveWriter returns a new Writer for live writing of messages with running status byte
func NewLiveWriter(output io.Writer) Writer {
	return &liveWriter{output, 0}
}

type smfwriter struct {
	status byte
}

func (w *smfwriter) ResetStatus() {
	w.status = 0
}

// Write writes the given message with running status
func (w *smfwriter) Write(msg midi.Message) []byte {
	raw := msg.Raw()
	// fmt.Printf("should write %s (% X)\n", msg, raw)
	firstByte := raw[0]

	// for non channel messages, reset status and write whole message
	if !midilib.IsChannelMessage(firstByte) {
		//fmt.Printf("is no channel message, resetting status\n")
		w.status = 0
		return raw
	}

	// for a different status, store runningStatus and write whole message
	if firstByte != w.status {
		// fmt.Printf("setting status to: % X (was: % X)\n", firstByte, w.status)
		w.status = firstByte
		return raw
	}

	// we got the same status as runningStatus, so omit the status byte when writing
	// fmt.Printf("taking running status (% X), writing: % X\n", w.status, raw[1:])
	return raw[1:]
}

func (w *liveWriter) runningstatus() {

}

func (w *liveWriter) write(b []byte) (n int, err error) {
	return w.output.Write(b)
}

type liveWriter struct {
	output io.Writer
	status byte
}

// Write writes the given message with running status
func (w *liveWriter) Write(msg []byte) (int, error) {
	// fmt.Printf("should write % X\n", msg)
	// for realtime system messages, don't affect status and write the whole message
	if msg[0] > 0xF7 {
		return w.write(msg)
	}

	// for non channel messages, reset status and write whole message
	if !midilib.IsChannelMessage(msg[0]) {
		// fmt.Printf("is no channel message, resetting status\n")
		w.status = 0
		return w.write(msg)
	}

	// for a different status, store runningStatus and write whole message
	if msg[0] != w.status {
		// fmt.Printf("setting status to: % X (was: % X)\n", msg[0], w.status)
		w.status = msg[0]
		return w.write(msg)
	}

	// we got the same status as runningStatus, so omit the status byte when writing
	// fmt.Printf("taking running status (% X), writing: % X\n", w.status, msg[1:])
	return w.write(msg[1:])
}
