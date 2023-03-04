package commands

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/emseers/eelbot"
)

func init() {
	commands["badjoke"] = badjokeFromConfig
}

func badjokeFromConfig(opts map[string]any, db *sql.DB, dbTimeout time.Duration) (*eelbot.Command, error) {
	if db == nil {
		return nil, requiresDatabaseErr("badjoke")
	}
	delay, err := getDuration(opts["delay"], 3*time.Second)
	if err != nil {
		return nil, err
	}
	return JokeCommand(db, dbTimeout, delay), nil
}

// JokeCommand returns an *eelbot.Command that reads and replies with a joke from the given db. Two line jokes use the
// given delay between replies.
func JokeCommand(db *sql.DB, dbTimeout, delay time.Duration) *eelbot.Command {
	return &eelbot.Command{
		MinArgs: 1,
		MaxArgs: 1,
		Summary: "Posts a joke.",
		Usage: `/%[1]s NUM

Posts a joke from the database. NUM can either be a valid joke number from the database, or "me" for a random joke.

Examples:
  /%[1]s me
  /%[1]s 42
`,
		Eval: func(s eelbot.Session, m *discordgo.MessageCreate, args []string) error {
			var (
				row    *sql.Row
				cancel context.CancelFunc
				line1  string
				line2  sql.NullString
			)
			if args[0] == "me" {
				row, cancel = queryRow(db, dbTimeout, randRowQuery("jokes", []string{"text", "punchline"}))
			} else if num, err := strconv.ParseUint(args[0], 10, 64); err == nil {
				row, cancel = queryRow(db, dbTimeout, "SELECT text, punchline FROM jokes WHERE id=$1;", num)
			} else {
				return unknownDirectiveErr(args[0])
			}
			defer cancel()
			if err := row.Scan(&line1, &line2); err != nil {
				return err
			}
			s.ChannelMessageSend(m.ChannelID, line1)
			if line2.Valid {
				s.ChannelTyping(m.ChannelID)
				time.Sleep(delay)
				s.ChannelMessageSend(m.ChannelID, line2.String)
			}
			return nil
		},
	}
}
