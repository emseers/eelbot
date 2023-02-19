package commands_test

import (
	"math/rand"
	"strings"
	"testing"

	"github.com/emseers/eelbot/commands"
	"github.com/stretchr/testify/require"
)

func TestFlip(t *testing.T) {
	s := newTestSession()
	f := commands.FlipCommand().Eval
	commands.Rand = rand.New(rand.NewSource((42)))

	const numFlips = 100000
	for i := 0; i < numFlips; i++ {
		require.NoError(t, f(s, newMsgCreate("", testChannelID), []string{}))
	}

	results := strings.Split(strings.TrimSpace(s.messages[testChannelID].String()), "\n")
	require.Len(t, results, numFlips)

	count := map[string]int{}
	for _, r := range results {
		count[r]++
	}

	require.Equal(t, 50147, count["Heads"])
	require.Equal(t, 49832, count["Tails"])
	require.Equal(t, 21, count["Landed on edge"])
}
