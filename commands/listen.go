package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/emseers/eelbot"
)

// ListenCommand returns an *eelbot.Command that sets the bot's status to listen to the given song.
func ListenCommand() *eelbot.Command {
	return &eelbot.Command{
		MinArgs: 0,
		MaxArgs: 1,
		Summary: "Listens to a song.",
		Eval: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
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
