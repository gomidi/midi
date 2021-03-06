package smf

/*
const (
	// SMF0 represents the singletrack SMF format (0).
	SMF0 = format(0)

	// SMF1 represents the multitrack SMF format (1).
	SMF1 = format(1)

	// SMF2 represents the sequential track SMF format (2).
	SMF2 = format(2)
)

var (
	_ Format = SMF0
)

// Format is the common interface of all SMF file formats
type Format interface {

	// String returns the string representation of the SMF format.
	String() string

	// Type returns the type of the SMF file: 0 for SMF0, 1 for SMF1 and 2 for SMF2
	Type() uint16

	smfformat() // make the implementation exclusive to this package
}

// format is an implementation of Format
type format uint16

func (f format) Type() uint16 {
	return uint16(f)
}

func (f format) smfformat() {}

func (f format) String() string {
	switch f {
	case SMF0:
		return "SMF0 (singletrack)"
	case SMF1:
		return "SMF1 (multitrack)"
	case SMF2:
		return "SMF2 (sequential tracks)"
	}
	panic("unreachable")
}
*/
