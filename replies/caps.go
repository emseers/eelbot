package replies

import (
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
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

func init() {
	replies["caps"] = capsFromConfig
}

func capsFromConfig(opts map[string]any, percent int, minDelay, maxDelay time.Duration) (*eelbot.Reply, error) {
	minLen, ok := opts["min_len"].(float64)
	if !ok {
		minLen = 5
	}
	return CapsReply(percent, minDelay, maxDelay, int(minLen)), nil
}

// CapsReply returns an *eelbot.Reply that has the given percent chance to trigger a reply on an all caps message that
// is at least minLen characters long.
func CapsReply(percent int, minDelay, maxDelay time.Duration, minLen int) *eelbot.Reply {
	return &eelbot.Reply{
		Eval: func(s eelbot.Session, m *discordgo.MessageCreate) bool {
			msg := toAlphabetsOnly(m.Content)
			if len(msg) >= minLen && msg == strings.ToUpper(msg) && roll(percent) {
				asyncReply(s, m.ChannelID, randElem(capsReplies), minDelay, maxDelay)
				return true
			}
			return false
		},
	}
}
