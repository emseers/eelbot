// Package eelbot contains a simple bot that can listen to commands and do things.
package eelbot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// A Bot is an instance of a discord bot that can listen for commands and do various things.
type Bot struct {
	sess    Session
	cmds    map[string]*Command
	replies []*Reply
}

// New creates a new Bot instance.
func New(sess Session) *Bot {
	bot := &Bot{
		sess: sess,
		cmds: map[string]*Command{},
	}

	bot.sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		// Ignore all messages created by the bot itself.
		if m.Author.ID == s.State.User.ID {
			return
		}

		defer func() {
			if r := recover(); r != nil {
				bot.sess.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error: %v", r))
			}
		}()

		if strings.HasPrefix(m.Content, "/") {
			// strings.FieldsFunc is used over strings.Split to avoid empty values in the resulting slice.
			args := strings.FieldsFunc(m.Content, func(c rune) bool { return c == ' ' })
			cmd := strings.ToLower(args[0][1:])
			args = args[1:]
			if c, ok := bot.cmds[cmd]; ok {
				if err := evalCmd(cmd, c, bot.sess, m, args); err != nil {
					bot.sess.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Error: %s", err.Error()))
				}
			}
			return
		}

		for _, reply := range bot.replies {
			if reply.m.hasChannel(m.ChannelID) {
				continue
			}
			if reply.Eval(bot.sess, m) {
				if reply.Timeout > 0 {
					reply.m.addChannel(m.ChannelID, reply.Timeout)
				}
				return
			}
		}
	})

	return bot
}

// Start starts the bot and opens a connection to Discord.
func (bot *Bot) Start() error {
	bot.createHelpCmd()
	return bot.sess.Open()
}

// Stop stops the bot and closes any open connections.
func (bot *Bot) Stop() error {
	return bot.sess.Close()
}
