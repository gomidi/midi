package handler

import "fmt"

// Option configures the handler
type Option func(*Handler)

// SetLogger allows to set a custom logger for the handler
func SetLogger(l Logger) Option {
	return func(h *Handler) {
		h.logger = l
	}
}

// NoLogger is an option to disable the defaut logging of a handler
func NoLogger() Option {
	return func(h *Handler) {
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
