package replies_test

import (
	"testing"

	"github.com/emseers/eelbot"
	"github.com/emseers/eelbot/replies"
	"github.com/stretchr/testify/require"
	"gopkg.in/ini.v1"
)

const (
	cfg = `[Replies]
caps_enable = true
caps_min_len = 5
caps_percent = 17
caps_timeout = 120

hello_enable = true
hello_percent = 33
hello_timeout = 600

goodbye_enable = true
goodbye_percent = 33
goodbye_timeout = 600

laugh_enable = true
laugh_percent = 17
laugh_timeout = 10

question_enable = true
question_percent = 17
question_timeout = 10
`
)

func TestRegister(t *testing.T) {
	bot := eelbot.New(newTestSession())
	config, err := ini.InsensitiveLoad([]byte(cfg))
	require.NoError(t, err)

	s := config.Section("Replies")
	require.NoError(t, replies.Register(bot, s))

	// Should be fault tolerant with invalid keys.
	s.Key("caps_min_len").SetValue("foo")
	s.Key("hello_percent").SetValue("bar")
	require.NoError(t, replies.Register(bot, s))
}
