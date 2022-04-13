package commands

import (
	"math/rand"

	"github.com/bwmarrin/discordgo"
	"github.com/emseers/eelbot"
)

// FlipCommand returns an *eelbot.Command that replies with the result of a coin toss.
func FlipCommand() *eelbot.Command {
	return &eelbot.Command{
		MinArgs: 0,
		MaxArgs: 0,
		Summary: "Flips a coin.",
		Eval: func(s *discordgo.Session, m *discordgo.MessageCreate, _ []string) error {
			var result string
			switch {
			case rand.Intn(6000) == 0:
				result = "Landed on edge"
			case rand.Intn(2) == 0:
				result = "Heads"
			default:
				result = "Tails"
			}
			s.ChannelMessageSend(m.ChannelID, result)
			return nil
		},
	}
}
