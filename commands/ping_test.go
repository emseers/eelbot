package commands_test

import (
	"strings"
	"testing"

	"github.com/emseers/eelbot/commands"
	"github.com/stretchr/testify/require"
)

func TestPing(t *testing.T) {
	s := newTestSession()
	f := commands.PingCommand().Eval

	require.NoError(t, f(s, newMsgCreate("", testChannelID), []string{}))
	require.Equal(t, "Pong", strings.TrimSpace(s.messages[testChannelID].String()))
}
