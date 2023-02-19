package eelbot

import (
	"io"

	"github.com/bwmarrin/discordgo"
)

// Session is a discord session.
type Session interface {
	AddHandler(handler any) func()

	Channel(channelID string, options ...discordgo.RequestOption) (*discordgo.Channel, error)
	ChannelMessage(channelID, messageID string, options ...discordgo.RequestOption) (*discordgo.Message, error)
	ChannelMessages(channelID string, limit int, beforeID, afterID, aroundID string,
		options ...discordgo.RequestOption) ([]*discordgo.Message, error)
	ChannelMessagesPinned(channelID string, options ...discordgo.RequestOption) ([]*discordgo.Message, error)

	ChannelTyping(channelID string, options ...discordgo.RequestOption) (err error)
	ChannelMessageSend(channelID, content string, options ...discordgo.RequestOption) (*discordgo.Message, error)
	ChannelMessageSendEmbeds(channelID string, embeds []*discordgo.MessageEmbed,
		options ...discordgo.RequestOption) (*discordgo.Message, error)
	ChannelMessageSendTTS(channelID, content string, options ...discordgo.RequestOption) (*discordgo.Message, error)
	ChannelFileSend(channelID, name string, r io.Reader, options ...discordgo.RequestOption) (*discordgo.Message, error)
	ChannelFileSendWithMessage(channelID, content, name string, r io.Reader,
		options ...discordgo.RequestOption) (*discordgo.Message, error)

	ChannelMessageEdit(channelID, messageID, content string,
		options ...discordgo.RequestOption) (*discordgo.Message, error)
	ChannelMessageEditEmbeds(channelID, messageID string, embeds []*discordgo.MessageEmbed,
		options ...discordgo.RequestOption) (*discordgo.Message, error)

	ChannelMessageDelete(channelID, messageID string, options ...discordgo.RequestOption) error
	ChannelMessagesBulkDelete(channelID string, messages []string, options ...discordgo.RequestOption) error

	ChannelMessagePin(channelID, messageID string, options ...discordgo.RequestOption) error
	ChannelMessageUnpin(channelID, messageID string, options ...discordgo.RequestOption) error

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
