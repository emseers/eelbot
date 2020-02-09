package msg

import (
	"math/rand"
	"regexp"
)

// Flip returns "Heads" or "Tails".
func Flip() string {
	if rand.Int63()&(1<<62) == 0 {
		return "Heads"
	}

	return "Tails"
}

// ToValidAscii converts a UTF-8 string to ASCII.
func ToValidAscii(s string) string {
	re := regexp.MustCompile("[[:^ascii:]]")
	t := re.ReplaceAllLiteralString(s, "")
	return t
}

// ToAlphabetOnly removes any non alphanumeric characters from the given string.
func ToAlphabetOnly(s string) string {
	re := regexp.MustCompile("[^a-zA-Z]+")
	t := re.ReplaceAllLiteralString(s, "")
	return t
}
