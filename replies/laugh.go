package replies

import (
	"regexp"

	"github.com/bwmarrin/discordgo"
	"github.com/emseers/eelbot"
)

var (
	laughExps = []*regexp.Regexp{
		regexp.MustCompile(`(?i)\bl+[o0]+l+\b`),
		regexp.MustCompile(`(?i)\bl+m+f*a+[o0]+\b`),
		regexp.MustCompile(`(?i)\br+[o0]+t*f+l+\b`),
	}
)

func init() {
	replies["laugh"] = laughFromConfig
}

func laughFromConfig(_ map[string]any, percent int) (*eelbot.Reply, error) {
	return LaughReply(percent), nil
}

// LaughReply returns an *eelbot.Reply that has the given percent chance to trigger a reply on valid matches.
func LaughReply(percent int) *eelbot.Reply {
	return &eelbot.Reply{
		Eval: func(s eelbot.Session, m *discordgo.MessageCreate) bool {
			if match(m.Content, laughExps...) && roll(percent) {
				s.ChannelMessageSend(m.ChannelID, "lol")
				return true
			}
			return false
		},
	}
}
