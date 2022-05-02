package commands

import (
	"bytes"
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/emseers/eelbot"
)

func init() {
	commands["taunt"] = tauntFromConfig
}

func tauntFromConfig(_ map[string]any, db *sql.DB, dbTimeout time.Duration) (*eelbot.Command, error) {
	if db == nil {
		return nil, requiresDatabaseErr("taunt")
	}
	_, err := exec(db, dbTimeout, `CREATE TABLE IF NOT EXISTS taunts (
  id   integer PRIMARY KEY,
  name text NOT NULL,
  file bytea NOT NULL
);`)
	if err != nil {
		return nil, err
	}
	return TauntCommand(db, dbTimeout), nil
}

// TauntCommand returns an *eelbot.Command that reads and replies with a taunt from the given db.
func TauntCommand(db *sql.DB, dbTimeout time.Duration) *eelbot.Command {
	return &eelbot.Command{
		MinArgs: 1,
		MaxArgs: 1,
		Summary: "Posts a taunt.",
		Usage: `/%[1]s NUM

Posts a taunt from the database. NUM can either be a valid taunt number from the database, or "me" for a random taunt.

Examples:
  /%[1]s me
  /%[1]s 42
`,
		Eval: func(s eelbot.Session, m *discordgo.MessageCreate, args []string) error {
			var (
				row    *sql.Row
				cancel context.CancelFunc
				name   string
				file   []byte
			)
			if args[0] == "me" {
				row, cancel = queryRow(db, dbTimeout, randRowQuery("taunts", []string{"name", "file"}))
			} else if num, err := strconv.ParseUint(args[0], 10, 64); err == nil {
				row, cancel = queryRow(db, dbTimeout, "SELECT name, file FROM taunts WHERE id=$1;", num)
			} else {
				return unknownDirectiveErr(args[0])
			}
			defer cancel()
			if err := row.Scan(&name, &file); err != nil {
				return err
			}
			s.ChannelFileSend(m.ChannelID, name, bytes.NewReader(file))
			return nil
		},
	}
}
