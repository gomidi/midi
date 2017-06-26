package channel

func clearBitU16(n uint16, pos uint16) uint16 {
	mask := ^(uint16(1) << pos)
	n &= mask
	return n
}

func parseStatus(b byte) (messageType uint8, messageChannel uint8) {
	messageType = (b & 0xF0) >> 4
	messageChannel = b & 0x0F
	return
}

func msbLsbSigned(n int16) uint16 {
	if n > 8191 {
		panic("n must not overflow 14bits (max 8191)")
	}
	if n < -8191 {
		panic("n must not overflow 14bits (min -8191)")
	}
	return msbLsbUnsigned(uint16(n + 8192))
}

// takes a 14bit uint and pads it to 16 bit like in the specs for e.g. pitchbend
func msbLsbUnsigned(n uint16) uint16 {
	if n > 16383 {
		panic("n must not overflow 14bits (max 16383)")
	}

	lsb := n << 8
	lsb = clearBitU16(lsb, 15)
	lsb = clearBitU16(lsb, 7)

	// 0x7f = 127 = 0000000001111111
	msb := 0x7f & (n >> 7)
	return lsb | msb
}

func clearBitU8(n uint8, pos uint8) uint8 {
	mask := ^(uint8(1) << pos)
	n &= mask
	return n
}

// parseUint7 parses a 7-bit bit integer from a byte, ignoring the high bit.
func parseUint7(b byte) uint8 {
	return b & 0x7f
}

// parseTwoUint7 parses two 7-bit bit integer stored in two bytes, ignoring the high bit in each.
func parseTwoUint7(b1, b2 byte) (uint8, uint8) {
	return (b1 & 0x7f), (b2 & 0x7f)
}
