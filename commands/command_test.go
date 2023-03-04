package commands_test

import (
	"testing"
	"time"

	"github.com/emseers/eelbot"
	"github.com/emseers/eelbot/commands"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

var (
	cfg = map[string]any{
		"badjoke": map[string]any{"enable": true, "delay": 3},
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
	require.ErrorContains(t, commands.Register(bot, cfg, nil, 0), "command requires a database")
	require.NoError(t, commands.Register(bot, cfg, db, time.Second))

	require.EqualError(t, commands.Register(bot, cfgJokeOnly, nil, 0), "/badjoke command requires a database")
	require.NoError(t, commands.Register(bot, cfgJokeOnly, db, time.Second))

	require.EqualError(t, commands.Register(bot, cfgEelOnly, nil, 0), "/eel command requires a database")
	require.NoError(t, commands.Register(bot, cfgEelOnly, db, time.Second))

	require.EqualError(t, commands.Register(bot, cfgTauntOnly, nil, 0), "/taunt command requires a database")
	require.NoError(t, commands.Register(bot, cfgTauntOnly, db, time.Second))

	cfgJokeOnly["badjoke"].(map[string]any)["delay"] = "foo"
	require.ErrorContains(t, commands.Register(bot, cfgJokeOnly, db, time.Second), "invalid duration")
	cfgTauntOnly["taunt"].(map[string]any)["enable"] = "bar" // Should treat non bool values as false.
	require.NoError(t, commands.Register(bot, cfgTauntOnly, db, time.Second))
}
