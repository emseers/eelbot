package commands

import (
	"database/sql"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/emseers/eelbot"
	"gopkg.in/ini.v1"
)

func init() {
	commands["say"] = sayFromConfig
}

func sayFromConfig(*ini.Section, *sql.DB) (*eelbot.Command, error) {
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
		Eval: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
			s.ChannelMessageDelete(m.ChannelID, m.ID)
			s.ChannelMessageSend(m.ChannelID, strings.Join(args, " "))
			return nil
		},
	}
}
