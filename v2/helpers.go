package midi

/*
b, err := midilib.ReadVarLengthData(rd)

	if err != nil {
		return "", err
	}

	return string(b), nil
*/

// bin2decDenom converts the binary denominator to the decimal
func bin2decDenom(bin uint8) uint8 {
	if bin == 0 {
		return 1
	}
	return 2 << (bin - 1)
}
