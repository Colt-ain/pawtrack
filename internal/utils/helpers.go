package utils

import "strconv"

// Atoi converts a string to int, returns 0 on error
func Atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

// Max returns the maximum of two numbers
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// ClampAtoi converts a string to int and clamps the value
func ClampAtoi(s string, lo, hi int) int {
	i := Atoi(s)
	if i < lo {
		return lo
	}
	if i > hi {
		return hi
	}
	return i
}
