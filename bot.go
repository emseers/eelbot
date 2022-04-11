package eelbot

import (
	"fmt"
	"io"
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

		meta := &Meta{
			ChannelID: m.ChannelID,
			MessageID: m.ID,
		}

		if strings.HasPrefix(m.Content, "/") {
			// strings.FieldsFunc is used over strings.Split to avoid empty values in the resulting slice.
			args := strings.FieldsFunc(m.Content, func(c rune) bool { return c == ' ' })
			cmd := strings.ToLower(args[0][1:])
			args = args[1:]
			if c, ok := bot.cmds[cmd]; ok {
				if err := bot.evalCmd(cmd, c, meta, args); err != nil {
					bot.SendMsg(m.ChannelID, fmt.Sprintf("Error: %s", err.Error()))
				}
			}
			return
		}

		for _, reply := range bot.replies {
			if reply.m.hasChannel(meta.ChannelID) {
				continue
			}
			if reply.Eval(bot, meta, m.Content) {
				if reply.Timeout > 0 {
					reply.m.addChannelWithTimedReset(meta.ChannelID, reply.Timeout)
				}
				return
			}
		}
	})
	return
}

// UpdateGameStatus updates the bot's game status.
func (bot *Bot) UpdateGameStatus(game string) {
	_ = bot.dg.UpdateGameStatus(0, game)
}

// UpdateListeningStatus updates the bot's listening status.
func (bot *Bot) UpdateListeningStatus(song string) {
	_ = bot.dg.UpdateListeningStatus(song)
}

// SendMsg sends the given message to the given channel.
func (bot *Bot) SendMsg(channelID, msg string) {
	_, _ = bot.dg.ChannelMessageSend(channelID, msg)
}

// SendFile sends the given file to the given channel.
func (bot *Bot) SendFile(channelID, filename string, file io.Reader) {
	_, _ = bot.dg.ChannelFileSend(channelID, filename, file)
}

// DeleteMsg deletes the given message ID from the given channel.
func (bot *Bot) DeleteMsg(channelID, msgID string) {
	_ = bot.dg.ChannelMessageDelete(channelID, msgID)
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
