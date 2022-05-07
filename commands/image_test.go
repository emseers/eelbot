package commands_test

import (
	"path"
	"testing"
	"time"

	"github.com/emseers/eelbot/commands"
	"github.com/stretchr/testify/require"
)

func TestImage(t *testing.T) {
	s := newTestSession()
	f := commands.ImageCommand(db, time.Second).Eval

	require.NoError(t, f(s, newMsgCreate("", testChannelID), []string{"me"}))
	require.Len(t, s.files[testChannelID], 1)

	require.NoError(t, f(s, newMsgCreate("", testChannelID), []string{"1"}))
	require.Equal(t, []byte(testFile1), s.files[testChannelID][path.Base(testFileName1)])

	require.NoError(t, f(s, newMsgCreate("", testChannelID), []string{"2"}))
	require.Equal(t, []byte(testFile2), s.files[testChannelID][path.Base(testFileName2)])

	require.EqualError(t, f(s, newMsgCreate("", testChannelID), []string{"b"}), "unknown directive: b")
}
