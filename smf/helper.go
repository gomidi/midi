package smf

import (
	"strconv"
)

// roundFloat rounds the given float by the given decimals after the dot
func roundFloat(x float64, decimals int) float64 {
	// return roundFloat(x, numDig(x)+decimals)
	frep := strconv.FormatFloat(x, 'f', decimals, 64)
	f, _ := strconv.ParseFloat(frep, 64)
	return f
}
