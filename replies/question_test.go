package replies_test

import (
	"strings"
	"testing"

	"github.com/emseers/eelbot/replies"
	"github.com/stretchr/testify/require"
)

func TestQuestion(t *testing.T) {
	s := newTestSession()
	f := replies.QuestionReply(100).Eval

	require.False(t, f(s, newMsgCreate("¿que pasa?", testChannelID)))
	require.Nil(t, s.messages[testChannelID])

	require.True(t, f(s, newMsgCreate("?", testChannelID)))
	require.Equal(t, "Don't questionmark me", strings.TrimSpace(s.messages[testChannelID].String()))
}
