// Package replies contains various eelbot replies.
package replies

import (
	"time"

	"github.com/emseers/eelbot"
	"gopkg.in/ini.v1"
)

type replyFromConfigFunc func(s *ini.Section, percent int) (*eelbot.Reply, error)

var replies = map[string]replyFromConfigFunc{}

// Register eelbot replies based on the given config.
func Register(bot *eelbot.Bot, s *ini.Section) error {
	for reply, f := range replies {
		enable, err := s.Key(reply + "_enable").Bool()
		if err == nil && enable {
			percent, err2 := s.Key(reply + "_percent").Int()
			if err2 != nil {
				percent = 100
			}
			r, err3 := f(s, percent)
			if err3 != nil {
				return err3
			}
			timeout, _ := s.Key(reply + "_timeout").Int64()
			r.Timeout = time.Second * time.Duration(timeout)
			bot.RegisterReply(*r)
		}
	}
	return nil
}
