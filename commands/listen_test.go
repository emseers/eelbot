package commands_test

import (
	"strings"
	"testing"

	"github.com/emseers/eelbot/commands"
	"github.com/stretchr/testify/require"
)

func TestListen(t *testing.T) {
	s := newTestSession()
	f := commands.ListenCommand().Eval

	const name = "Hot in Herre"
	require.NoError(t, f(s, newMsgCreate("/listen "+name, testChannelID), strings.Split(name, " ")))
	require.Equal(t, name, s.statusListen)

	require.NoError(t, f(s, newMsgCreate("", testChannelID), []string{}))
	require.Equal(t, "", s.statusListen)
}
