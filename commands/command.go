// Package commands contains various eelbot commands.
package commands

import (
	"context"
	"database/sql"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/emseers/eelbot"
)

type commandFromConfigFunc func(map[string]any, *sql.DB, time.Duration) (*eelbot.Command, error)

var (
	commands = map[string]commandFromConfigFunc{}

	// Rand is the random number source that's used by some commands. It is public for the purposes of testing.
	Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// Register eelbot commands based on the given opts. The expected format for the opts is a key for each command name
// with the value being a map[string]any for the options of the command. The value must have the "enable" key with the
// value being a bool, along with any other command specific options. An example config is as follows (JSONified):
//
//	{
//	  "command_a": {
//	    "enable": true,
//	    "command_a_opt_1": "foo",
//	    "command_a_opt_2": 2
//	  },
//	  "command_b": {
//	    "enable": false
//	  }
//	}
func Register(bot *eelbot.Bot, opts map[string]any, db *sql.DB, dbTimeout time.Duration) error {
	for cmd, f := range commands {
		if cmdOpts, ok := opts[cmd].(map[string]any); ok {
			if enable, _ := cmdOpts["enable"].(bool); enable {
				c, err := f(cmdOpts, db, dbTimeout)
				if err != nil {
					return err
				}
				bot.RegisterCommand(cmd, *c)
			}
		}
	}
	return nil
}

func getDuration(v any, def time.Duration) (dur time.Duration, err error) {
	switch d := v.(type) {
	case int:
		dur = time.Second * time.Duration(d)
	case int8:
		dur = time.Second * time.Duration(d)
	case int16:
		dur = time.Second * time.Duration(d)
	case int32:
		dur = time.Second * time.Duration(d)
	case int64:
		dur = time.Second * time.Duration(d)
	case uint:
		dur = time.Second * time.Duration(d)
	case uint8:
		dur = time.Second * time.Duration(d)
	case uint16:
		dur = time.Second * time.Duration(d)
	case uint32:
		dur = time.Second * time.Duration(d)
	case uint64:
		dur = time.Second * time.Duration(d)
	case float32:
		dur = time.Second * time.Duration(d)
	case float64:
		dur = time.Second * time.Duration(d)
	case string:
		dur, err = time.ParseDuration(d)
	case time.Duration:
		dur = d
	default:
		dur = def
	}
	return
}

func queryRow(db *sql.DB, dbTimeout time.Duration, query string, args ...any) (*sql.Row, context.CancelFunc) {
	ctx, cancel := context.Background(), func() {}
	if dbTimeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, dbTimeout)
	}
	return db.QueryRowContext(ctx, query, args...), cancel
}

// Returns a query to select a random row from a table that contains an integer primary key column called 'id' that
// maintains gapless sequential values. This is much more performant than "ORDER BY RANDOM()" since it doesn't require
// ordering all rows.
func randRowQuery(table string, cols []string) string {
	return fmt.Sprintf(
		"SELECT %s FROM %[2]s WHERE id=(SELECT (MIN(id) + trunc(random()*(MAX(id)-MIN(id)))::integer) FROM %[2]s);",
		strings.Join(cols, ", "),
		table,
	)
}
