package commands_test

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
)

const (
	testChannelID      = "123456789"
	testJoke1          = "A steak pun is a rare medium well done."
	testJoke2          = "What do you call a fish with no eyes?"
	testJoke2Punchline = "Fsh"
	testFile1          = "testdata/file1.png"
	testFile2          = "testdata/file2.png"
)

var (
	db *sql.DB
)

func setup() (err error) {
	db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		return
	}

	_, err = db.Exec(fmt.Sprintf(`
CREATE TABLE "jokes" (
  "id"        INTEGER NOT NULL,
  "text"      TEXT NOT NULL,
  "punchline" TEXT,
  PRIMARY KEY("id" AUTOINCREMENT)
);
CREATE TABLE "images" (
	"id"   INTEGER NOT NULL,
	"path" TEXT NOT NULL,
	PRIMARY KEY("id" AUTOINCREMENT)
  );
CREATE TABLE "taunts" (
  "id"   INTEGER NOT NULL,
  "path" TEXT NOT NULL,
  PRIMARY KEY("id" AUTOINCREMENT)
);

INSERT INTO "jokes" ("text", "punchline") VALUES
  (%q, null),
  (%q, %q);
INSERT INTO "images" ("path") VALUES
  (%[4]q),
  (%[5]q);
INSERT INTO "taunts" ("path") VALUES
  (%[4]q),
  (%[5]q);
`, testJoke1, testJoke2, testJoke2Punchline, testFile1, testFile2))
	return
}

func TestMain(m *testing.M) {
	if err := setup(); err != nil {
		log.Fatalln(err)
	}
	os.Exit(m.Run())
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

func (s *testsession) AddHandler(handler any) func() {
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
