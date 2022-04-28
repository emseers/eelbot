package eelbot_test

import (
	"strings"
	"testing"

	"github.com/emseers/eelbot"
	"github.com/emseers/eelbot/commands"
	"github.com/stretchr/testify/require"
)

func TestCommand(t *testing.T) {
	s := newTestSession()
	bot := eelbot.New(s)
	bot.RegisterCommand("channel", *commands.ChannelCommand())
	bot.RegisterCommand("roll", *commands.RollCommand())
	bot.RegisterCommand("badjoke", *commands.JokeCommand(nil, 0))

	c := *commands.PlayCommand()
	c.MinArgs = -1
	bot.RegisterCommand("play", c)

	require.NoError(t, bot.Start())

	s.send(newMsg("/channel", testChannelID, "", ""))
	require.Equal(t, testChannelID, strings.TrimSpace(s.messages[testChannelID].String()))

	s.messages[testChannelID].Reset()
	s.send(newMsg("/channel 1", testChannelID, "", ""))
	require.Equal(t, "Error: channel requires at most 0 arguments",
		strings.TrimSpace(s.messages[testChannelID].String()))

	s.messages[testChannelID].Reset()
	s.send(newMsg("/roll 1 1", testChannelID, "", ""))
	require.Equal(t, "1", strings.TrimSpace(s.messages[testChannelID].String()))

	s.messages[testChannelID].Reset()
	s.send(newMsg("/roll", testChannelID, "", ""))
	require.Equal(t, "Error: roll requires at least 1 arguments", strings.TrimSpace(s.messages[testChannelID].String()))

	s.messages[testChannelID].Reset()
	s.send(newMsg("/roll -9223372036854775807 9223372036854775807", testChannelID, "", "")) // Should cause a panic.
	require.Equal(t, "Error: invalid argument to Int63n", strings.TrimSpace(s.messages[testChannelID].String()))

	s.messages[testChannelID].Reset()
	s.send(newMsg("/help", testChannelID, "", ""))
	require.Equal(t, `Available commands:
`+"```"+`
/badjoke: Posts a joke.
/channel: Posts the current channel ID.
/help   : Displays the summary of available commands or details for a specific command.
/play   : Plays a game.
/roll   : Rolls a die.
`+"```", strings.TrimSpace(s.messages[testChannelID].String()))

	s.messages[testChannelID].Reset()
	s.send(newMsg("/help play", testChannelID, "", ""))
	require.Equal(t, "```"+`
/play [ARGS...]

Plays a game.
`+"```", strings.TrimSpace(s.messages[testChannelID].String()))

	s.messages[testChannelID].Reset()
	s.send(newMsg("/help roll", testChannelID, "", ""))
	require.Equal(t, "```"+`
/roll ARG1 [ARG2]

Rolls a die.
`+"```", strings.TrimSpace(s.messages[testChannelID].String()))

	s.messages[testChannelID].Reset()
	s.send(newMsg("/help badjoke", testChannelID, "", ""))
	require.Equal(t, "```"+`
/badjoke NUM

Posts a joke from the database. NUM can either be a valid joke number from the database, or "me" for a random joke.

Examples:
  /badjoke me
  /badjoke 42
`+"```", strings.TrimSpace(s.messages[testChannelID].String()))

	s.messages[testChannelID].Reset()
	s.send(newMsg("/help fish", testChannelID, "", ""))
	require.Equal(t, "Error: unknown command: fish", strings.TrimSpace(s.messages[testChannelID].String()))

	require.NoError(t, bot.Stop())
}
