package smfwriter

import "errors"

var (
	// ErrFinished is the error returned, if the last track has been written
	ErrFinished = errors.New("SMF has been written")
)
