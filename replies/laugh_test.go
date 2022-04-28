package replies_test

import (
	"strings"
	"testing"
	"time"

	"github.com/emseers/eelbot/replies"
	"github.com/stretchr/testify/require"
)

func TestLaugh(t *testing.T) {
	s := newTestSession()
	f := replies.LaughReply(100, 0, 0).Eval

	require.False(t, f(s, newMsgCreate("lawl", testChannelID)))
	require.Nil(t, s.messages[testChannelID])

	require.True(t, f(s, newMsgCreate("lol", testChannelID)))
	time.Sleep(10 * time.Millisecond) // Small delay for goroutine to finish the write.
	require.Equal(t, "lol", strings.TrimSpace(s.messages[testChannelID].String()))
}
