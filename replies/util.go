package replies

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"time"

	"github.com/emseers/eelbot"
)

func getInt(v any, def int) int {
	switch i := v.(type) {
	case int:
		return i
	case int8:
		return int(i)
	case int16:
		return int(i)
	case int32:
		return int(i)
	case int64:
		return int(i)
	case uint:
		return int(i)
	case uint8:
		return int(i)
	case uint16:
		return int(i)
	case uint32:
		return int(i)
	case uint64:
		return int(i)
	case float32:
		return int(i)
	case float64:
		return int(i)
	default:
		return def
	}
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

func roll(percent int) bool {
	return Rand.Intn(100) < percent
}

func match(s string, rexps ...*regexp.Regexp) bool {
	for _, rexp := range rexps {
		if rexp.MatchString(s) {
			return true
		}
	}
	return false
}

// Returns a random reply from a reply table that contains an integer primary key column called 'id' that maintains
// gapless sequential values. This is much more performant than "ORDER BY RANDOM()" since it doesn't require ordering
// all rows.
func fetchReply(db *sql.DB, dbTimeout time.Duration, table string) (text string, err error) {
	ctx := context.Background()
	if dbTimeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, dbTimeout)
		defer cancel()
	}
	row := db.QueryRowContext(ctx, fmt.Sprintf(
		`SELECT text FROM reply.%[1]s
		 WHERE id=(SELECT (MIN(id) + trunc(random()*(MAX(id)-MIN(id)))::integer) FROM reply.%[1]s);`,
		table))
	err = row.Scan(&text)
	return
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
