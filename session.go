package eelbot

import (
	"io"

	"github.com/bwmarrin/discordgo"
)

// Session is a discord session.
type Session interface {
	AddHandler(handler interface{}) func()

	Channel(channelID string) (*discordgo.Channel, error)
	ChannelMessage(channelID, messageID string) (*discordgo.Message, error)
	ChannelMessages(channelID string, limit int, beforeID, afterID, aroundID string) ([]*discordgo.Message, error)
	ChannelMessagesPinned(channelID string) ([]*discordgo.Message, error)

	ChannelTyping(channelID string) (err error)
	ChannelMessageSend(channelID, content string) (*discordgo.Message, error)
	ChannelMessageSendEmbeds(channelID string, embeds []*discordgo.MessageEmbed) (*discordgo.Message, error)
	ChannelMessageSendTTS(channelID, content string) (*discordgo.Message, error)
	ChannelFileSend(channelID, name string, r io.Reader) (*discordgo.Message, error)
	ChannelFileSendWithMessage(channelID, content, name string, r io.Reader) (*discordgo.Message, error)

	ChannelMessageEdit(channelID, messageID, content string) (*discordgo.Message, error)
	ChannelMessageEditEmbeds(channelID, messageID string, embeds []*discordgo.MessageEmbed) (*discordgo.Message, error)

	ChannelMessageDelete(channelID, messageID string) error
	ChannelMessagesBulkDelete(channelID string, messages []string) error

	ChannelMessagePin(channelID, messageID string) error
	ChannelMessageUnpin(channelID, messageID string) error

	UpdateGameStatus(idle int, name string) error
	UpdateStreamingStatus(idle int, name, url string) error
	UpdateListeningStatus(name string) error

	Open() error
	Close() error
}

// NewSession creates a new session from a bot token.
func NewSession(token string) (Session, error) {
	return discordgo.New("Bot " + token)
}
