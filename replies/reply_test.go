package replies_test

import (
	"testing"

	"github.com/emseers/eelbot"
	"github.com/emseers/eelbot/replies"
	"github.com/stretchr/testify/require"
)

var (
	cfg = map[string]any{
		"caps": map[string]any{
			"enable":  true,
			"min_len": 5,
			"percent": 17,
			"timeout": 120,
		},
		"hello": map[string]any{
			"enable":  true,
			"percent": 33,
			"timeout": 600,
		},
		"goodbye": map[string]any{
			"enable":  true,
			"percent": 33,
			"timeout": 600,
		},
		"laugh": map[string]any{
			"enable":  true,
			"percent": 17,
			"timeout": 10,
		},
		"question": map[string]any{
			"enable":  true,
			"percent": 17,
			"timeout": 10,
		},
	}
)

func TestRegister(t *testing.T) {
	bot := eelbot.New(newTestSession())
	require.NoError(t, replies.Register(bot, cfg))

	// Should be fault tolerant with invalid keys.
	cfg["caps"].(map[string]any)["min_len"] = "foo"
	cfg["hello"].(map[string]any)["percent"] = "bar"
	require.NoError(t, replies.Register(bot, cfg))
}
