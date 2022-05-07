package commands

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/emseers/eelbot"
)

func init() {
	commands["roll"] = rollFromConfig
}

func rollFromConfig(map[string]any, *sql.DB, time.Duration) (*eelbot.Command, error) {
	return RollCommand(), nil
}

// RollCommand returns an *eelbot.Command that replies with the result of a dice roll.
func RollCommand() *eelbot.Command {
	return &eelbot.Command{
		MinArgs: 1,
		MaxArgs: 2,
		Summary: "Rolls a die.",
		Eval: func(s eelbot.Session, m *discordgo.MessageCreate, args []string) error {
			var (
				lower int64
				upper int64
				err   error
			)
			switch len(args) {
			case 1:
				if upper, err = strconv.ParseInt(args[0], 10, 64); err != nil {
					return err
				}
			case 2:
				if lower, err = strconv.ParseInt(args[0], 10, 64); err != nil {
					return err
				}
				if upper, err = strconv.ParseInt(args[1], 10, 64); err != nil {
					return err
				}
			}
			if upper < lower {
				lower, upper = upper, lower
			}
			upper++
			s.ChannelMessageSend(m.ChannelID, fmt.Sprint(rand.Int63n(upper-lower)+lower))
			return nil
		},
	}
}
