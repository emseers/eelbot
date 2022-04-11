package commands

import (
	"github.com/emseers/eelbot"
)

// PingCommand returns an *eelbot.Command that replies with "Pong".
func PingCommand() *eelbot.Command {
	return &eelbot.Command{
		MinArgs: 0,
		MaxArgs: 0,
		Summary: "Replies with \"Pong\".",
		Eval: func(bot *eelbot.Bot, meta *eelbot.Meta, _ []string) error {
			bot.SendMsg(meta.ChannelID, "Pong")
			return nil
		},
	}
}
