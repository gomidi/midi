package reader

import (
	"fmt"

	"gitlab.com/gomidi/midi/midireader"
)

// IgnoreMIDIClock lets the reader not use MIDI clock to calculate the tempo
func IgnoreMIDIClock() Option {
	return func(r *Reader) {
		r.ignoreMIDIClock = true
	}
}

// ReaderOption configures the reader
type Option func(*Reader)

// SetLogger allows to set a custom logger for the Reader
func SetLogger(l Logger) Option {
	return func(r *Reader) {
		r.logger = l
	}
}

func Options(options ...midireader.Option) Option {
	return func(r *Reader) {
		r.midiReaderOptions = append(r.midiReaderOptions, options...)
	}
}

// NoLogger is an option to disable the defaut logging of a Reader
func NoLogger() Option {
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
