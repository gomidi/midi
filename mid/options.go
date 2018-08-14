package mid

import "fmt"

// ReaderOption configures the reader
type ReaderOption func(*Reader)

// SetLogger allows to set a custom logger for the handler
func SetLogger(l Logger) ReaderOption {
	return func(h *Reader) {
		h.logger = l
	}
}

// NoLogger is an option to disable the defaut logging of a handler
func NoLogger() ReaderOption {
	return func(h *Reader) {
		h.logger = nil
	}
}

// Logger is the inferface used by Handler for logging incoming messages.
type Logger interface {
	Printf(format string, vals ...interface{})
}

type logfunc func(format string, vals ...interface{})

func (l logfunc) Printf(format string, vals ...interface{}) {
	l(format, vals...)
}

func printf(format string, vals ...interface{}) {
	fmt.Printf(format, vals...)
}
