package commands

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/bwmarrin/discordgo"
	"github.com/emseers/eelbot"
	"gopkg.in/ini.v1"
)

func init() {
	commands["eel"] = imageFromConfig
}

func imageFromConfig(_ *ini.Section, db *sql.DB) (*eelbot.Command, error) {
	if db == nil {
		return nil, requiresDatabaseErr("eel")
	}
	return ImageCommand(db), nil
}

// ImageCommand returns an *eelbot.Command that reads and replies with an image from the given db.
func ImageCommand(db *sql.DB) *eelbot.Command {
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
		Eval: func(s *discordgo.Session, m *discordgo.MessageCreate, args []string) error {
			var (
				query string
				path  string
			)
			if args[0] == "me" {
				query = "SELECT path FROM images ORDER BY RANDOM() LIMIT 1;"
			} else if num, err := strconv.ParseUint(args[0], 10, 64); err == nil {
				query = fmt.Sprintf("SELECT path FROM images WHERE id=%d;", num)
			} else {
				return unknownDirectiveErr(args[0])
			}
			row := db.QueryRow(query)
			if err := row.Scan(&path); err != nil {
				return err
			}
			name := filepath.Base(path)
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			s.ChannelFileSend(m.ChannelID, name, file)
			return nil
		},
	}
}
