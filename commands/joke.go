package commands

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/emseers/eelbot"
	"gopkg.in/ini.v1"
)

func init() {
	commands["badjoke"] = badjokeFromConfig
}

func badjokeFromConfig(s *ini.Section, db *sql.DB) (*eelbot.Command, error) {
	if db == nil {
		return nil, requiresDatabaseErr("badjoke")
	}
	delay, err := s.Key("badjoke_delay").Int64()
	if err != nil {
		delay = 3
	}
	return JokeCommand(db, time.Second*time.Duration(delay)), nil
}

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
		Eval: func(s eelbot.Session, m *discordgo.MessageCreate, args []string) error {
			var (
				query string
				line1 string
				line2 sql.NullString
			)
			if args[0] == "me" {
				query = "SELECT text, punchline FROM jokes ORDER BY RANDOM() LIMIT 1;"
			} else if num, err := strconv.ParseUint(args[0], 10, 64); err == nil {
				query = fmt.Sprintf("SELECT text, punchline FROM jokes WHERE id=%d;", num)
			} else {
				return unknownDirectiveErr(args[0])
			}
			row := db.QueryRow(query)
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
