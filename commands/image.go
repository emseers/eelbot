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
	commands["eel"] = imageFromConfig
}

func imageFromConfig(_ map[string]any, db *sql.DB, dbTimeout time.Duration) (*eelbot.Command, error) {
	if db == nil {
		return nil, requiresDatabaseErr("eel")
	}
	return ImageCommand(db, dbTimeout), nil
}

// ImageCommand returns an *eelbot.Command that reads and replies with an image from the given db.
func ImageCommand(db *sql.DB, dbTimeout time.Duration) *eelbot.Command {
	return &eelbot.Command{
		MinArgs: 1,
		MaxArgs: 1,
		Summary: "Posts an image.",
		Usage: `/%[1]s NUM

Posts an image from the database. NUM can either be a valid image number from the database, or "me" for a random image.

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
				row, cancel = queryRow(db, dbTimeout, randRowQuery("images", []string{"name", "file"}))
			} else if num, err := strconv.ParseUint(args[0], 10, 64); err == nil {
				row, cancel = queryRow(db, dbTimeout, "SELECT name, file FROM images WHERE id=$1;", num)
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
