// Package helpers implements helpers to startup services, mainly for testing purposes.
package helpers

import (
	"context"
	"database/sql"
	"time"
)

// A CloseFunc tells an instance to close/shutdown. It should be called after the instance is no longer needed.
type CloseFunc func()

func pingWithTimeout(db *sql.DB, timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return db.PingContext(ctx)
}
