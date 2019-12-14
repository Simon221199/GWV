package main

import (
	"strconv"
)

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func isNumber(str string) bool {

	if _, err := strconv.Atoi(str); err == nil {
		return true
	}

	return false
}
