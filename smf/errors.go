package smf

import "errors"

var (
	ErrUnsupportedSMFFormat = errors.New("The SMF format was not expected.")
	ErrExpectedMthd         = errors.New("Expected SMF Midi header.")
	ErrBadSizeChunk         = errors.New("Chunk was an unexpected size.")
)
