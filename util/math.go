package util

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
