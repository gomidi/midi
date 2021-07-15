package drivers

var (
	firstDriver string

	// REGISTRY is the registry for MIDI drivers
	REGISTRY = map[string]Driver{}
)

// RegisterDriver register a driver
func Register(d Driver) {
	if len(REGISTRY) == 0 {
		firstDriver = d.String()
	}
	REGISTRY[d.String()] = d
}

// Get returns the first available driver
func Get() Driver {
	if len(REGISTRY) == 0 {
		return nil
	}
	return REGISTRY[firstDriver]
}

// Close closes the first available driver
func Close() {
	d := Get()
	if d != nil {
		d.Close()
	}
}

// Driver is a driver for MIDI connections.
// It may provide the timing delta to the previous message in micro seconds.
// It must send the given MIDI data immediately.
type Driver interface {

	// Ins returns the available MIDI input ports.
	Ins() ([]In, error)

	// Outs returns the available MIDI output ports.
	Outs() ([]Out, error)

	// String returns the name of the driver.
	String() string

	// Close closes the driver. Must be called for cleanup at the end of a session.
	Close() error
}
