package commands

import (
	"database/sql"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/emseers/eelbot"
)

func init() {
	commands["say"] = sayFromConfig
}

func sayFromConfig(map[string]any, *sql.DB) (*eelbot.Command, error) {
	return SayCommand(), nil
}

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
		Eval: func(s eelbot.Session, m *discordgo.MessageCreate, args []string) error {
			s.ChannelMessageDelete(m.ChannelID, m.ID)

			// Don't use args and rather use the raw input directly.
			s.ChannelMessageSend(m.ChannelID, strings.SplitN(m.Content, " ", 2)[1])
			return nil
		},
	}
}
