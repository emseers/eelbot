package replies_test

import (
	"io"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
)

const (
	testChannelID = "123456789"
)

// Creates a Message with the provided arguments.
//  0: Content
//  1: ChannelID
//  2: MessageID
func newMsg(args ...string) *discordgo.Message {
	m := new(discordgo.Message)
	if len(args) > 0 {
		m.Content = args[0]
	}
	if len(args) > 1 {
		m.ChannelID = args[1]
	}
	if len(args) > 2 {
		m.ID = args[2]
	}
	return m
}

// Creates a MessageCreate with the provided arguments. Arguments are the same as newMsg.
func newMsgCreate(args ...string) *discordgo.MessageCreate {
	return &discordgo.MessageCreate{Message: newMsg(args...)}
}

func newTestSession() *testsession {
	return &testsession{
		messages: make(map[string]*strings.Builder),
	}
}

type testsession struct {
	counter  int
	messages map[string]*strings.Builder
}

func (s *testsession) AddHandler(handler interface{}) func() {
	return func() {}
}

func (s *testsession) Channel(channelID string) (*discordgo.Channel, error) {
	panic("Channel not implemented")
}

func (s *testsession) ChannelMessage(channelID, messageID string) (*discordgo.Message, error) {
	panic("ChannelMessage not implemented")
}

func (s *testsession) ChannelMessages(channelID string, limit int, beforeID, afterID, aroundID string) (
	[]*discordgo.Message, error) {
	panic("ChannelMessages not implemented")
}

func (s *testsession) ChannelMessagesPinned(channelID string) ([]*discordgo.Message, error) {
	panic("ChannelMessagesPinned not implemented")
}

func (s *testsession) ChannelTyping(channelID string) (err error) {
	panic("ChannelTyping not implemented")
}

func (s *testsession) ChannelMessageSend(channelID string, content string) (*discordgo.Message, error) {
	b, ok := s.messages[channelID]
	if !ok {
		b = new(strings.Builder)
		s.messages[channelID] = b
	}
	b.WriteString(content)
	b.WriteByte('\n')
	s.counter++
	return newMsg(content, channelID, strconv.Itoa(s.counter)), nil
}

func (s *testsession) ChannelMessageSendEmbeds(channelID string, embeds []*discordgo.MessageEmbed) (
	*discordgo.Message, error) {
	panic("ChannelMessageSendEmbeds not implemented")
}

func (s *testsession) ChannelMessageSendTTS(channelID string, content string) (*discordgo.Message, error) {
	panic("ChannelMessageSendTTS not implemented")
}

func (s *testsession) ChannelFileSend(channelID, name string, r io.Reader) (*discordgo.Message, error) {
	panic("ChannelFileSend not implemented")
}

func (s *testsession) ChannelFileSendWithMessage(channelID, content string, name string, r io.Reader) (
	*discordgo.Message, error) {
	panic("ChannelFileSendWithMessage not implemented")
}

func (s *testsession) ChannelMessageEdit(channelID, messageID, content string) (*discordgo.Message, error) {
	panic("ChannelMessageEdit not implemented")
}

func (s *testsession) ChannelMessageEditEmbeds(channelID, messageID string, embeds []*discordgo.MessageEmbed) (
	*discordgo.Message, error) {
	panic("ChannelMessageEditEmbeds not implemented")
}

func (s *testsession) ChannelMessageDelete(channelID, messageID string) error {
	panic("ChannelMessageDelete not implemented")
}

func (s *testsession) ChannelMessagesBulkDelete(channelID string, messages []string) error {
	panic("ChannelMessagesBulkDelete not implemented")
}

func (s *testsession) ChannelMessagePin(channelID, messageID string) error {
	panic("ChannelMessagePin not implemented")
}

func (s *testsession) ChannelMessageUnpin(channelID, messageID string) error {
	panic("ChannelMessageUnpin not implemented")
}

func (s *testsession) UpdateGameStatus(idle int, name string) error {
	panic("UpdateGameStatus not implemented")
}

func (s *testsession) UpdateStreamingStatus(idle int, name string, url string) error {
	panic("UpdateStreamingStatus not implemented")
}

func (s *testsession) UpdateListeningStatus(name string) error {
	panic("UpdateListeningStatus not implemented")
}

func (s *testsession) Open() error {
	panic("Open not implemented")
}

func (s *testsession) Close() error {
	panic("Close not implemented")
}
