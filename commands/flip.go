package commands

import (
	"database/sql"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/emseers/eelbot"
)

func init() {
	commands["flip"] = flipFromConfig
}

func flipFromConfig(map[string]any, *sql.DB, time.Duration) (*eelbot.Command, error) {
	return FlipCommand(), nil
}

// FlipCommand returns an *eelbot.Command that replies with the result of a coin toss.
func FlipCommand() *eelbot.Command {
	return &eelbot.Command{
		MinArgs: 0,
		MaxArgs: 0,
		Summary: "Flips a coin.",
		Eval: func(s eelbot.Session, m *discordgo.MessageCreate, _ []string) error {
			var result string
			switch {
			case Rand.Intn(6000) == 0:
				result = "Landed on edge"
			case Rand.Intn(2) == 0:
				result = "Heads"
			default:
				result = "Tails"
			}
			s.ChannelMessageSend(m.ChannelID, result)
			return nil
		},
	}
}
