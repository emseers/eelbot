package commands_test

import (
	"testing"

	"github.com/emseers/eelbot"
	"github.com/emseers/eelbot/commands"
	"github.com/stretchr/testify/require"
	"gopkg.in/ini.v1"
)

const (
	cfg = `[Commands]
badjoke_enable = true
badjoke_delay = 3
eel_enable = true
taunt_enable = true
channel_enable = true
flip_enable = true
listen_enable = true
play_enable = true
ping_enable = true
roll_enable = true
say_enable = true
saychan_enable = true
`
	cfgJokeOnly = `[Commands]
badjoke_enable = true
`
	cfgEelOnly = `[Commands]
eel_enable = true
`
	cfgTauntOnly = `[Commands]
taunt_enable = true
`
)

func TestRegister(t *testing.T) {
	bot := eelbot.New(newTestSession())
	config, err := ini.InsensitiveLoad([]byte(cfg))
	require.NoError(t, err)

	s := config.Section("Commands")
	require.ErrorContains(t, commands.Register(bot, s, nil), "command requires a database")
	require.NoError(t, commands.Register(bot, s, db))

	// Should be fault tolerant with invalid keys.
	s.Key("badjoke_delay").SetValue("foo")
	s.Key("channel_enable").SetValue("bar")
	require.NoError(t, commands.Register(bot, s, db))

	config, err = ini.InsensitiveLoad([]byte(cfgJokeOnly))
	require.NoError(t, err)
	s = config.Section("Commands")
	require.EqualError(t, commands.Register(bot, s, nil), "/badjoke command requires a database")
	require.NoError(t, commands.Register(bot, s, db))

	config, err = ini.InsensitiveLoad([]byte(cfgEelOnly))
	require.NoError(t, err)
	s = config.Section("Commands")
	require.EqualError(t, commands.Register(bot, s, nil), "/eel command requires a database")
	require.NoError(t, commands.Register(bot, s, db))

	config, err = ini.InsensitiveLoad([]byte(cfgTauntOnly))
	require.NoError(t, err)
	s = config.Section("Commands")
	require.EqualError(t, commands.Register(bot, s, nil), "/taunt command requires a database")
	require.NoError(t, commands.Register(bot, s, db))
}
