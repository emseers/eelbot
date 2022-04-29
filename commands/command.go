// Package commands contains various eelbot commands.
package commands

import (
	"database/sql"

	"github.com/emseers/eelbot"
)

type commandFromConfigFunc func(map[string]any, *sql.DB) (*eelbot.Command, error)

var commands = map[string]commandFromConfigFunc{}

// Register eelbot commands based on the given opts. The expected format for the opts is a key for each command name
// with the value being a map[string]any for the options of the command. The value must have the "enable" key with the
// value being a bool, along with any other command specific options. An example config is as follows (JSONified):
//  {
//    "command_a": {
//      "enable": true,
//      "command_a_opt_1": "foo",
//      "command_a_opt_2": 2
//    },
//    "command_b": {
//      "enable": false
//    }
//  }
func Register(bot *eelbot.Bot, opts map[string]any, db *sql.DB) error {
	for cmd, f := range commands {
		if cmdOpts, ok := opts[cmd].(map[string]any); ok {
			if enable, _ := cmdOpts["enable"].(bool); enable {
				c, err := f(cmdOpts, db)
				if err != nil {
					return err
				}
				bot.RegisterCommand(cmd, *c)
			}
		}
	}
	return nil
}
