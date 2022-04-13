package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/emseers/eelbot"
)

// PingCommand returns an *eelbot.Command that replies with "Pong".
func PingCommand() *eelbot.Command {
	return &eelbot.Command{
		MinArgs: 0,
		MaxArgs: 0,
		Summary: "Replies with \"Pong\".",
		Eval: func(s *discordgo.Session, m *discordgo.MessageCreate, _ []string) error {
			s.ChannelMessageSend(m.ChannelID, "Pong")
			return nil
		},
	}
}
