package msg

import (
	"math/rand"
	"regexp"
)

func flip() string {
	if rand.Int63()&(1<<62) == 0 {
		return "Heads"
	}

	return "Tails"
}

func toValidASCII(in string) (out string) {
	rexp := regexp.MustCompile("[[:^ascii:]]")
	out = rexp.ReplaceAllLiteralString(in, "")
	return
}

func toAlphabetOnly(in string) (out string) {
	rexp := regexp.MustCompile("[^a-zA-Z]+")
	out = rexp.ReplaceAllLiteralString(in, "")
	return
}
