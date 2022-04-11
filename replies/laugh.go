package replies

import (
	"regexp"

	"github.com/emseers/eelbot"
)

var (
	laughExps = []*regexp.Regexp{
		regexp.MustCompile(`(?i)\bl+[o0]+l+\b`),
		regexp.MustCompile(`(?i)\bl+m+f*a+[o0]+\b`),
		regexp.MustCompile(`(?i)\br+[o0]+t*f+l+\b`),
	}
)

// LaughReply returns an *eelbot.Reply that has the given percent chance to trigger a reply on valid matches.
func LaughReply(percent int) *eelbot.Reply {
	return &eelbot.Reply{
		Eval: func(bot *eelbot.Bot, meta *eelbot.Meta, msg string) bool {
			if match(msg, laughExps...) && roll(percent) {
				bot.SendMsg(meta.ChannelID, "lol")
				return true
			}
			return false
		},
	}
}
