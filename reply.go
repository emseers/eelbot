package eelbot

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

// Reply defines a reply for a non-command message.
type Reply struct {
	// Timeout specifies how much time needs to pass before this reply can be triggered again on the same channel.
	Timeout time.Duration

	// Eval should return whether the reply is applicable to the given message. If false is returned, it is expected
	// that no replies were sent.
	Eval func(s Session, m *discordgo.MessageCreate) bool

	m *chanMap
}

// RegisterReply registers a new reply.
func (bot *Bot) RegisterReply(r Reply) {
	r.m = newChanMap()
	bot.replies = append(bot.replies, &r)
}
