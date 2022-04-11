package commands

import (
	"strings"

	"github.com/emseers/eelbot"
)

// SayCommand returns an *eelbot.Command that replies with the given message.
func SayCommand() *eelbot.Command {
	return &eelbot.Command{
		MinArgs: 1,
		MaxArgs: -1,
		Summary: "Says the given message.",
		Usage: `/%[1]s MSG

Says the given message on the current channel.

Examples:
  /%[1]s Hello world!
`,
		Eval: func(bot *eelbot.Bot, meta *eelbot.Meta, args []string) error {
			bot.DeleteMsg(meta.ChannelID, meta.MessageID)
			bot.SendMsg(meta.ChannelID, strings.Join(args, " "))
			return nil
		},
	}
}
