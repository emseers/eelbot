package commands

import (
	"database/sql"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/emseers/eelbot"
)

func init() {
	commands["play"] = playFromConfig
}

func playFromConfig(map[string]any, *sql.DB, time.Duration) (*eelbot.Command, error) {
	return PlayCommand(), nil
}

// PlayCommand returns an *eelbot.Command that sets the bot's status to play the given game.
func PlayCommand() *eelbot.Command {
	return &eelbot.Command{
		MinArgs: 0,
		MaxArgs: -1,
		Summary: "Plays a game.",
		Eval: func(s eelbot.Session, m *discordgo.MessageCreate, args []string) error {
			s.ChannelMessageDelete(m.ChannelID, m.ID)
			var game string
			if len(args) > 0 {
				// Don't use args and rather use the raw input directly.
				game = strings.SplitN(m.Content, " ", 2)[1]
			}
			s.UpdateGameStatus(0, game)
			return nil
		},
	}
}
