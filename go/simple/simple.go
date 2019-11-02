package simple

import "strconv"

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func IsNumber(str string) bool {

	if _, err := strconv.Atoi(str); err == nil {
		return true
	}

	return false
}