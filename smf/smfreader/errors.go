package smfreader

import "errors"

var (
	errUnsupportedSMFFormat  = errors.New("The SMF format was not expected.")
	errExpectedMthd          = errors.New("Expected SMF Midi header.")
	errBadSizeChunk          = errors.New("Chunk was an unexpected size.")
	errInterruptedByCallback = errors.New("interrupted by callback")
	ErrFinished              = errors.New("SMF has been read")
	ErrMissing               = errors.New("incomplete, tracks missing")
)
