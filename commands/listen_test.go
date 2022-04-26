package commands_test

import (
	"testing"

	"github.com/emseers/eelbot/commands"
	"github.com/stretchr/testify/require"
)

func TestListen(t *testing.T) {
	s := newTestSession()
	f := commands.ListenCommand().Eval

	require.NoError(t, f(s, newMsgCreate("", testChannelID), []string{"Money"}))
	require.Equal(t, "Money", s.statusListen)

	require.NoError(t, f(s, newMsgCreate("", testChannelID), []string{}))
	require.Equal(t, "", s.statusListen)
}
