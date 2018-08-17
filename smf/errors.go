package smf

import "errors"

// ErrFinished indicates that the read or write operation has been finished successfully
var ErrFinished = errors.New("SMF action finished successfully")
