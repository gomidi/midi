package midi

import (
	"gitlab.com/gomidi/midi/v2/drivers"
)

// FindInPort returns the number of the midi in port with the given name
// It returns -1, if the port can't be found.
func FindInPort(name string) int {
	in, err := drivers.InByName(name)
	if err != nil {
		return -1
	}
	return in.Number()
}

// FindOutPort returns the number of the midi out port with the given name
// It returns -1, if the port can't be found.
func FindOutPort(name string) int {
	out, err := drivers.OutByName(name)
	if err != nil {
		return -1
	}
	return out.Number()
}
