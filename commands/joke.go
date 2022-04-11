package commands

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/emseers/eelbot"
)

// JokeCommand returns an *eelbot.Command that reads and replies with a joke from the given db. Two line jokes use the
// given delay between replies.
func JokeCommand(db *sql.DB, delay time.Duration) *eelbot.Command {
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
		Eval: func(bot *eelbot.Bot, meta *eelbot.Meta, args []string) error {
			var (
				query string
				line1 string
				line2 sql.NullString
			)
			if args[0] == "me" {
				query = "SELECT text, punchline FROM jokes ORDER BY RANDOM() LIMIT 1;"
			} else if num, err := strconv.ParseUint(args[0], 10, 0); err == nil {
				query = fmt.Sprintf("SELECT text, punchline FROM jokes WHERE id=%d;", num)
			} else {
				return unknownDirectiveErr(args[0])
			}
			row := db.QueryRow(query)
			if err := row.Scan(&line1, &line2); err != nil {
				return err
			}
			bot.SendMsg(meta.ChannelID, line1)
			if line2.Valid {
				time.Sleep(delay)
				bot.SendMsg(meta.ChannelID, line2.String)
			}
			return nil
		},
	}
}
