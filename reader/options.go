package reader

import (
	"fmt"
)

// IgnoreMIDIClock lets the reader not use MIDI clock to calculate the tempo
func IgnoreMIDIClock() func(*Reader) {
	return func(r *Reader) {
		r.ignoreMIDIClock = true
	}
}

// SetLogger allows to set a custom logger for the Reader
func SetLogger(l Logger) func(*Reader) {
	return func(r *Reader) {
		r.logger = l
	}
}

// NoLogger is an option to disable the defaut logging of a Reader
func NoLogger() func(*Reader) {
	return func(r *Reader) {
		r.logger = nil
	}
}

// Logger is the inferface used by Reader for logging incoming messages.
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
