package lib

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

const limit = byte('\n')

func Read(rd io.Reader) (out []byte, err error) {
	for {
		b, errRd := read(rd)
		if errRd != nil {
			return nil, errRd
		}

		if b == limit {
			return out, err
		}
		out = append(out, b)
	}
}

func ReadAndConvert(rd io.Reader) (out []byte, err error) {
	out, err = Read(rd)
	if err != nil {
		return
	}

	return convert(out)
}
