package commands_test

import (
	"strings"
	"testing"

	"github.com/emseers/eelbot/commands"
	"github.com/stretchr/testify/require"
)

func TestChannel(t *testing.T) {
	s := newTestSession()
	f := commands.ChannelCommand().Eval

	require.NoError(t, f(s, newMsgCreate("", testChannelID), []string{}))
	require.Equal(t, testChannelID, strings.TrimSpace(s.messages[testChannelID].String()))
}
