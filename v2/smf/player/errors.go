package player

import "errors"

var (
	ErrNoSMFData  = errors.New("no SMF data is set")
	ErrIsPlaying  = errors.New("already playing")
	ErrIsStopped  = errors.New("already stopped")
	ErrInvalidSMF = errors.New("failed to read tracks from SMF data")
	ErrNoOutPort  = errors.New("failed to get sound out port")
	errDone       = errors.New("done playing SMF data")
	errStopped    = errors.New("stopped playing")
	errPaused     = errors.New("paused playing")
)
