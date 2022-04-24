package commands_test

import (
	"strings"
	"testing"

	"github.com/emseers/eelbot/commands"
	"github.com/stretchr/testify/require"
)

func TestSay(t *testing.T) {
	s := newTestSession()
	f := commands.SayCommand().Eval

	const msg = "Time is an illusion, lunchtime doubly so."
	require.NoError(t, f(s, newMsgCreate("", testChannelID), strings.Split(msg, " ")))
	require.Equal(t, msg, strings.TrimSpace(s.messages[testChannelID].String()))
}
