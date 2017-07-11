package midiwriter

type config struct {
	noRunningStatus bool
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
