package commands

import (
	"database/sql"

	"github.com/bwmarrin/discordgo"
	"github.com/emseers/eelbot"
	"gopkg.in/ini.v1"
)

func init() {
	commands["channel"] = channelFromConfig
}

func channelFromConfig(*ini.Section, *sql.DB) (*eelbot.Command, error) {
	return ChannelCommand(), nil
}

// ChannelCommand returns an *eelbot.Command that replies with the channel ID of the command message.
func ChannelCommand() *eelbot.Command {
	return &eelbot.Command{
		MinArgs: 0,
		MaxArgs: 0,
		Summary: "Posts the current channel ID.",
		Eval: func(s *discordgo.Session, m *discordgo.MessageCreate, _ []string) error {
			s.ChannelMessageSend(m.ChannelID, m.ChannelID)
			return nil
		},
	}
}
