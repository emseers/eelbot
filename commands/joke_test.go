package commands_test

import (
	"strings"
	"testing"

	"github.com/emseers/eelbot/commands"
	"github.com/stretchr/testify/require"
)

func TestJoke(t *testing.T) {
	s := newTestSession()
	f := commands.JokeCommand(db, 0).Eval

	require.NoError(t, f(s, newMsgCreate("", testChannelID), []string{"me"}))
	require.NotEmpty(t, s.messages[testChannelID].String())

	s.messages[testChannelID].Reset()
	require.NoError(t, f(s, newMsgCreate("", testChannelID), []string{"1"}))
	require.Equal(t, testJoke1, strings.TrimSpace(s.messages[testChannelID].String()))

	s.messages[testChannelID].Reset()
	require.NoError(t, f(s, newMsgCreate("", testChannelID), []string{"2"}))
	require.Equal(t, testJoke2+"\n"+testJoke2Punchline, strings.TrimSpace(s.messages[testChannelID].String()))

	require.EqualError(t, f(s, newMsgCreate("", testChannelID), []string{"b"}), "unknown directive: b")
}
