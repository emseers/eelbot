package commands_test

import (
	"database/sql"
	"encoding/base64"
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
	testFileName1      = "file1.png"
	testFileName2      = "file2.png"

	testFile1 = "iVBORw0KGgoAAAANSUhEUgAAAAIAAAACCAIAAAD91JpzAAABhGlDQ1BJQ0MgcHJvZmlsZQAAKJF9kT1Iw1AUhU/TSkUqDnYQcchQ" +
		"O1kQFXHUKhShQqgVWnUweekfNGlIUlwcBdeCgz+LVQcXZ10dXAVB8AfE0clJ0UVKvC8ptIjxwuN9nHfP4b37AKFZZZoVGgc03TYzqaSYy6+K" +
		"4VeEMIQAYojLzDLmJCkN3/q6p06quwTP8u/7s/rVgsWAgEg8ywzTJt4gnt60Dc77xFFWllXic+Ixky5I/Mh1xeM3ziWXBZ4ZNbOZeeIosVjq" +
		"YqWLWdnUiKeIY6qmU76Q81jlvMVZq9ZZ+578hZGCvrLMdVojSGERS5AgQkEdFVRhI0G7ToqFDJ0nffzDrl8il0KuChg5FlCDBtn1g//B79la" +
		"xckJLymSBHpeHOdjFAjvAq2G43wfO07rBAg+A1d6x19rAjOfpDc6WuwIGNgGLq47mrIHXO4AQ0+GbMquFKQlFIvA+xl9Ux4YvAX61ry5tc9x" +
		"+gBkaVbpG+DgEIiXKHvd59293XP7t6c9vx9wdHKmBJh2ogAAAAlwSFlzAAAuIwAALiMBeKU/dgAAAAd0SU1FB+YEGAEnEd1Tcy8AAAAZdEVY" +
		"dENvbW1lbnQAQ3JlYXRlZCB3aXRoIEdJTVBXgQ4XAAAAFUlEQVQI1wXBAQEAAACAEP9PF1CpMCnkBftjnTYAAAAAAElFTkSuQmCC"
	testFile2 = "iVBORw0KGgoAAAANSUhEUgAAAAIAAAACCAIAAAD91JpzAAABhGlDQ1BJQ0MgcHJvZmlsZQAAKJF9kT1Iw1AUhU/TSkUqDnYQcchQ" +
		"O1kQFXHUKhShQqgVWnUweekfNGlIUlwcBdeCgz+LVQcXZ10dXAVB8AfE0clJ0UVKvC8ptIjxwuN9nHfP4b37AKFZZZoVGgc03TYzqaSYy6+K" +
		"4VeEMIQAYojLzDLmJCkN3/q6p06quwTP8u/7s/rVgsWAgEg8ywzTJt4gnt60Dc77xFFWllXic+Ixky5I/Mh1xeM3ziWXBZ4ZNbOZeeIosVjq" +
		"YqWLWdnUiKeIY6qmU76Q81jlvMVZq9ZZ+578hZGCvrLMdVojSGERS5AgQkEdFVRhI0G7ToqFDJ0nffzDrl8il0KuChg5FlCDBtn1g//B79la" +
		"xckJLymSBHpeHOdjFAjvAq2G43wfO07rBAg+A1d6x19rAjOfpDc6WuwIGNgGLq47mrIHXO4AQ0+GbMquFKQlFIvA+xl9Ux4YvAX61ry5tc9x" +
		"+gBkaVbpG+DgEIiXKHvd59293XP7t6c9vx9wdHKmBJh2ogAAAAlwSFlzAAAuIwAALiMBeKU/dgAAAAd0SU1FB+YEGAEnJmXu1iAAAAAZdEVY" +
		"dENvbW1lbnQAQ3JlYXRlZCB3aXRoIEdJTVBXgQ4XAAAAEklEQVQI12P4//8/AwT8//8fACnkBft7DmIIAAAAAElFTkSuQmCC"
)

var (
	db *sql.DB
)

func setup() (err error) {
	db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		return
	}

	var f1, f2 []byte
	if f1, err = base64.StdEncoding.DecodeString(testFile1); err != nil {
		return
	}
	if f2, err = base64.StdEncoding.DecodeString(testFile2); err != nil {
		return
	}
	if err = os.WriteFile(testFileName1, f1, 0644); err != nil {
		return
	}
	if err = os.WriteFile(testFileName2, f2, 0644); err != nil {
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
`, testJoke1, testJoke2, testJoke2Punchline, testFileName1, testFileName2))
	return
}

func teardown() {
	_ = os.Remove(testFileName1)
	_ = os.Remove(testFileName2)
}

func TestMain(m *testing.M) {
	if err := setup(); err != nil {
		log.Fatalln(err)
	}
	code := m.Run()
	teardown()
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
	return nil
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
	return nil
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
	s.statusPlay = name
	return nil
}

func (s *testsession) UpdateStreamingStatus(idle int, name string, url string) error {
	panic("UpdateStreamingStatus not implemented")
}

func (s *testsession) UpdateListeningStatus(name string) error {
	s.statusListen = name
	return nil
}

func (s *testsession) Open() error {
	panic("Open not implemented")
}

func (s *testsession) Close() error {
	panic("Close not implemented")
}
