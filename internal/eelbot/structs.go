package eelbot

import (
	"github.com/Emseers/Eelbot/internal/config"
	"github.com/Emseers/Eelbot/internal/msg"
	"github.com/bwmarrin/discordgo"
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
