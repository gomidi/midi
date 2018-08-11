package key

import (
	"github.com/gomidi/midi/midimessage/meta"
)

// ok
func CMaj() meta.KeySignature {
	return key(0, 0, true, false)
}

// ok
func DMaj() meta.KeySignature {
	return key(2, 2, true, false)
}

// ok
func EMaj() meta.KeySignature {
	return key(4, 4, true, false)
}

// ok
func FSharpMaj() meta.KeySignature {
	return key(6, 6, true, false)
}

// ok
func GMaj() meta.KeySignature {
	return key(7, 1, true, false)
}

// ok
func AMaj() meta.KeySignature {
	return key(9, 3, true, false)
}

// ok
func BMaj() meta.KeySignature {
	return key(11, 5, true, false)
}

// ok
func FMaj() meta.KeySignature {
	return key(5, 1, true, true)
}

// ok
func BFlatMaj() meta.KeySignature {
	return key(10, 2, true, true)
}

// ok
func EFlatMaj() meta.KeySignature {
	return key(3, 3, true, true)
}

// ok
func AFlatMaj() meta.KeySignature {
	return key(8, 4, true, true)
}

// ok
func DFlatMaj() meta.KeySignature {
	return key(1, 5, true, true)
}

// ok
func GFlatMaj() meta.KeySignature {
	return key(6, 6, true, true)
}

/*
func CFlatMaj() meta.KeySignature {
	return key(11, 5, true, false)
}

func CSharpMaj() meta.KeySignature {
	return key(1, 5, true, true)
}

func DSharpMaj() meta.KeySignature {
	return key(3, 3, true, true)
}

func ESharpMaj() meta.KeySignature {
	return key(5, 1, true, true)
}

func FFlatMaj() meta.KeySignature {
	return key(4, 4, true, false)
}

func GSharpMaj() meta.KeySignature {
	return key(8, 4, true, true)
}

func ASharpMaj() meta.KeySignature {
	return key(10, 2, true, true)
}

func BSharpMaj() meta.KeySignature {
	return key(0, 0, true, false)
}
*/
