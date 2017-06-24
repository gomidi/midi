package smfreader

import "errors"

var (
	ErrUnsupportedSMFFormat  = errors.New("The SMF format was not expected.")
	ErrExpectedMthd          = errors.New("Expected SMF Midi header.")
	ErrBadSizeChunk          = errors.New("Chunk was an unexpected size.")
	ErrInterruptedByCallback = errors.New("interrupted by callback")
)
