package vlq

// stolen from http://stackoverflow.com/questions/19239449/how-do-i-reverse-an-array-in-go
func reverse(b []byte) {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
}

const (
	vlqContinue = 128
	vlqMask     = 127
)

// limit the largest possible value to int32
/*
The largest number which is allowed is 0FFFFFFF so that the variable-length representations must fit in 32
bits in a routine to write variable-length numbers. Theoretically, larger numbers are possible, but 2 x 10 8
96ths of a beat at a fast tempo of 500 beats per minute is four days, long enough for any delta-time!
*/

// Variable-Length Quantity (VLQ) is an way of representing arbitrary
// see https://blogs.infosupport.com/a-primer-on-vlq/
// we use the variant of the midi-spec
// stolen and converted to go from https://github.com/dvberkel/VLQKata/blob/master/src/main/java/nl/dvberkel/kata/Kata.java#L12

// Encode encodes the given value as variable length quantity
func Encode(n uint32) (out []byte) {
	var quo, rem uint32
	quo = n / vlqContinue
	rem = n % vlqContinue

	out = append(out, byte(rem))

	for quo > 0 {
		out = append(out, byte(quo)|vlqContinue)
		quo = quo / vlqContinue
		// rem = quo % vlqContinue
	}

	reverse(out)
	return
}

// Decode decodes a variable length quantity
func Decode(source []byte) (num uint32) {

	for i := 0; i < len(source); i++ {
		var n = uint32(source[i] & vlqMask)
		for (source[i] & vlqContinue) != 0 {
			i++
			n *= 128
			n += uint32(source[i] & vlqMask)
		}
		num += n
	}

	return
}
