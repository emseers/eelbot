package msg

import (
	"math/rand"
	"regexp"
)

func Flip() string {
	if rand.Int63()&(1<<62) == 0 {
		return "Heads"
	}

	return "Tails"
}

func ToValidAscii(s string) string {
	re := regexp.MustCompile("[[:^ascii:]]")
	t := re.ReplaceAllLiteralString(s, "")
	return t
}

func ToAlphabetOnly(s string) string {
	re := regexp.MustCompile("[^a-zA-Z]+")
	t := re.ReplaceAllLiteralString(s, "")
	return t
}

func YellResponse() string {
	responses := [5]string{
		"Stop yelling",
		"Saying things louder doesn't make you right",
		"Chillax",
		"Why so serious?",
		"Calm down",
	}

	return responses[rand.Intn(len(responses))]
}
