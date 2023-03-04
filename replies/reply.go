// Package replies contains various eelbot replies.
package replies

import (
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"regexp"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/emseers/eelbot"
)

var (
	// Rand is the random number source that's used for replies. It is public for the purposes of testing.
	Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// Register eelbot replies based on the given config. The expected format for the opts is an emtry for reply type with
// with the value being a map[string]any for the options of the reply. The value must have a "tableName" key with the
// value being a string (and a table with that name must exist in the "reply" schema of the database) and a "regexps"
// key with the value being a []any (but individual values must be valid golang flavor regular expression strings) that
// specifies the condition(s) to trigger the reply. The value can also optionally have keys for "percent" (that controls
// the percent of the time the reply triggers on a valid match), "min_delay" (the minimum time to wait before replying,
// in seconds), "max_delay" (the maximum time to wait before replying, in seconds) and "timeout" (time in seconds for
// which this reply won't trigger again after triggering successfully). An example config is as follows (JSONified):
//
//	[
//	  {
//	    "tableName": "reply_a",
//	    "regexps": ["a{3,6}", "b{4,7}"],
//	    "percent": 72,
//	    "min_delay": 3,
//	    "max_delay": 6,
//	    "timeout": 120
//	  },
//	  {
//	    "tableName": "reply_b",
//	    "regexps": ["foobar"]
//	  }
//	]
func Register(bot *eelbot.Bot, replies []any, db *sql.DB, dbTimeout time.Duration) error {
	for _, reply := range replies {
		if replyOpts, ok := reply.(map[string]any); ok {
			if db == nil {
				return fmt.Errorf("reply requires a database")
			}
			tableName, _ := replyOpts["table_name"].(string)
			if tableName == "" {
				return fmt.Errorf("missing table_name for reply")
			}
			// Fetch a reply from the table to ensure it exists and has the right schema.
			if _, err := fetchReply(db, dbTimeout, tableName); err != nil && !errors.Is(err, sql.ErrNoRows) {
				return err
			}
			rexps := []*regexp.Regexp{}
			regexStrs, _ := replyOpts["regexps"].([]any)
			for _, r := range regexStrs {
				if regexStr, _ := r.(string); len(regexStr) != 0 {
					regex, err := regexp.Compile(regexStr)
					if err != nil {
						return err
					}
					rexps = append(rexps, regex)
				}
			}
			if len(rexps) == 0 {
				return fmt.Errorf("one or more valid regexps must be provided for reply %s", tableName)
			}
			percent := getInt(replyOpts["percent"], 100)
			var minDelay, maxDelay, timeout time.Duration
			var err error
			if minDelay, err = getDuration(replyOpts["min_delay"], 0); err != nil {
				return err
			}
			if maxDelay, err = getDuration(replyOpts["max_delay"], 0); err != nil {
				return err
			}
			if timeout, err = getDuration(replyOpts["timeout"], 0); err != nil {
				return err
			}

			bot.RegisterReply(eelbot.Reply{
				Eval:    EvalFn(db, dbTimeout, tableName, rexps, percent, minDelay, maxDelay),
				Timeout: timeout,
			})
		}
	}
	return nil
}

// EvalFn returns a function suitable for use as a *eelbot.Reply's Eval function.
func EvalFn(db *sql.DB, dbTimeout time.Duration, tableName string, rexps []*regexp.Regexp, percent int,
	minDelay, maxDelay time.Duration) func(s eelbot.Session, m *discordgo.MessageCreate) bool {
	return func(s eelbot.Session, m *discordgo.MessageCreate) bool {
		if match(m.Content, rexps...) && roll(percent) {
			reply, err := fetchReply(db, dbTimeout, tableName)
			if err != nil {
				return false
			}
			asyncReply(s, m.ChannelID, reply, minDelay, maxDelay)
			return true
		}
		return false
	}
}
