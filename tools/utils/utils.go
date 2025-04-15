// Package utils provides utility functions for the AWX package.
package utils

import "strconv"

// AtoiDefault takes a string and a defaultValue. If the string cannot be converted, defaultValue is returned.
func AtoiDefault(s string, defaultValue *int) *int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return defaultValue
	}
	return &n
}

// ItoaDefault takes an int pointer and a defaultValue. If the int pointer is nil, defaultValue is returned.
func ItoaDefault(i *int, defaultValue string) string {
	if i == nil {
		return defaultValue
	}
	return strconv.Itoa(*i)
}
