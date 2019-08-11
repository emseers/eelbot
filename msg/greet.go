package msg

import (
	"math/rand"
	"strings"
)

func IsMorningGreet(msg string) bool {
	greetings := [5]string{
		"mornin",
		"g'mornin",
		"gmornin",
		"gud mornin",
		"hello",
	}

	for _, greeting := range greetings {
		if strings.HasPrefix(msg, greeting) {
			return true
		}
	}

	return false
}

func IsGoodbyeGreet(msg string) bool {
	greetings := [12]string{
		"gnight",
		"g'night",
		"gnite",
		"g'nite",
		"good night",
		"goodnight",
		"good nite",
		"goodnite",
		"bye",
		"good bye",
		"gbye",
		"goodbye",
	}

	for _, greeting := range greetings {
		if strings.HasPrefix(msg, greeting) {
			return true
		}
	}

	return false
}

func MorningGreet() string {
	greetings := [10]string{
		"Morning",
		"Good morning",
		"Mornin'",
		"Good mornin'",
		"G'morning",
		"Top of the morning to you",
		"Morning to you too",
		"Hi",
		"Hello",
		"Hiya",
	}

	return greetings[rand.Intn(len(greetings))]
}

func GoodbyeGreet() string {
	greetings := [10]string{
		"Night",
		"Good night",
		"Nite",
		"Good nite",
		"G'night",
		"Have a good night",
		"Have a good evening",
		"Bye",
		"Goodbye",
		"Bye bye :wave:",
	}

	return greetings[rand.Intn(len(greetings))]
}
