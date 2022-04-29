package replies_test

import (
	"strings"
	"testing"
	"time"

	"github.com/emseers/eelbot/replies"
	"github.com/stretchr/testify/require"
)

func TestCaps(t *testing.T) {
	s := newTestSession()
	f := replies.CapsReply(100, 0, 0, 5).Eval

	require.False(t, f(s, newMsgCreate("HOME", testChannelID)))
	require.Nil(t, s.messages[testChannelID])

	require.True(t, f(s, newMsgCreate("HELLO", testChannelID)))
	time.Sleep(10 * time.Millisecond) // Small delay for goroutine to finish the write.
	require.NotEmpty(t, strings.TrimSpace(s.messages[testChannelID].String()))
}
