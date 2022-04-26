package commands_test

import (
	"strings"
	"testing"

	"github.com/emseers/eelbot/commands"
	"github.com/stretchr/testify/require"
)

func TestSayChan(t *testing.T) {
	s := newTestSession()
	f := commands.SayChanCommand().Eval

	const chanID = "42"
	const msg = "Time is an illusion, lunchtime doubly so."
	require.NoError(t, f(s, newMsgCreate("", testChannelID), append([]string{chanID}, strings.Split(msg, " ")...)))
	require.Equal(t, msg, strings.TrimSpace(s.messages[chanID].String()))
}
