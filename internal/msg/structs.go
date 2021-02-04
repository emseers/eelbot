package msg

import (
	"database/sql"
	"io"
	"sync"
	"time"
)

// UpdateStatusFunc is a callback function that will be used by the Interpreter to update the bot's status.
type UpdateStatusFunc func(game string)

// SendMsgFunc is a callback function that will be used by the Interpreter to send messages.
type SendMsgFunc func(channelID, msg string)

// SendFileFunc is a callback function that will be used by the Interpreter to send files.
type SendFileFunc func(channelID, filename string, file io.Reader)

// DeleteMsgFunc is a callback function that will be used by the Interpreter to delete messages.
type DeleteMsgFunc func(channelID, msgID string)

// CallBackCtx is the context to hold all callback functions for the message interpreter.
type CallBackCtx struct {
	UpdateStatus UpdateStatusFunc
	SendMsg      SendMsgFunc
	SendFile     SendFileFunc
	DeleteMsg    DeleteMsgFunc
}

// NewInterpreterCtx is the context needed to create a message interpreter.
type NewInterpreterCtx struct {
	MsgTimeout         time.Duration
	MultiLineJokeDelay time.Duration
	SQLiteDB           string
	TauntsFolder       string
}

// An Interpreter can parse and respond to certain messages.
type Interpreter struct {
	msgTimeout         time.Duration
	multiLineJokeDelay time.Duration
	sqliteDB           *sql.DB
	tauntsFolder       string
	flagAllCaps        *flagMap
	flagQuestion       *flagMap
	flagHello          *flagMap
	flagGoodbye        *flagMap
}

// A flagMap is a wrapper around a sync.Map to provide type safety to act like a map[string]struct{}.
type flagMap struct {
	syncMap *sync.Map
}
