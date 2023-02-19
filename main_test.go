package eelbot_test

import (
	"io"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	testUserID    = "1011089811111610"
	testChannelID = "123456789"
)

// Creates a Message with the provided arguments.
//
//	0: Content
//	1: ChannelID
//	2: MessageID
//	3: AuthorID
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
	if len(args) > 3 {
		m.Author = &discordgo.User{ID: args[3]}
	}
	return m
}

func newTestSession() *testsession {
	s := &testsession{
		messages: make(map[string]*strings.Builder),
		sess:     &discordgo.Session{State: discordgo.NewState()},
	}
	s.sess.State.User = &discordgo.User{ID: testUserID}
	return s
}

type testsession struct {
	counter  int
	sess     *discordgo.Session
	handler  func(s *discordgo.Session, m *discordgo.MessageCreate)
	messages map[string]*strings.Builder
}

func (s *testsession) send(m *discordgo.Message) {
	s.handler(s.sess, &discordgo.MessageCreate{Message: m})
}

func (s *testsession) AddHandler(handler any) func() {
	s.handler = handler.(func(s *discordgo.Session, m *discordgo.MessageCreate))
	return func() {}
}

func (s *testsession) ChannelTyping(string, ...discordgo.RequestOption) (err error) {
	return nil
}

func (s *testsession) ChannelMessageSend(channelID, content string,
	options ...discordgo.RequestOption) (*discordgo.Message, error) {
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

func (s *testsession) Open() error {
	return nil
}

func (s *testsession) Close() error {
	return nil
}

// The following functions exist for the purpose of satisfying the eelbot.Session interface.

func (s *testsession) Channel(string, ...discordgo.RequestOption) (*discordgo.Channel, error) {
	panic("Channel not implemented")
}

func (s *testsession) ChannelMessage(string, string, ...discordgo.RequestOption) (*discordgo.Message, error) {
	panic("ChannelMessage not implemented")
}

func (s *testsession) ChannelMessages(string, int, string, string, string,
	...discordgo.RequestOption) ([]*discordgo.Message, error) {
	panic("ChannelMessages not implemented")
}

func (s *testsession) ChannelMessagesPinned(string, ...discordgo.RequestOption) ([]*discordgo.Message, error) {
	panic("ChannelMessagesPinned not implemented")
}

func (s *testsession) ChannelMessageSendEmbeds(string, []*discordgo.MessageEmbed,
	...discordgo.RequestOption) (*discordgo.Message, error) {
	panic("ChannelMessageSendEmbeds not implemented")
}

func (s *testsession) ChannelMessageSendTTS(string, string, ...discordgo.RequestOption) (*discordgo.Message, error) {
	panic("ChannelMessageSendTTS not implemented")
}

func (s *testsession) ChannelFileSend(string, string, io.Reader,
	...discordgo.RequestOption) (*discordgo.Message, error) {
	panic("ChannelFileSend not implemented")
}

func (s *testsession) ChannelFileSendWithMessage(string, string, string, io.Reader,
	...discordgo.RequestOption) (*discordgo.Message, error) {
	panic("ChannelFileSendWithMessage not implemented")
}

func (s *testsession) ChannelMessageEdit(string, string, string,
	...discordgo.RequestOption) (*discordgo.Message, error) {
	panic("ChannelMessageEdit not implemented")
}

func (s *testsession) ChannelMessageEditEmbeds(string, string, []*discordgo.MessageEmbed,
	...discordgo.RequestOption) (*discordgo.Message, error) {
	panic("ChannelMessageEditEmbeds not implemented")
}

func (s *testsession) ChannelMessageDelete(string, string, ...discordgo.RequestOption) error {
	panic("ChannelMessageDelete not implemented")
}

func (s *testsession) ChannelMessagesBulkDelete(string, []string, ...discordgo.RequestOption) error {
	panic("ChannelMessagesBulkDelete not implemented")
}

func (s *testsession) ChannelMessagePin(string, string, ...discordgo.RequestOption) error {
	panic("ChannelMessagePin not implemented")
}

func (s *testsession) ChannelMessageUnpin(string, string, ...discordgo.RequestOption) error {
	panic("ChannelMessageUnpin not implemented")
}

func (s *testsession) UpdateGameStatus(int, string) error {
	panic("UpdateGameStatus not implemented")
}

func (s *testsession) UpdateStreamingStatus(int, string, string) error {
	panic("UpdateStreamingStatus not implemented")
}

func (s *testsession) UpdateListeningStatus(string) error {
	panic("UpdateListeningStatus not implemented")
}
