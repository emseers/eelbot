package replies

import (
	"regexp"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/emseers/eelbot"
)

var (
	questionExp = regexp.MustCompile(`^(\?\s*)+$`)
)

func init() {
	replies["question"] = questionFromConfig
}

func questionFromConfig(_ map[string]any, percent int, minDelay, maxDelay time.Duration) (*eelbot.Reply, error) {
	return QuestionReply(percent, minDelay, maxDelay), nil
}

// QuestionReply returns an *eelbot.Reply that has the given percent chance to trigger a reply on valid matches.
func QuestionReply(percent int, minDelay, maxDelay time.Duration) *eelbot.Reply {
	return &eelbot.Reply{
		Eval: func(s eelbot.Session, m *discordgo.MessageCreate) bool {
			if questionExp.MatchString(m.Content) && roll(percent) {
				asyncReply(s, m.ChannelID, "Don't questionmark me", minDelay, maxDelay)
				return true
			}
			return false
		},
	}
}
