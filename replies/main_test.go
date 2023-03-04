package replies_test

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/emseers/eelbot/helpers"
)

const (
	initdbScriptsDir     = "../initdb"
	incorrectFormatTable = "incorrect_format"
	sampleTable          = "sample"
	testChannelID        = "123456789"
)

var (
	db *sql.DB
)

func setup(url string) (err error) {
	db, err = sql.Open("pgx", url)
	if err != nil {
		return
	}

	// Load initdb scripts and run them to setup the schema.
	init := new(bytes.Buffer)
	err = filepath.WalkDir(initdbScriptsDir, func(path string, f fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		var b []byte
		if b, err = os.ReadFile(path); err != nil {
			return err
		}
		init.Write(b)
		init.WriteByte('\n')
		return nil
	})
	if err != nil {
		return
	}

	_, err = db.Exec(init.String())
	if err != nil {
		return
	}

	_, err = db.Exec(fmt.Sprintf("CREATE TABLE reply.%s (eye_dee integer PRIMARY KEY, file bytea NOT NULL);",
		incorrectFormatTable))
	if err != nil {
		return
	}

	_, err = db.Exec(fmt.Sprintf("CREATE TABLE reply.%s (id integer PRIMARY KEY, text text NOT NULL);", sampleTable))
	if err != nil {
		return
	}

	_, err = db.Exec(fmt.Sprintf("INSERT INTO reply.%s (id, text) VALUES (1, $1);", sampleTable), "hi")
	return
}

func TestMain(m *testing.M) {
	url, close, err := helpers.StartPostgreSQL()
	if err != nil {
		log.Fatalln(err)
	}

	if err = setup(url); err != nil {
		close()
		log.Fatalln(err)
	}

	code := m.Run()

	if db != nil {
		_ = db.Close()
	}
	close()

	os.Exit(code)
}

// Creates a Message with the provided arguments.
//
//	0: Content
//	1: ChannelID
//	2: MessageID
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

func (s *testsession) AddHandler(any) func() {
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

func (s *testsession) Open() error {
	panic("Open not implemented")
}

func (s *testsession) Close() error {
	panic("Close not implemented")
}
