package eelbot_test

import (
	"strings"
	"testing"
	"time"

	"github.com/emseers/eelbot"
	"github.com/emseers/eelbot/replies"
	"github.com/stretchr/testify/require"
)

func TestReply(t *testing.T) {
	s := newTestSession()
	bot := eelbot.New(s)
	r := replies.LaughReply(100, 0, 0)
	r.Timeout = 500 * time.Millisecond
	bot.RegisterReply(*r)

	s.send(newMsg("lol", testChannelID, "", testUserID))
	require.Nil(t, s.messages[testChannelID])

	for i := 0; i < 20; i++ {
		s.send(newMsg("lol", testChannelID, "", ""))
	}
	time.Sleep(10 * time.Millisecond) // Small delay for goroutines to finish writes.
	require.Equal(t, "lol", strings.TrimSpace(s.messages[testChannelID].String()))

	s.messages[testChannelID].Reset()
	time.Sleep(500 * time.Millisecond)

	for i := 0; i < 20; i++ {
		s.send(newMsg("lol", testChannelID, "", ""))
	}
	time.Sleep(10 * time.Millisecond) // Small delay for goroutine to finish the write.
	require.Equal(t, "lol", strings.TrimSpace(s.messages[testChannelID].String()))
}
