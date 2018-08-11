package key

import (
	"github.com/gomidi/midi/midimessage/meta"
)

// ok
func AMin() meta.KeySignature {
	return key(9, 0, false, false)
}

// ok
func BMin() meta.KeySignature {
	return key(11, 2, false, false)
}

// ok
func CSharpMin() meta.KeySignature {
	return key(1, 4, false, false)
}

// ok
func DSharpMin() meta.KeySignature {
	return key(3, 6, false, false)
}

// ok
func EMin() meta.KeySignature {
	return key(4, 1, false, false)
}

// ok
func FSharpMin() meta.KeySignature {
	return key(6, 3, false, false)
}

// ok
func GSharpMin() meta.KeySignature {
	return key(8, 5, false, false)
}

// ok
func DMin() meta.KeySignature {
	return key(2, 1, false, true)
}

// ok
func GMin() meta.KeySignature {
	return key(7, 2, false, true)
}

// ok
func CMin() meta.KeySignature {
	return key(0, 3, false, true)
}

// ok
func FMin() meta.KeySignature {
	return key(5, 4, false, true)
}

// ok
func BFlatMin() meta.KeySignature {
	return key(10, 5, false, true)
}

// ok
func EFlatMin() meta.KeySignature {
	return key(3, 6, false, true)
}

/*
func CFlatMin() meta.KeySignature {
	return key(11, 2, false, false)
}


func ESharpMin() meta.KeySignature {
	return key(5, 4, false, true)
}

func FFlatMin() meta.KeySignature {
	return key(4, 1, false, false)
}

func GFlatMin() meta.KeySignature {
	return key(5, 3, false, false)
}

func AFlatMin() meta.KeySignature {
	return key(8, 5, false, false)
}

func ASharpMin() meta.KeySignature {
	return key(10, 5, false, true)
}

func BSharpMin() meta.KeySignature {
	return key(0, 3, false, true)
}

func DFlatMin() meta.KeySignature {
	return key(1, 4, false, false)
}



*/
