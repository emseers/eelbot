// Package replies contains various eelbot replies.
package replies

import (
	"math/rand"
	"time"

	"github.com/emseers/eelbot"
)

type replyFromConfigFunc func(opts map[string]any, percent int, minDelay, maxDelay time.Duration) (*eelbot.Reply, error)

var (
	replies = map[string]replyFromConfigFunc{}

	// Rand is the random number source that's used by some commands. It is public for the purposes of testing.
	Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// Register eelbot replies based on the given config. The expected format for the opts is a key for each reply type
// with the value being a map[string]any for the options of the reply. The value must have the "enable" key with the
// value being a bool, along with any other command specific options. The value can also optionally have keys for
// "percent" (that controls the percent of the time the reply triggers on a valid match), "min_delay" (the minimum time
// to wait before replying, in seconds), "max_delay" (the maximum time to wait before replying, in seconds) and
// "timeout" (time in seconds for which this reply won't trigger again after triggering successfully). An example config
// is as follows (JSONified):
//
//	{
//	  "reply_a": {
//	    "enable": true,
//	    "percent": 72,
//	    "min_delay": 3,
//	    "max_delay": 6,
//	    "timeout": 120,
//	    "reply_a_opt_1": "foo",
//	    "reply_a_opt_2": 2
//	  },
//	  "reply_b": {
//	    "enable": false
//	  }
//	}
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
			time.Sleep(time.Duration(Rand.Int63n(int64(maxDelay - minDelay))))
		}
		s.ChannelMessageSend(channelID, msg)
	}()
}
