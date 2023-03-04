package replies_test

import (
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/emseers/eelbot"
	"github.com/emseers/eelbot/replies"
	"github.com/stretchr/testify/require"
)

var (
	cfg = []any{
		map[string]any{},
		map[string]any{
			"table_name": "nonexistent",
		},
		map[string]any{
			"table_name": incorrectFormatTable,
		},
		map[string]any{
			"table_name": sampleTable,
		},
		map[string]any{
			"table_name": sampleTable,
			"regexps":    []any{`/\`},
		},
		map[string]any{
			"table_name": sampleTable,
			"regexps":    []any{"hello"},
			"min_delay":  "foo",
		},
		map[string]any{
			"table_name": sampleTable,
			"regexps":    []any{"hello"},
			"max_delay":  "bar",
		},
		map[string]any{
			"table_name": sampleTable,
			"regexps":    []any{"hello"},
			"timeout":    "baz",
		},
		map[string]any{
			"table_name": sampleTable,
			"regexps":    []any{"hello"},
			"percent":    100,
			"min_delay":  "50ms",
			"max_delay":  "150ms",
			"timeout":    10,
		},
	}
)

func TestRegister(t *testing.T) {
	bot := eelbot.New(newTestSession())
	require.EqualError(t, replies.Register(bot, cfg, nil, 0), "reply requires a database")

	// Missing table name.
	require.EqualError(t, replies.Register(bot, cfg, db, time.Second), "missing table_name for reply")

	// Nonexistent table.
	cfg := cfg[1:]
	require.ErrorContains(t, replies.Register(bot, cfg, db, time.Second),
		"relation \"reply.nonexistent\" does not exist")

	// Table with incorrect format.
	cfg = cfg[1:]
	require.ErrorContains(t, replies.Register(bot, cfg, db, time.Second), "column \"text\" does not exist")

	// Missing regexps.
	cfg = cfg[1:]
	require.ErrorContains(t, replies.Register(bot, cfg, db, time.Second),
		"one or more valid regexps must be provided for reply")

	// Invalid regex.
	cfg = cfg[1:]
	require.ErrorContains(t, replies.Register(bot, cfg, db, time.Second), "error parsing regexp")

	// Invalid min delay.
	cfg = cfg[1:]
	require.ErrorContains(t, replies.Register(bot, cfg, db, time.Second), "invalid duration")

	// Invalid max delay.
	cfg = cfg[1:]
	require.ErrorContains(t, replies.Register(bot, cfg, db, time.Second), "invalid duration")

	// Invalid timeout.
	cfg = cfg[1:]
	require.ErrorContains(t, replies.Register(bot, cfg, db, time.Second), "invalid duration")

	// Valid config.
	cfg = cfg[1:]
	require.NoError(t, replies.Register(bot, cfg, db, time.Second))
}

func TestReply(t *testing.T) {
	s := newTestSession()
	f := replies.EvalFn(db, time.Second, sampleTable, []*regexp.Regexp{regexp.MustCompile("hello")}, 100,
		50*time.Millisecond, 150*time.Millisecond)

	require.False(t, f(s, newMsgCreate("goodbye", testChannelID)))
	require.True(t, f(s, newMsgCreate("hello", testChannelID)))
	time.Sleep(10 * time.Millisecond) // Not enough delay for goroutine to finish the write.
	require.Nil(t, s.messages[testChannelID])
	time.Sleep(200 * time.Millisecond) // Enough delay for goroutine to finish the write.
	require.Equal(t, "hi", strings.TrimSpace(s.messages[testChannelID].String()))

	s.messages[testChannelID].Reset()
	f = replies.EvalFn(db, time.Second, sampleTable, []*regexp.Regexp{regexp.MustCompile("hello")}, 100,
		50*time.Millisecond, 0) // maxDelay is less than minDelay.
	require.True(t, f(s, newMsgCreate("hello", testChannelID)))
	time.Sleep(10 * time.Millisecond) // Not enough delay for goroutine to finish the write.
	require.Equal(t, "", strings.TrimSpace(s.messages[testChannelID].String()))
	time.Sleep(100 * time.Millisecond) // Enough delay for goroutine to finish the write.
	require.Equal(t, "hi", strings.TrimSpace(s.messages[testChannelID].String()))
}
