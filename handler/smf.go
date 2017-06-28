// ReadSMFFile open, reads and closes a complete SMF file.
// If the read content was a valid midi file, nil is returned.
//
package handler

import (
	"io"

	"github.com/gomidi/midi/smf"
	"github.com/gomidi/midi/smf/smfreader"
)

// The messages are dispatched to the corresponding attached functions of the handler.
//
// They must be attached before Handler.ReadSMF is called
// and they must not be unset or replaced until ReadSMF returns.
//
// The *Pos parameter that is passed to the functions is always, because we are reading a file.
func (h *Handler) ReadSMFFile(file string, options ...smfreader.Option) error {
	h.errSMF = nil
	h.pos = &SMFPosition{}
	err := smfreader.ReadFile(file, h.readSMF, options...)
	if err != nil {
		return err
	}
	return h.errSMF
}

// ReadSMF reads midi messages from src (which is supposed to be the content of a standard midi file (SMF))
// until an error or io.EOF happens.
//
// ReadSMF does not close the src.
//
// If the read content was a valid midi file, nil is returned.
//
// The messages are dispatched to the corresponding attached functions of the handler.
//
// They must be attached before Handler.ReadSMF is called
// and they must not be unset or replaced until ReadSMF returns.
//
// The *Pos parameter that is passed to the functions is always, because we are reading a file.
func (h *Handler) ReadSMF(src io.Reader, options ...smfreader.Option) error {
	h.errSMF = nil
	h.pos = &SMFPosition{}
	rd := smfreader.New(src, options...)
	h.readSMF(rd)
	return h.errSMF
}

func (h *Handler) readSMF(rd smf.Reader) {
	hd, err := rd.ReadHeader()

	if err != nil {
		h.errSMF = err
		return
	}

	if h.SMFHeader != nil {
		h.SMFHeader(hd)
	}

	// use err here
	err = h.read(rd)
	if err != io.EOF {
		h.errSMF = err
	}

	return
}
