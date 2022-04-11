package commands

import (
	"github.com/emseers/eelbot"
)

// ListenCommand returns an *eelbot.Command that sets the bot's status to listen to the given song.
func ListenCommand() *eelbot.Command {
	return &eelbot.Command{
		MinArgs: 0,
		MaxArgs: 1,
		Summary: "Listens to a song.",
		Eval: func(bot *eelbot.Bot, meta *eelbot.Meta, args []string) error {
			bot.DeleteMsg(meta.ChannelID, meta.MessageID)
			var song string
			if len(args) == 1 {
				song = args[0]
			}
			bot.UpdateListeningStatus(song)
			return nil
		},
	}
}
