package eelbot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/emseers/eelbot/internal/config"
	"github.com/emseers/eelbot/internal/msg"
)

// NewBotCtx is the context needed to create an instance of the bot.
type NewBotCtx struct {
	Token string
	Cfg   config.Reader
}

// A Bot is an instance of a discord bot that can listen for commands and do various things.
type Bot struct {
	dgBot          *discordgo.Session
	msgInterpreter *msg.Interpreter
}
