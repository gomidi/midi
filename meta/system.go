package meta

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"lib"
)

var (
	_ SysEvent = SysSongPositionPointer(0)
	_ SysEvent = SysSongSelect(0)
	_ SysEvent = SysTuneRequest
	_ SysEvent = sysUnknown(nil)
	_ SysEvent = SysEx(nil)
	_ SysEvent = SysMTCQuarterFrame(0)

/*
	statusSysEx                     = byte(0xF0)
	statusSysMTCQuarterFrameMessage = byte(0xF1)
	statusSongPositionPointer       = byte(0xF2)
	statusSongSelect                = byte(0xF3)
	statusTuneRequest               = byte(0xF6)
*/
)

type SysEvent interface {
	Event
	sysCommon()
}

type SysMTCQuarterFrame uint8

func (m SysMTCQuarterFrame) Number() uint8 {
	return uint8(m)
}

func (m SysMTCQuarterFrame) readFrom(rd io.Reader) (Event, error) {
	// TODO TEST
	length, err := lib.ReadVarLength(rd)

	if err != nil {
		return nil, err
	}

	if length != 2 {
		err = lib.UnexpectedEventLengthError("MTC QuarterFrameMessage expected length 1")
		return nil, err
	}

	var b1 byte

	b1, err = lib.ReadByte(rd)

	if err != nil {
		return nil, err
	}

	return SysMTCQuarterFrame(b1), nil

}
func (m SysMTCQuarterFrame) meta() {}

func (m SysMTCQuarterFrame) sysCommon() {}

func (m SysMTCQuarterFrame) String() string {
	return fmt.Sprintf("%T: %v", m, uint8(m))
}

func (m SysMTCQuarterFrame) Raw() []byte {
	return (&metaEvent{
		typ:  statusSysMTCQuarterFrameMessage,
		data: []byte{byte(m)},
	}).Bytes()

}

func (m SysSongPositionPointer) readFrom(rd io.Reader) (Event, error) {
	// TODO TEST
	length, err := lib.ReadVarLength(rd)

	if err != nil {
		return nil, err
	}

	if length != 2 {
		err = lib.UnexpectedEventLengthError("SongPositionPointer expected length 2")
		return nil, err
	}
	var b1, b2 byte

	b1, err = lib.ReadByte(rd)

	if err != nil {
		return nil, err
	}

	b2, err = lib.ReadByte(rd)

	if err != nil {
		return nil, err
	}

	val := uint16((b2)&0x7f) << 7
	val |= uint16(b1) & 0x7f

	return SysSongPositionPointer(val), nil
}

func (m SysSongPositionPointer) meta() {}

type SysSongPositionPointer uint16

func (m SysSongPositionPointer) Number() uint16 {
	return uint16(m)
}

func (m SysSongPositionPointer) String() string {
	return fmt.Sprintf("%T: %v", m, uint16(m))
}

// TODO test
func (m SysSongPositionPointer) Raw() []byte {
	r := lib.MsbLsbUnsigned(uint16(m))

	var bf bytes.Buffer
	//	binary.Write(&bf, binary.BigEndian, uint16(change))
	binary.Write(&bf, binary.BigEndian, r)
	return (&metaEvent{
		typ:  statusSongPositionPointer,
		data: bf.Bytes(),
	}).Bytes()

}
func (m SysSongPositionPointer) sysCommon() {}

func (m SysSongSelect) Pos() uint16 {
	return uint16(m)
}

// TODO Test
func (m SysSongSelect) Raw() []byte {
	return (&metaEvent{
		typ:  statusSongSelect,
		data: []byte{byte(m)},
	}).Bytes()
}

type SysSongSelect uint8

func (m SysSongSelect) Number() uint8 {
	return uint8(m)
}

func (m SysSongSelect) String() string {
	return fmt.Sprintf("%T: %v", m, uint8(m))
}

func (m SysSongSelect) sysCommon() {}

func (m SysSongSelect) meta() {}

func (m SysSongSelect) readFrom(rd io.Reader) (Event, error) {
	// TODO TEST
	length, err := lib.ReadVarLength(rd)

	if err != nil {
		return nil, err
	}

	if length != 2 {
		err = lib.UnexpectedEventLengthError("SongSelect expected length 1")
		return nil, err
	}

	var b1 byte

	b1, err = lib.ReadByte(rd)

	if err != nil {
		return nil, err
	}

	return SysSongSelect(b1), nil

}

type tuneRequest bool

const (
	SysTuneRequest = tuneRequest(false)
)

func (m tuneRequest) meta() {}

func (m tuneRequest) String() string {
	return fmt.Sprintf("%T", m)
}

func (m tuneRequest) readFrom(rd io.Reader) (Event, error) {
	return m, nil
}

func (m tuneRequest) sysCommon() {}

// TODO test
func (m tuneRequest) Raw() []byte {
	return (&metaEvent{
		typ:  statusTuneRequest,
		data: []byte{},
	}).Bytes()
}

type sysUnknown []byte

func (m sysUnknown) String() string {
	return fmt.Sprintf("%T len: %v", m, len(m))
}

func (m sysUnknown) Len() int {
	return len(m)
}

func (m sysUnknown) meta() {}

func (m sysUnknown) Bytes() []byte {
	return []byte(m)
}

// TODO: don't know if I should implement
func (m sysUnknown) readFrom(rd io.Reader) (Event, error) {
	return m, nil
}

func (m sysUnknown) sysCommon() {}

func (m sysUnknown) Raw() []byte {
	panic("not implemented")
}

type SysEx []byte

func (m SysEx) Bytes() []byte {
	return []byte(m)
}

// TODO: implement
func (m SysEx) readFrom(rd io.Reader) (Event, error) {
	return m, nil
}

func (m SysEx) meta() {}

func (m SysEx) String() string {
	return fmt.Sprintf("%T len: %v", m, len(m))
}

func (m SysEx) Len() int {
	return len(m)
}

func (m SysEx) sysCommon() {}

func (m SysEx) Raw() []byte {
	var b = []byte{0xF0}
	b = append(b, []byte(m)...)
	b = append(b, 0xF7)
	return b
}
