package msg

import (
	"math/rand"
	"regexp"
	"strings"
)

func isHelloGreet(msg string) bool {
	greetings := []string{
		"mornin",
		"g'mornin",
		"gmornin",
		"gud mornin",
		"good mornin",
		"hello",
	}

	for _, greeting := range greetings {
		if strings.HasPrefix(msg, greeting) {
			return true
		}
	}

	return false
}

func isGoodbyeGreet(msg string) bool {
	greetings := []string{
		"gnight",
		"g'night",
		"gnite",
		"g'nite",
		"good night",
		"goodnight",
		"good nite",
		"goodnite",
		"gud night",
		"gudnight",
		"gud nite",
		"gudnite",
		"bye",
		"good bye",
		"goodbye",
		"gbye",
		"g'bye",
	}

	for _, greeting := range greetings {
		if strings.HasPrefix(msg, greeting) {
			return true
		}
	}

	return false
}

func isLaugh(msg string) bool {
	rexp := regexp.MustCompile("l[o0]+l")
	if rexp.Match([]byte(msg)) && rand.Intn(100) < 17 {
		return true
	}

	return false
}

func helloGreet() string {
	greetings := []string{
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
		"Hello there",
		"How's it going?",
		"How're you doing?",
		"Jesus Christ be praised",
	}

	return greetings[rand.Intn(len(greetings))]
}

func goodbyeGreet() string {
	greetings := []string{
		"Night",
		"Good night",
		"Nite",
		"Good nite",
		"G'night",
		"Have a good night",
		"Have a good evening",
		"Bye",
		"Bye now",
		"Goodbye",
		"Bye bye",
		"Bye bye 👋",
		"God be with you",
		"auf Wiedersehen",
		"Aloha",
		"Adieu",
		"Cheerio",
		"Ciao",
		"Farewell",
		"See you later",
		"See you again later",
		"Catch you later",
		"Till next time",
		"Take care",
	}

	return greetings[rand.Intn(len(greetings))]
}

func yellResponse() string {
	responses := []string{
		"Stop yelling",
		"Stop squeaking",
		"Stop squealing",
		"Stop howling",
		"Not so loud",
		"Saying things louder doesn't make you right",
		"Why so serious?",
		"Calm down",
		"Simmer down",
		"Settle down",
		"Chillax",
		"Chill, man",
		"Chill out, dude",
		"Yur alright, boahhh",
	}

	return responses[rand.Intn(len(responses))]
}
