package replies_test

import (
	"strings"
	"testing"
	"time"

	"github.com/emseers/eelbot/replies"
	"github.com/stretchr/testify/require"
)

func TestQuestion(t *testing.T) {
	s := newTestSession()
	f := replies.QuestionReply(100, 0, 0).Eval

	require.False(t, f(s, newMsgCreate("¿que pasa?", testChannelID)))
	require.Nil(t, s.messages[testChannelID])

	require.True(t, f(s, newMsgCreate("?", testChannelID)))
	time.Sleep(10 * time.Millisecond) // Small delay for goroutine to finish the write.
	require.Equal(t, "Don't questionmark me", strings.TrimSpace(s.messages[testChannelID].String()))
}
