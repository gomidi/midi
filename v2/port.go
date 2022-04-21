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
	in.Close()
	return in.Number()
}

// CloseInPort closes the in port
func CloseInPort(num int) error {
	in, err := drivers.InByNumber(num)
	if err != nil {
		return err
	}
	in.Close()
	return nil
}

// FindOutPort returns the number of the midi out port with the given name
// It returns -1, if the port can't be found.
func FindOutPort(name string) int {
	out, err := drivers.OutByName(name)
	if err != nil {
		return -1
	}
	out.Close()
	return out.Number()
}

// CloseOutPort closes the out port
func CloseOutPort(num int) error {
	out, err := drivers.OutByNumber(num)
	if err != nil {
		return err
	}
	out.Close()
	return nil
}
