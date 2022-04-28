package commands

import (
	"database/sql"

	"github.com/bwmarrin/discordgo"
	"github.com/emseers/eelbot"
)

func init() {
	commands["listen"] = listenFromConfig
}

func listenFromConfig(map[string]any, *sql.DB) (*eelbot.Command, error) {
	return ListenCommand(), nil
}

// ListenCommand returns an *eelbot.Command that sets the bot's status to listen to the given song.
func ListenCommand() *eelbot.Command {
	return &eelbot.Command{
		MinArgs: 0,
		MaxArgs: 1,
		Summary: "Listens to a song.",
		Eval: func(s eelbot.Session, m *discordgo.MessageCreate, args []string) error {
			s.ChannelMessageDelete(m.ChannelID, m.ID)
			var song string
			if len(args) == 1 {
				song = args[0]
			}
			s.UpdateListeningStatus(song)
			return nil
		},
	}
}
