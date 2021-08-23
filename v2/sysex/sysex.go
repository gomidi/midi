package sysex

import (
	"bytes"
	"fmt"
)

// see https://www.2writers.com/eddie/TutSysEx.htm
type Manufacturer struct {
	ManufacturerID ManufacturerID
	DeviceID       byte
	ModelID        byte
	InfoRequest    bool // true: requesting infos, false: sending infos
	Address        [3]byte
	SendingData    []byte
	NumReqBytes    [3]byte
}

var GMReset = Manufacturer{
	//F0   41   10   42   12   40007F   00   41   F7
	ManufacturerID: 0x41,
	DeviceID:       0x10,
	ModelID:        0x42,
	InfoRequest:    false,
	Address:        [3]byte{0x40, 0x00, 0x7F},
	SendingData:    []byte{0x00},
}

func init() {
	bt := fmt.Sprintf("% X", GMReset.SysEx())
	if bt != "F0 41 10 42 12 40 00 7F 00 41 F7" {
		panic(bt)
	}
}

func Parse(bt []byte) (*Manufacturer, error) {
	if len(bt) < 11 {
		return nil, fmt.Errorf("sysex message too short (must be 11 bytes minimum")
	}

	if bt[0] != 0xF0 {
		return nil, fmt.Errorf("missing start byte 0xF0")
	}

	var s Manufacturer

	s.ManufacturerID = ManufacturerID(bt[1])
	s.DeviceID = bt[2]
	s.ModelID = bt[3]
	switch bt[4] {
	case 0x11:
		s.InfoRequest = true
	case 0x12:
		s.InfoRequest = false
	default:
		return nil, fmt.Errorf("invalid send/req byte")
	}

	s.Address[0] = bt[5]
	s.Address[1] = bt[6]
	s.Address[2] = bt[7]

	if s.InfoRequest {
		if len(bt) < 13 {
			return nil, fmt.Errorf("sysex message for requesting data too short (must be 13 bytes minimum")
		}
		s.NumReqBytes[0] = bt[8]
		s.NumReqBytes[1] = bt[9]
		s.NumReqBytes[2] = bt[10]
	} else {
		s.SendingData = bt[8 : len(bt)-1]
	}

	checksum := bt[len(bt)-2]

	if checksum != s.Checksum() {
		return nil, fmt.Errorf("invalid checksum")
	}

	if bt[len(bt)-1] != 0xF7 {
		return nil, fmt.Errorf("missing end byte 0xF7")
	}

	return &s, nil
}

func (s Manufacturer) Checksum() (sum byte) {

	/*
				1. Convert hex to decimal:
				   40h = 64
				   11h = 17
				   00h = 0
				   41h = 65
				   63h = 99

				2. Add values:
				   64 + 17 + 0 + 65 + 99 = 245

				3. Divide by 128
				   245 / 128 = 1 remainder 117

				4. Subtract remainder from 128, if remainder is not 0*
				   128 - 117 = 11

				5. Covert to hex:
				   11 = 0Bh

		        *If the remainder is 0 then the checksum is 0.
	*/

	var bt = []byte{s.Address[0], s.Address[1], s.Address[2]}

	if s.InfoRequest {
		bt = append(bt, s.NumReqBytes[0], s.NumReqBytes[1], s.NumReqBytes[2])
	} else {
		bt = append(bt, s.SendingData...)
	}

	var su int32

	for _, b := range bt {
		su += int32(b)
	}

	rem := su % 128

	if rem == 0 {
		return 0
	}

	return byte(128 - rem)
}

func (s Manufacturer) SysEx() []byte {
	var bf bytes.Buffer

	bf.WriteByte(0xF0)
	bf.WriteByte(byte(s.ManufacturerID))
	bf.WriteByte(s.DeviceID)
	bf.WriteByte(s.ModelID)
	if s.InfoRequest {
		bf.WriteByte(0x11)
	} else {
		bf.WriteByte(0x12)
	}
	bf.WriteByte(s.Address[0])
	bf.WriteByte(s.Address[1])
	bf.WriteByte(s.Address[2])

	if s.InfoRequest {
		bf.WriteByte(s.NumReqBytes[0])
		bf.WriteByte(s.NumReqBytes[1])
		bf.WriteByte(s.NumReqBytes[2])
	} else {
		bf.Write(s.SendingData)
	}

	bf.WriteByte(s.Checksum())

	bf.WriteByte(0xF7)

	return bf.Bytes()
}
