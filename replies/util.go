package replies

import (
	"regexp"
	"strings"
)

var (
	nonAlphaExp = regexp.MustCompile("[^a-zA-Z]+")
)

func roll(percent int) bool {
	return Rand.Intn(100) < percent
}

func randElem(s []string) string {
	return s[Rand.Intn(len(s))]
}

func toAlphabetsOnly(s string) string {
	return nonAlphaExp.ReplaceAllLiteralString(s, "")
}

func hasPrefix(s string, prefixes ...string) bool {
	s = strings.ToLower(s)
	for _, prefix := range prefixes {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}

func match(s string, rexps ...*regexp.Regexp) bool {
	for _, rexp := range rexps {
		if rexp.MatchString(s) {
			return true
		}
	}
	return false
}
