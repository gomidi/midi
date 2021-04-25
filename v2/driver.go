package midi

var (
	firstDriver string

	// DRIVERS is the registry for MIDI drivers
	DRIVERS = map[string]Driver{}
)

// RegisterDriver register a driver
func RegisterDriver(d Driver) {
	if len(DRIVERS) == 0 {
		firstDriver = d.String()
	}
	DRIVERS[d.String()] = d
}

// GetDriver returns the first available driver
func GetDriver() Driver {
	if len(DRIVERS) == 0 {
		return nil
	}
	return DRIVERS[firstDriver]
}

// CloseDriver closes the first available driver
func CloseDriver() {
	d := GetDriver()
	if d != nil {
		d.Close()
	}
}

// Driver is a driver for MIDI connections.
// Apart from system exclusive data the MIDI bytes must be provided in complete messages.
// However for channel messages the status from a `running status`.
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
