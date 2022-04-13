// Package commands contains various eelbot commands.
package commands

import (
	"database/sql"

	"github.com/emseers/eelbot"
	"gopkg.in/ini.v1"
)

type commandFromConfigFunc func(*ini.Section, *sql.DB) (*eelbot.Command, error)

var commands = map[string]commandFromConfigFunc{}

// Register eelbot commands based on the given config.
func Register(bot *eelbot.Bot, s *ini.Section, db *sql.DB) error {
	for cmd, f := range commands {
		enable, err := s.Key(cmd + "_enable").Bool()
		if err == nil && enable {
			c, err2 := f(s, db)
			if err2 != nil {
				return err2
			}
			bot.RegisterCommand(cmd, *c)
		}
	}
	return nil
}
