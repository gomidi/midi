package midiwriter

type config struct {
	noRunningStatus        bool
	checkMessageType       bool
	ignoreWrongMessageType bool
}

// Option is a configuration option for a writer
type Option func(*config)

// NoRunningStatus is an option for the writer that prevents it from
// using the running status.
func NoRunningStatus() Option {
	return func(c *config) {
		c.noRunningStatus = true
	}
}

// CheckMessageType is an option for Writer to check, if the message can
// be send as live MIDI data to an instrument. If not, an error will be returned.
func CheckMessageType() Option {
	return func(c *config) {
		c.checkMessageType = true
	}
}

// SkipNonLiveMessages skips MIDI messages that can't be send over the wire (live).
func SkipNonLiveMessages() Option {
	return func(c *config) {
		c.ignoreWrongMessageType = true
	}
}
