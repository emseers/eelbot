package replies

import (
	"regexp"

	"github.com/bwmarrin/discordgo"
	"github.com/emseers/eelbot"
)

var (
	questionExp = regexp.MustCompile(`^(\?\s*)+$`)
)

func init() {
	replies["question"] = questionFromConfig
}

func questionFromConfig(_ map[string]any, percent int) (*eelbot.Reply, error) {
	return QuestionReply(percent), nil
}

// QuestionReply returns an *eelbot.Reply that has the given percent chance to trigger a reply on valid matches.
func QuestionReply(percent int) *eelbot.Reply {
	return &eelbot.Reply{
		Eval: func(s eelbot.Session, m *discordgo.MessageCreate) bool {
			if questionExp.MatchString(m.Content) && roll(percent) {
				s.ChannelMessageSend(m.ChannelID, "Don't questionmark me")
				return true
			}
			return false
		},
	}
}
