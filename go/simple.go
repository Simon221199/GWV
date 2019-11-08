package main

import (
	"crypto/rand"
	"fmt"
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

func Uuid() string {

	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%X-%X-%X-%X-%X", bytes[0:4], bytes[4:6], bytes[6:8], bytes[8:10], bytes[10:])
}