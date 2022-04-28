package replies_test

import (
	"strings"
	"testing"
	"time"

	"github.com/emseers/eelbot"
	"github.com/emseers/eelbot/replies"
	"github.com/stretchr/testify/require"
)

var (
	cfg = map[string]any{
		"caps": map[string]any{
			"enable":  true,
			"min_len": 5,
			"percent": 17,
			"timeout": 120,
		},
		"hello": map[string]any{
			"enable":  true,
			"percent": 33,
			"timeout": 600,
		},
		"goodbye": map[string]any{
			"enable":  true,
			"percent": 33,
			"timeout": 600,
		},
		"laugh": map[string]any{
			"enable":  true,
			"percent": 17,
			"timeout": 10,
		},
		"question": map[string]any{
			"enable":  true,
			"percent": 17,
			"timeout": 10,
		},
	}
)

func TestRegister(t *testing.T) {
	bot := eelbot.New(newTestSession())
	require.NoError(t, replies.Register(bot, cfg))

	// Should be fault tolerant with invalid keys.
	cfg["caps"].(map[string]any)["min_len"] = "foo"
	cfg["hello"].(map[string]any)["percent"] = "bar"
	require.NoError(t, replies.Register(bot, cfg))
}

func TestReply(t *testing.T) {
	s := newTestSession()
	f := replies.LaughReply(100, 50*time.Millisecond, 150*time.Millisecond).Eval

	require.True(t, f(s, newMsgCreate("lol", testChannelID)))
	time.Sleep(10 * time.Millisecond) // Not enough delay for goroutine to finish the write.
	require.Nil(t, s.messages[testChannelID])
	time.Sleep(200 * time.Millisecond) // Enough delay for goroutine to finish the write.
	require.Equal(t, "lol", strings.TrimSpace(s.messages[testChannelID].String()))

	s.messages[testChannelID].Reset()
	f = replies.LaughReply(100, 50*time.Millisecond, 0).Eval // maxDelay is less than minDelay.
	require.True(t, f(s, newMsgCreate("lol", testChannelID)))
	time.Sleep(10 * time.Millisecond) // Not enough delay for goroutine to finish the write.
	require.Equal(t, "", strings.TrimSpace(s.messages[testChannelID].String()))
	time.Sleep(100 * time.Millisecond) // Enough delay for goroutine to finish the write.
	require.Equal(t, "lol", strings.TrimSpace(s.messages[testChannelID].String()))
}
