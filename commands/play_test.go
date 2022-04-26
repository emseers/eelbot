package commands_test

import (
	"testing"

	"github.com/emseers/eelbot/commands"
	"github.com/stretchr/testify/require"
)

func TestPlay(t *testing.T) {
	s := newTestSession()
	f := commands.PlayCommand().Eval

	require.NoError(t, f(s, newMsgCreate("", testChannelID), []string{"Doom"}))
	require.Equal(t, "Doom", s.statusPlay)

	require.NoError(t, f(s, newMsgCreate("", testChannelID), []string{}))
	require.Equal(t, "", s.statusPlay)
}
