// Package replies contains various eelbot replies.
package replies

import (
	"math/rand"
	"time"

	"github.com/emseers/eelbot"
)

type replyFromConfigFunc func(opts map[string]any, percent int, minDelay, maxDelay time.Duration) (*eelbot.Reply, error)

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
				minDelay, _ := replyOpts["min_delay"].(float64)
				maxDelay, _ := replyOpts["min_delay"].(float64)
				r, err := f(
					replyOpts,
					int(percent),
					time.Second*time.Duration(minDelay),
					time.Second*time.Duration(maxDelay),
				)
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

func asyncReply(s eelbot.Session, channelID, msg string, minDelay, maxDelay time.Duration) {
	if maxDelay < minDelay {
		maxDelay = minDelay
	}
	go func() {
		s.ChannelTyping(channelID)
		time.Sleep(minDelay)
		if maxDelay > minDelay {
			time.Sleep(time.Duration(rand.Int63n(int64(maxDelay - minDelay))))
		}
		s.ChannelMessageSend(channelID, msg)
	}()
}
