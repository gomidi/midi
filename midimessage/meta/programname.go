package meta

import (
	"fmt"
	"io"
)

type ProgramName string

/*
ProgramName

FF 08 length text

This optional event is used to embed the patch/program name that is called up by the immediately subsequent Bank Select and Program Change messages. It serves to aid the end user in making an intelligent program choice when using different hardware.

This event may appear anywhere in a track, and there may be multiple occurrences within a track.
*/
func (m ProgramName) String() string {
	return fmt.Sprintf("%T: %#v", m, m.Text())
}

func (m ProgramName) Text() string {
	return string(m)
}

func (m ProgramName) Raw() []byte {
	return (&metaMessage{
		Typ:  byte(byteProgramName),
		Data: []byte(m),
	}).Bytes()
}

func (m ProgramName) readFrom(rd io.Reader) (Message, error) {
	text, err := readText(rd)

	if err != nil {
		return nil, err
	}

	return ProgramName(text), nil
}

func (m ProgramName) meta() {}
