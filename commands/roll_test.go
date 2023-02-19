package commands_test

import (
	"math/rand"
	"strconv"
	"strings"
	"testing"

	"github.com/emseers/eelbot/commands"
	"github.com/stretchr/testify/require"
)

func TestRoll(t *testing.T) {
	s := newTestSession()
	f := commands.RollCommand().Eval
	commands.Rand = rand.New(rand.NewSource((42)))

	const numRolls = 100000
	for i := 0; i < numRolls; i++ {
		require.NoError(t, f(s, newMsgCreate("", testChannelID), []string{"999"}))
	}

	results := strings.Split(strings.TrimSpace(s.messages[testChannelID].String()), "\n")
	require.Len(t, results, numRolls)

	count := map[string]int{}
	for _, r := range results {
		count[r]++
	}

	require.Len(t, count, 1000)
	for i := 0; i <= 999; i++ {
		require.Greater(t, count[strconv.Itoa(i)], 0, "number of rolls that resulted in %d", i)
	}

	s.messages[testChannelID].Reset()
	for i := 0; i < numRolls; i++ {
		require.NoError(t, f(s, newMsgCreate("", testChannelID), []string{"1999", "1000"}))
	}

	results = strings.Split(strings.TrimSpace(s.messages[testChannelID].String()), "\n")
	require.Len(t, results, numRolls)

	count = map[string]int{}
	for _, r := range results {
		count[r]++
	}

	require.Len(t, count, 1000)
	for i := 1000; i <= 1999; i++ {
		require.Greater(t, count[strconv.Itoa(i)], 0, "number of rolls that resulted in %d", i)
	}

	require.Error(t, f(s, newMsgCreate("", testChannelID), []string{"b"}))
	require.Error(t, f(s, newMsgCreate("", testChannelID), []string{"0", "b"}))
	require.Error(t, f(s, newMsgCreate("", testChannelID), []string{"a", "b"}))
}
