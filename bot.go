package eelbot

import (
	"io"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/emseers/eelbot/msg"
)

// Options is used to specify options needed to create an instance of a Bot.
type Options struct {
	Token              string
	MsgTimeout         time.Duration
	MultiLineJokeDelay time.Duration
	DBName             string
}

// A Bot is an instance of a discord bot that can listen for commands and do various things.
type Bot struct {
	dg             *discordgo.Session
	msgInterpreter *msg.Interpreter
}

// New creates a new Bot instance.
func New(opts *Options) (bot *Bot, err error) {
	var dg *discordgo.Session
	if dg, err = discordgo.New("Bot " + opts.Token); err != nil {
		return
	}

	var interpreter *msg.Interpreter
	if interpreter, err = msg.NewInterpreter(opts.DBName); err != nil {
		return
	}
	interpreter.MsgTimeout = opts.MsgTimeout
	interpreter.MultiLineJokeDelay = opts.MultiLineJokeDelay

	bot = &Bot{
		dg:             dg,
		msgInterpreter: interpreter,
	}

	bot.dg.AddHandler(bot.guildCreateHandler)
	bot.dg.AddHandler(bot.messageCreateHandler)
	return
}

// Start starts the bot and open a connection to Discord.
func (bot *Bot) Start() (err error) {
	err = bot.dg.Open()
	return
}

// UpdateStatus updates the bot's status.
func (bot *Bot) UpdateStatus(game string) {
	_ = bot.dg.UpdateStatus(0, game)
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

// Stop stops the bot and closes any open connections.
func (bot *Bot) Stop() (err error) {
	err = bot.msgInterpreter.Stop()
	if err != nil {
		return
	}

	err = bot.dg.Close()
	return
}
