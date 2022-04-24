package commands_test

import (
	"encoding/base64"
	"testing"

	"github.com/emseers/eelbot/commands"
	"github.com/stretchr/testify/require"
)

func TestImage(t *testing.T) {
	s := newTestSession()
	f := commands.ImageCommand(db).Eval

	require.NoError(t, f(s, newMsgCreate("", testChannelID), []string{"me"}))
	require.Len(t, s.files[testChannelID], 1)

	f1, err1 := base64.StdEncoding.DecodeString(testFile1)
	f2, err2 := base64.StdEncoding.DecodeString(testFile2)
	require.NoError(t, err1)
	require.NoError(t, err2)

	require.NoError(t, f(s, newMsgCreate("", testChannelID), []string{"1"}))
	require.Equal(t, f1, s.files[testChannelID][testFileName1])

	require.NoError(t, f(s, newMsgCreate("", testChannelID), []string{"2"}))
	require.Equal(t, f2, s.files[testChannelID][testFileName2])

	require.EqualError(t, f(s, newMsgCreate("", testChannelID), []string{"b"}), "unknown directive: b")
}
