package midi

var DRIVERS = map[string]Driver{}
var firstDriver string

// RegisterDriver register a driver
func RegisterDriver(d Driver) {
	DRIVERS[d.String()] = d
	if len(DRIVERS) == 0 {
		firstDriver = d.String()
	}
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

