package msg

import (
	"math/rand"
	"strings"
)

// IsMorningGreet checks if a string begins with a possible morning greeting.
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

// IsGoodbyeGreet checks if a string begins with a possible goodbye greeting.
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

// MorningGreet returns a morning greeting.
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

// GoodbyeGreet returns a goodbye message.
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

// YellResponse returns a yell response message.
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
