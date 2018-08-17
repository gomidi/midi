package smfreader

import "errors"

var (
	errUnsupportedSMFFormat  = errors.New("The SMF format was not expected.")
	errExpectedMthd          = errors.New("Expected SMF Midi header.")
	errBadSizeChunk          = errors.New("Chunk was an unexpected size.")
	errInterruptedByCallback = errors.New("interrupted by callback")
	// ErrFinished is the error returned, if the SMF has been read
	ErrFinished = errors.New("SMF has been read")
	// ErrMissing is the error returned, if there is no more data, but tracks are missing
	ErrMissing = errors.New("incomplete, tracks missing")
)
