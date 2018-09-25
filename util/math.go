package util

import "math"

func Min2(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Min3(a, b, c int) int {
	mi := a
	if b < mi {
		mi = b
	}
	if c < mi {
		mi = c
	}
	return mi
}

// Round rounds a given number to the nearest hundred (usually for prices)
func Round(x float64) float64 {
	return math.Round(x*100) / 100
}
