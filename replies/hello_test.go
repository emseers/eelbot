package replies_test

import (
	"strings"
	"testing"
	"time"

	"github.com/emseers/eelbot/replies"
	"github.com/stretchr/testify/require"
)

func TestHello(t *testing.T) {
	s := newTestSession()
	f := replies.HelloReply(100, 0, 0).Eval

	require.False(t, f(s, newMsgCreate("Goodbye", testChannelID)))
	require.Nil(t, s.messages[testChannelID])

	require.True(t, f(s, newMsgCreate("Hello", testChannelID)))
	time.Sleep(10 * time.Millisecond) // Small delay for goroutine to finish the write.
	require.NotEmpty(t, strings.TrimSpace(s.messages[testChannelID].String()))
}
