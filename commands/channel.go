package commands

import (
	"github.com/emseers/eelbot"
)

// ChannelCommand returns an *eelbot.Command that replies with the channel ID of the command message.
func ChannelCommand() *eelbot.Command {
	return &eelbot.Command{
		MinArgs: 0,
		MaxArgs: 0,
		Summary: "Posts the current channel ID.",
		Eval: func(bot *eelbot.Bot, meta *eelbot.Meta, _ []string) error {
			bot.SendMsg(meta.ChannelID, meta.ChannelID)
			return nil
		},
	}
}
