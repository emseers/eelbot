package commands

import (
	"database/sql"

	"github.com/bwmarrin/discordgo"
	"github.com/emseers/eelbot"
)

func init() {
	commands["play"] = playFromConfig
}

func playFromConfig(map[string]any, *sql.DB) (*eelbot.Command, error) {
	return PlayCommand(), nil
}

// PlayCommand returns an *eelbot.Command that sets the bot's status to play the given game.
func PlayCommand() *eelbot.Command {
	return &eelbot.Command{
		MinArgs: 0,
		MaxArgs: 1,
		Summary: "Plays a game.",
		Eval: func(s eelbot.Session, m *discordgo.MessageCreate, args []string) error {
			s.ChannelMessageDelete(m.ChannelID, m.ID)
			var game string
			if len(args) == 1 {
				game = args[0]
			}
			s.UpdateGameStatus(0, game)
			return nil
		},
	}
}
