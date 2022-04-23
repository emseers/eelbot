package commands

import (
	"database/sql"

	"github.com/bwmarrin/discordgo"
	"github.com/emseers/eelbot"
	"gopkg.in/ini.v1"
)

func init() {
	commands["ping"] = pingFromConfig
}

func pingFromConfig(*ini.Section, *sql.DB) (*eelbot.Command, error) {
	return PingCommand(), nil
}

// PingCommand returns an *eelbot.Command that replies with "Pong".
func PingCommand() *eelbot.Command {
	return &eelbot.Command{
		MinArgs: 0,
		MaxArgs: 0,
		Summary: "Replies with \"Pong\".",
		Eval: func(s eelbot.Session, m *discordgo.MessageCreate, _ []string) error {
			s.ChannelMessageSend(m.ChannelID, "Pong")
			return nil
		},
	}
}
