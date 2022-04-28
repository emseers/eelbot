// Package replies contains various eelbot replies.
package replies

import (
	"time"

	"github.com/emseers/eelbot"
)

type replyFromConfigFunc func(opts map[string]any, percent int) (*eelbot.Reply, error)

var replies = map[string]replyFromConfigFunc{}

// Register eelbot replies based on the given config.
func Register(bot *eelbot.Bot, opts map[string]any) error {
	for reply, f := range replies {
		if replyOpts, ok := opts[reply].(map[string]any); ok {
			if enable, _ := replyOpts["enable"].(bool); enable {
				percent, ok2 := replyOpts["percent"].(float64)
				if !ok2 {
					percent = 100
				}
				r, err := f(replyOpts, int(percent))
				if err != nil {
					return err
				}
				timeout, _ := replyOpts["timeout"].(float64)
				r.Timeout = time.Second * time.Duration(timeout)
				bot.RegisterReply(*r)
			}
		}
	}
	return nil
}
