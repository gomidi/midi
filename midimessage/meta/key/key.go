package key

import (
	"github.com/gomidi/midi/midimessage/meta"
)

func key(key, num uint8, isMajor, isFlat bool) meta.KeySignature {
	return meta.KeySignature{Key: key, IsMajor: isMajor, Num: num, IsFlat: isFlat}
}
