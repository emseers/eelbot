package replies

import (
	"strings"

	"github.com/emseers/eelbot"
)

var (
	capsReplies = []string{
		"Stop yelling",
		"Stop squeaking",
		"Stop squealing",
		"Stop squirming",
		"Stop howling",
		"Not so loud",
		"Saying things louder doesn't make you right",
		"Why so serious?",
		"Calm down",
		"Simmer down",
		"Settle down",
		"Chillax",
		"Chill, dude",
		"Chill out, dude",
	}
)

// CapsReply returns an *eelbot.Reply that has the given percent chance to trigger a reply on an all caps message that
// is at least minLen characters long.
func CapsReply(minLen, percent int) *eelbot.Reply {
	return &eelbot.Reply{
		Eval: func(bot *eelbot.Bot, meta *eelbot.Meta, msg string) bool {
			msg = toAlphabetsOnly(msg)
			if len(msg) >= minLen && msg == strings.ToUpper(msg) && roll(percent) {
				bot.SendMsg(meta.ChannelID, randElem(capsReplies))
			}
			return false
		},
	}
}
