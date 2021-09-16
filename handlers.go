package eelbot

import (
	"github.com/bwmarrin/discordgo"
	"github.com/emseers/eelbot/msg"
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

	bot.msgInterpreter.ParseAndReply(event.ChannelID, event.ID, event.Content, msg.Callbacks{
		UpdateStatus: bot.UpdateStatus,
		SendMsg:      bot.SendMsg,
		SendFile:     bot.SendFile,
		DeleteMsg:    bot.DeleteMsg,
	})
}
