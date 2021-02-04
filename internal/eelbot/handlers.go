package eelbot

import (
	"github.com/Emseers/Eelbot/internal/msg"
	"github.com/bwmarrin/discordgo"
)

func (bot *Bot) guildCreateHandler(dg *discordgo.Session, event *discordgo.GuildCreate) {
	if event.Guild.Unavailable {
		return
	}

	for _, channel := range event.Guild.Channels {
		if channel.ID == event.Guild.ID {
			bot.SendMsg(channel.ID, bot.msgInterpreter.GetWelcomeMsg())
			return
		}
	}
}

func (bot *Bot) messageCreateHandler(dg *discordgo.Session, event *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself.
	if event.Author.ID == dg.State.User.ID {
		return
	}

	bot.msgInterpreter.ParseAndReply(event.ChannelID, event.ID, event.Content, msg.CallBackCtx{
		UpdateStatus: bot.UpdateStatus,
		SendMsg:      bot.SendMsg,
		SendFile:     bot.SendFile,
		DeleteMsg:    bot.DeleteMsg,
	})
}
