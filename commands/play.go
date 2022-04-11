package commands

import (
	"github.com/emseers/eelbot"
)

// PlayCommand returns an *eelbot.Command that sets the bot's status to play the given game.
func PlayCommand() *eelbot.Command {
	return &eelbot.Command{
		MinArgs: 0,
		MaxArgs: 1,
		Summary: "Plays a game.",
		Eval: func(bot *eelbot.Bot, meta *eelbot.Meta, args []string) error {
			bot.DeleteMsg(meta.ChannelID, meta.MessageID)
			var game string
			if len(args) == 1 {
				game = args[0]
			}
			bot.UpdateGameStatus(game)
			return nil
		},
	}
}
