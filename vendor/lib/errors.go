package lib

import "errors"

func UnexpectedEventLengthError(s string) error {
	return errors.New(s)
}

var ErrUnexpectedEOF = errors.New("Unexpected End of File found.")
