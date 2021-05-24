package eelbot

import (
	"io"

	"github.com/bwmarrin/discordgo"
	"github.com/emseers/eelbot/internal/msg"
)

// NewBot creates a new bot instance.
func NewBot(ctx NewBotCtx) (bot *Bot, err error) {
	genSettings, err := ctx.Cfg.GetGeneralSettings()
	if err != nil {
		return
	}

	jokeSettings, err := ctx.Cfg.GetJokeSettings()
	if err != nil {
		return
	}

	dbSettings, err := ctx.Cfg.GetDatabaseSettings()
	if err != nil {
		return
	}

	tauntSettings, err := ctx.Cfg.GetTauntSettings()
	if err != nil {
		return
	}

	dg, err := discordgo.New("Bot " + ctx.Token)
	if err != nil {
		return
	}

	interpreter, err := msg.NewInterpreter(msg.NewInterpreterCtx{
		MsgTimeout:         genSettings.MsgTimeout,
		MultiLineJokeDelay: jokeSettings.MultiLineJokeDelay,
		SQLiteDB:           dbSettings.DBName,
		TauntsFolder:       tauntSettings.TauntsFolder,
	})
	if err != nil {
		return
	}

	bot = &Bot{
		dgBot:          dg,
		msgInterpreter: interpreter,
	}

	bot.dgBot.AddHandler(bot.guildCreateHandler)
	bot.dgBot.AddHandler(bot.messageCreateHandler)
	return
}

// Start starts the bot and open a connection to Discord.
func (bot *Bot) Start() (err error) {
	err = bot.dgBot.Open()
	return
}

// UpdateStatus updates the bot's status.
func (bot *Bot) UpdateStatus(game string) {
	bot.dgBot.UpdateStatus(0, game)
}

// SendMsg sends the given message to the given channel.
func (bot *Bot) SendMsg(channelID, msg string) {
	bot.dgBot.ChannelMessageSend(channelID, msg)
}

// SendFile sends the given file to the given channel.
func (bot *Bot) SendFile(channelID, filename string, file io.Reader) {
	bot.dgBot.ChannelFileSend(channelID, filename, file)
}

// DeleteMsg deletes the given message ID from the given channel.
func (bot *Bot) DeleteMsg(channelID, msgID string) {
	bot.dgBot.ChannelMessageDelete(channelID, msgID)
}

// Stop stops the bot and closes any open connections.
func (bot *Bot) Stop() (err error) {
	err = bot.msgInterpreter.Stop()
	if err != nil {
		return
	}

	err = bot.dgBot.Close()
	return
}
