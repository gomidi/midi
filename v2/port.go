package midi

import (
	"gitlab.com/gomidi/midi/v2/drivers"
)

// FindInPort returns the midi in port that contains the given name
// and an error, if the port can't be found.
func FindInPort(name string) (drivers.In, error) {
	in, err := drivers.InByName(name)
	if err != nil {
		return nil, err
	}
	in.Close()
	return in, nil
}

// OutPort returns the midi out port for the given port number
func OutPort(portnumber int) (drivers.Out, error) {
	return drivers.OutByNumber(portnumber)
}

// InPort returns the midi in port for the given port number
func InPort(portnumber int) (drivers.In, error) {
	return drivers.InByNumber(portnumber)
}

// FindOutPort returns the midi out port that contains the given name
// and an error, if the port can't be found.
func FindOutPort(name string) (drivers.Out, error) {
	out, err := drivers.OutByName(name)
	if err != nil {
		return nil, err
	}
	out.Close()
	return out, nil
}
