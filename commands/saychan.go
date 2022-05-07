package commands

import (
	"database/sql"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/emseers/eelbot"
)

func init() {
	commands["saychan"] = saychanFromConfig
}

func saychanFromConfig(map[string]any, *sql.DB, time.Duration) (*eelbot.Command, error) {
	return SayChanCommand(), nil
}

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
		Eval: func(s eelbot.Session, m *discordgo.MessageCreate, args []string) error {
			s.ChannelMessageDelete(m.ChannelID, m.ID)

			// Don't use args and rather use the raw input directly.
			s.ChannelMessageSend(args[0], strings.SplitN(m.Content, " ", 3)[2])
			return nil
		},
	}
}
