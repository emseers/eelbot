package replies

import "github.com/emseers/eelbot"

var (
	helloPrefixes = []string{
		"mornin",
		"g'mornin",
		"gmornin",
		"gud mornin",
		"good mornin",
		"hello",
		"hiya",
	}

	helloReplies = []string{
		"Morning",
		"Good morning",
		"Mornin'",
		"Good mornin'",
		"G'morning",
		"Top of the morning to you",
		"Morning to you too",
		"G'day",
		"G'day mate",
		"Hi",
		"Hello",
		"Hiya",
		"Howdy",
		"Hello there",
		"How's it going?",
		"How's it hanging?",
		"How're you doing?",
		"How have you been?",
		"How do you do?",
		"How are things?",
		"How's life?",
		"How's your day?",
		"How's your day going?",
		"👋",
		"What's up?",
		"What's up",
		"Wassup",
		"What's going on?",
		"Good to see you",
		"Nice to see you",
		"Long time no see",
		"Yo",
		"Yo!",
		"Aloha",
		"Jesus Christ be praised",
	}
)

// HelloReply returns an *eelbot.Reply that has the given percent chance to trigger a reply on valid matches.
func HelloReply(percent int) *eelbot.Reply {
	return &eelbot.Reply{
		Eval: func(bot *eelbot.Bot, meta *eelbot.Meta, msg string) bool {
			if hasPrefix(msg, helloPrefixes...) && roll(percent) {
				bot.SendMsg(meta.ChannelID, randElem(helloReplies))
				return true
			}
			return false
		},
	}
}
