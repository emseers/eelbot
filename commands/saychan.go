package commands

import (
	"strings"

	"github.com/emseers/eelbot"
)

// SayChanCommand returns an *eelbot.Command that replies with the given message on the given channel.
func SayChanCommand() *eelbot.Command {
	return &eelbot.Command{
		MinArgs: 2,
		MaxArgs: -1,
		Summary: "Says the given message in the given channel.",
		Usage: `/%[1]s CHAN MSG

Says the given message on the channel ID specified by CHAN.

Examples:
  /%[1]s 123456789 Hello world!
`,
		Eval: func(bot *eelbot.Bot, meta *eelbot.Meta, args []string) error {
			bot.DeleteMsg(meta.ChannelID, meta.MessageID)
			bot.SendMsg(args[0], strings.Join(args[1:], " "))
			return nil
		},
	}
}
