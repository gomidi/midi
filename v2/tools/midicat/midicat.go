package midicat

import (
	"fmt"
	"io"
)

func read(rd io.Reader) (byte, error) {
	var b = make([]byte, 1)

	i, err := rd.Read(b)

	if err != nil {
		return 0, err
	}

	if i != 1 {
		return 0, err
	}

	return b[0], nil
}

func convert(b []byte) (out []byte, err error) {
	out = make([]byte, len(b)/2)

	_, err = fmt.Sscanf(string(b), "%X", &out)
	if err != nil {
		return nil, err
	}

	return out, nil

}

func convertDelta(b []byte) (deltams int32, err error) {
	_, err = fmt.Sscanf(string(b), "%d", &deltams)
	if err != nil {
		return -1, err
	}

	return deltams, nil

}

const limit = byte('\n')

func Read(rd io.Reader) (out []byte, deltams int32, err error) {
	var deltaRead bool
	var deltaBf []byte

	for {
		b, errRd := read(rd)
		if errRd != nil {
			return nil, -1, errRd
		}

		if b == ' ' {
			deltams, err = convertDelta(deltaBf)
			if err != nil {
				return
			}
			deltaRead = true
			continue
		}

		if b == limit {
			return out, deltams, err
		}

		if deltaRead {
			out = append(out, b)
		} else {
			deltaBf = append(deltaBf, b)
		}
	}
}

func ReadAndConvert(rd io.Reader) (out []byte, deltams int32, err error) {
	out, deltams, err = Read(rd)
	if err != nil {
		return
	}

	conv, err := convert(out)
	return conv, deltams, err
}
