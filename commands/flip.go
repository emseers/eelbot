package commands

import (
	"math/rand"

	"github.com/emseers/eelbot"
)

// FlipCommand returns an *eelbot.Command that replies with the result of a coin toss.
func FlipCommand() *eelbot.Command {
	return &eelbot.Command{
		MinArgs: 0,
		MaxArgs: 0,
		Summary: "Flips a coin.",
		Eval: func(bot *eelbot.Bot, meta *eelbot.Meta, _ []string) error {
			var result string
			if rand.Intn(6000) == 0 {
				result = "Landed on edge"
			} else if rand.Intn(2) == 0 {
				result = "Heads"
			} else {
				result = "Tails"
			}
			bot.SendMsg(meta.ChannelID, result)
			return nil
		},
	}
}
