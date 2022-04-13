package eelbot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Meta is the metadata passed in to a command or reply's Eval function.
type Meta struct {
	ChannelID string
	MessageID string
}

// A Bot is an instance of a discord bot that can listen for commands and do various things.
type Bot struct {
	dg      *discordgo.Session
	cmds    map[string]*Command
	replies []*Reply
}

// New creates a new Bot instance.
func New(token string) (bot *Bot, err error) {
	var dg *discordgo.Session
	if dg, err = discordgo.New("Bot " + token); err != nil {
		return
	}

	bot = &Bot{
		dg:   dg,
		cmds: map[string]*Command{},
	}

	bot.dg.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore all messages created by the bot itself.
		if m.Author.ID == s.State.User.ID {
			return
		}

		if strings.HasPrefix(m.Content, "/") {
			// strings.FieldsFunc is used over strings.Split to avoid empty values in the resulting slice.
			args := strings.FieldsFunc(m.Content, func(c rune) bool { return c == ' ' })
			cmd := strings.ToLower(args[0][1:])
			args = args[1:]
			if c, ok := bot.cmds[cmd]; ok {
				if err := evalCmd(cmd, c, s, m, args); err != nil {
					s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error: %s", err.Error()))
				}
			}
			return
		}

		for _, reply := range bot.replies {
			if reply.m.hasChannel(m.ChannelID) {
				continue
			}
			if reply.Eval(s, m) {
				if reply.Timeout > 0 {
					reply.m.addChannelWithTimedReset(m.ChannelID, reply.Timeout)
				}
				return
			}
		}
	})
	return
}

// Start starts the bot and opens a connection to Discord.
func (bot *Bot) Start() error {
	bot.createHelpCmd()
	return bot.dg.Open()
}

// Stop stops the bot and closes any open connections.
func (bot *Bot) Stop() error {
	return bot.dg.Close()
}
