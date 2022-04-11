package replies

import (
	"regexp"

	"github.com/emseers/eelbot"
)

var (
	questionExp = regexp.MustCompile(`^(\?\s*)+$`)
)

// QuestionReply returns an *eelbot.Reply that has the given percent chance to trigger a reply on valid matches.
func QuestionReply(percent int) *eelbot.Reply {
	return &eelbot.Reply{
		Eval: func(bot *eelbot.Bot, meta *eelbot.Meta, msg string) bool {
			if questionExp.MatchString(msg) && roll(percent) {
				bot.SendMsg(meta.ChannelID, "Don't questionmark me")
				return true
			}
			return false
		},
	}
}
