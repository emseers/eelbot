package commands_test

import (
	"strings"
	"testing"

	"github.com/emseers/eelbot/commands"
	"github.com/stretchr/testify/require"
)

func TestPlay(t *testing.T) {
	s := newTestSession()
	f := commands.PlayCommand().Eval

	const name = "Doom II: Hell on Earth"
	require.NoError(t, f(s, newMsgCreate("/play "+name, testChannelID), strings.Split(name, " ")))
	require.Equal(t, name, s.statusPlay)

	require.NoError(t, f(s, newMsgCreate("", testChannelID), []string{}))
	require.Equal(t, "", s.statusPlay)
}
