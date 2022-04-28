package replies

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/emseers/eelbot"
)

var (
	goodbyePrefixes = []string{
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

	goodbyeReplies = []string{
		"Night",
		"Good night",
		"G'night",
		"Have a nice night",
		"Have a good night",
		"Have a nice evening",
		"Have a good evening",
		"Have a nice day",
		"Have a good day",
		"Bye",
		"Bye now",
		"Goodbye",
		"Bye bye",
		"Bye bye 👋",
		"See ya",
		"👋",
		"God be with you",
		"auf Wiedersehen",
		"Aloha",
		"Adieu",
		"Cheerio",
		"Ciao",
		"Peace",
		"Peace out",
		"Farewell",
		"Later",
		"Laters",
		"See you later",
		"See you again later",
		"Catch you later",
		"Till next time",
		"Talk to you later",
		"Take care",
		"Take it easy",
		"Until next time",
		"Nice talking to you",
	}
)

func init() {
	replies["goodbye"] = goodbyeFromConfig
}

func goodbyeFromConfig(_ map[string]any, percent int, minDelay, maxDelay time.Duration) (*eelbot.Reply, error) {
	return GoodbyeReply(percent, minDelay, maxDelay), nil
}

// GoodbyeReply returns an *eelbot.Reply that has the given percent chance to trigger a reply on valid matches.
func GoodbyeReply(percent int, minDelay, maxDelay time.Duration) *eelbot.Reply {
	return &eelbot.Reply{
		Eval: func(s eelbot.Session, m *discordgo.MessageCreate) bool {
			if hasPrefix(m.Content, goodbyePrefixes...) && roll(percent) {
				asyncReply(s, m.ChannelID, randElem(goodbyeReplies), minDelay, maxDelay)
				return true
			}
			return false
		},
	}
}
