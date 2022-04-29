package commands_test

import (
	"testing"

	"github.com/emseers/eelbot"
	"github.com/emseers/eelbot/commands"
	"github.com/stretchr/testify/require"
)

var (
	cfg = map[string]any{
		"badjoke": map[string]any{"enable": true, "delay": float64(3)},
		"eel":     map[string]any{"enable": true},
		"taunt":   map[string]any{"enable": true},
		"channel": map[string]any{"enable": true},
		"flip":    map[string]any{"enable": true},
		"listen":  map[string]any{"enable": true},
		"play":    map[string]any{"enable": true},
		"ping":    map[string]any{"enable": true},
		"roll":    map[string]any{"enable": true},
		"say":     map[string]any{"enable": true},
		"saychan": map[string]any{"enable": true},
	}
	cfgJokeOnly = map[string]any{
		"badjoke": map[string]any{"enable": true},
	}
	cfgEelOnly = map[string]any{
		"eel": map[string]any{"enable": true},
	}
	cfgTauntOnly = map[string]any{
		"taunt": map[string]any{"enable": true},
	}
)

func TestRegister(t *testing.T) {
	bot := eelbot.New(newTestSession())

	require.ErrorContains(t, commands.Register(bot, cfg, nil), "command requires a database")
	require.NoError(t, commands.Register(bot, cfg, db))

	// Should be fault tolerant with invalid keys.
	cfg["badjoke"].(map[string]any)["delay"] = "foo"
	cfg["channel"].(map[string]any)["enable"] = "bar"
	require.NoError(t, commands.Register(bot, cfg, db))

	require.EqualError(t, commands.Register(bot, cfgJokeOnly, nil), "/badjoke command requires a database")
	require.NoError(t, commands.Register(bot, cfgJokeOnly, db))

	require.EqualError(t, commands.Register(bot, cfgEelOnly, nil), "/eel command requires a database")
	require.NoError(t, commands.Register(bot, cfgEelOnly, db))

	require.EqualError(t, commands.Register(bot, cfgTauntOnly, nil), "/taunt command requires a database")
	require.NoError(t, commands.Register(bot, cfgTauntOnly, db))
}
