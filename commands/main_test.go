package commands_test

import (
	"bytes"
	"database/sql"
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
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	initdbScriptsDir   = "../initdb"
	testChannelID      = "123456789"
	testJoke1          = "A steak pun is a rare medium well done."
	testJoke2          = "What do you call a fish with no eyes?"
	testJoke2Punchline = "Fsh"
	testFileName1      = "file1.png"
	testFile1          = "some_image_file"
	testFileName2      = "file2.png"
	testFile2          = "some_other_file"
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

	_, err = db.Exec("INSERT INTO jokes (id, text, punchline) VALUES (1, $1, null), (2, $2, $3);",
		testJoke1, testJoke2, testJoke2Punchline)
	if err != nil {
		return
	}

	_, err = db.Exec("INSERT INTO images (id, name, file) VALUES (1, $1, $2), (2, $3, $4);",
		testFileName1, []byte(testFile1), testFileName2, []byte(testFile2))
	if err != nil {
		return
	}

	_, err = db.Exec("INSERT INTO taunts (id, name, file) VALUES (1, $1, $2), (2, $3, $4);",
		testFileName1, []byte(testFile1), testFileName2, []byte(testFile2))
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
		files:    make(map[string]map[string][]byte),
	}
}

type testsession struct {
	counter      int
	messages     map[string]*strings.Builder
	files        map[string]map[string][]byte
	statusListen string
	statusPlay   string
}

func (s *testsession) AddHandler(any) func() {
	return func() {}
}

func (s *testsession) ChannelTyping(channelID string) (err error) {
	return nil
}

func (s *testsession) ChannelMessageSend(channelID, content string) (*discordgo.Message, error) {
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

func (s *testsession) ChannelFileSend(channelID, name string, r io.Reader) (*discordgo.Message, error) {
	b, ok := s.files[channelID]
	if !ok {
		b = make(map[string][]byte)
		s.files[channelID] = b
	}
	f, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	b[name] = f
	s.counter++
	return newMsg(name, channelID, strconv.Itoa(s.counter)), nil
}

func (s *testsession) UpdateGameStatus(idle int, name string) error {
	s.statusPlay = name
	return nil
}

func (s *testsession) UpdateListeningStatus(name string) error {
	s.statusListen = name
	return nil
}

func (s *testsession) ChannelMessageDelete(channelID, messageID string) error {
	return nil
}

// The following functions exist for the purpose of satisfying the eelbot.Session interface.

func (s *testsession) Channel(string) (*discordgo.Channel, error) {
	panic("Channel not implemented")
}

func (s *testsession) ChannelMessage(string, string) (*discordgo.Message, error) {
	panic("ChannelMessage not implemented")
}

func (s *testsession) ChannelMessages(string, int, string, string, string) ([]*discordgo.Message, error) {
	panic("ChannelMessages not implemented")
}

func (s *testsession) ChannelMessagesPinned(string) ([]*discordgo.Message, error) {
	panic("ChannelMessagesPinned not implemented")
}

func (s *testsession) ChannelMessageSendEmbeds(string, []*discordgo.MessageEmbed) (*discordgo.Message, error) {
	panic("ChannelMessageSendEmbeds not implemented")
}

func (s *testsession) ChannelMessageSendTTS(string, string) (*discordgo.Message, error) {
	panic("ChannelMessageSendTTS not implemented")
}

func (s *testsession) ChannelFileSendWithMessage(string, string, string, io.Reader) (*discordgo.Message, error) {
	panic("ChannelFileSendWithMessage not implemented")
}

func (s *testsession) ChannelMessageEdit(string, string, string) (*discordgo.Message, error) {
	panic("ChannelMessageEdit not implemented")
}

func (s *testsession) ChannelMessageEditEmbeds(string, string, []*discordgo.MessageEmbed) (*discordgo.Message, error) {
	panic("ChannelMessageEditEmbeds not implemented")
}

func (s *testsession) ChannelMessagesBulkDelete(string, []string) error {
	panic("ChannelMessagesBulkDelete not implemented")
}

func (s *testsession) ChannelMessagePin(string, string) error {
	panic("ChannelMessagePin not implemented")
}

func (s *testsession) ChannelMessageUnpin(string, string) error {
	panic("ChannelMessageUnpin not implemented")
}

func (s *testsession) UpdateStreamingStatus(int, string, string) error {
	panic("UpdateStreamingStatus not implemented")
}

func (s *testsession) Open() error {
	panic("Open not implemented")
}

func (s *testsession) Close() error {
	panic("Close not implemented")
}
