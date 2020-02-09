package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/Emseers/Eelbot/msg"

	"github.com/bwmarrin/discordgo"
)

var token string
var buffer = make([][]byte, 0)

var timeout = 10 * time.Second
var flagYl bool
var flagQm bool
var flagMg bool
var flagEg bool

var multiLineJokeDelay = 3 * time.Second

func setFlag(flag *bool) {
	time.Sleep(timeout)
	*flag = true
}

func init() {
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	if token == "" {
		fmt.Println("No token provided. Please run: Eelbot -t <bot token>")
		return
	}

	// Set flags
	flagYl = true
	flagQm = true
	flagMg = true
	flagEg = true

	// Create a new Discord session using the provided bot token
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("Error creating Discord session: ", err)
		return
	}
	defer dg.Close()

	// Register ready as a callback for the ready events
	dg.AddHandler(ready)

	// Register messageCreate as a callback for the messageCreate events
	dg.AddHandler(messageCreate)

	// Register guildCreate as a callback for the guildCreate events
	dg.AddHandler(guildCreate)

	// Open the websocket and begin listening
	err = dg.Open()
	if err != nil {
		fmt.Println("Error opening Discord session: ", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received
	fmt.Println("Eelbot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}

// This function will be called (due to AddHandler above) when the bot receives the "ready" event from Discord
func ready(s *discordgo.Session, event *discordgo.Ready) {
	// Set the playing status.
	s.UpdateStatus(0, "")
}

// This function will be called (due to AddHandler above) every time a new message is created that the bot has access to
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Remove emoji and other random characters and check if entire message is uppercase
	content := msg.ToValidAscii(m.Content)
	alphaContent := msg.ToAlphabetOnly(content)
	if alphaContent == strings.ToUpper(alphaContent) && len(alphaContent) >= 5 {
		if flagYl {
			s.ChannelMessageSend(m.ChannelID, msg.YellResponse())
			flagYl = false
			go setFlag(&flagYl)
			return
		}
	}

	// Convert message to lowercase and parse
	content = strings.ToLower(content)
	if strings.HasPrefix(content, "/") {
		cmd := content[1:]
		cmdSlices := strings.Split(cmd, " ")

		switch cmdSlices[0] {
		case "badjoke":
			if len(cmdSlices) > 1 {
				if cmdSlices[1] == "me" {
					partOne, partTwo, err := msg.Joke()
					if err != nil {
						s.ChannelMessageSend(m.ChannelID, err.Error())
					} else {
						s.ChannelMessageSend(m.ChannelID, partOne)
						time.Sleep(multiLineJokeDelay)
						s.ChannelMessageSend(m.ChannelID, partTwo)
					}
				} else if _, err := strconv.Atoi(cmdSlices[1]); err == nil {
					partOne, partTwo, err := msg.JokeSpecific(cmdSlices[1])
					if err != nil {
						s.ChannelMessageSend(m.ChannelID, err.Error())
					} else {
						s.ChannelMessageSend(m.ChannelID, partOne)
						time.Sleep(multiLineJokeDelay)
						s.ChannelMessageSend(m.ChannelID, partTwo)
					}
				}
			}
		case "channel":
			s.ChannelMessageSend(m.ChannelID, "Channel ID: "+m.ChannelID)
		case "eel":
			if len(cmdSlices) > 1 {
				if cmdSlices[1] == "me" {
					eelPic, err := msg.EelPic()
					if err != nil {
						s.ChannelMessageSend(m.ChannelID, err.Error())
					} else {
						s.ChannelFileSend(m.ChannelID, "eel.png", eelPic)
					}
				}
			}
		case "flip":
			s.ChannelMessageSend(m.ChannelID, msg.Flip())
		case "help":
			s.ChannelMessageSend(m.ChannelID, "The following commands are available for use:\n"+
				"```\n"+
				"/badjoke me             : Tell a joke\n"+
				"/channel                : Get channel info\n"+
				"/eel me                 : Post an eel pic\n"+
				"/flip                   : Flip a coin\n"+
				"/help                   : Display this message\n"+
				"/ping                   : Pong\n"+
				"/play <game>            : Play the given game\n"+
				"/say <msg>              : Say the given message\n"+
				"/saychan <chan> <msg>   : Say the given message in the given channel\n"+
				"/taunt <tauntID>        : Post a taunt given a taunt ID\n"+
				"```")
		case "ping":
			s.ChannelMessageSend(m.ChannelID, "Pong")
		case "play":
			if len(m.Content) > 6 {
				s.UpdateStatus(0, m.Content[6:])
			} else {
				s.UpdateStatus(0, "")
			}
		case "say":
			if len(cmdSlices) > 1 {
				s.ChannelMessageDelete(m.ChannelID, m.ID)
				s.ChannelMessageSend(m.ChannelID, strings.Join(cmdSlices[1:], " "))
			}
		case "saychan":
			if len(cmdSlices) > 2 {
				s.ChannelMessageDelete(m.ChannelID, m.ID)
				s.ChannelMessageSend(cmdSlices[1], strings.Join(cmdSlices[2:], " "))
			}
		case "taunt":
			if len(cmdSlices) > 1 {
				taunt, err := strconv.Atoi(cmdSlices[1])
				taunt--
				if err == nil {
					reader, fileName, err := msg.PlayTaunt(taunt)
					if err == nil {
						s.ChannelFileSend(m.ChannelID, fileName, reader)
					} else {
						s.ChannelMessageSend(m.ChannelID, err.Error())
					}
				}
			}
		}
	} else if strings.HasPrefix(content, "?") {
		if flagQm {
			s.ChannelMessageSend(m.ChannelID, "Don't you questionmark me")
			flagQm = false
			go setFlag(&flagQm)
		}
	} else if msg.IsMorningGreet(content) {
		if flagMg {
			s.ChannelMessageSend(m.ChannelID, msg.MorningGreet())
			flagMg = false
			go setFlag(&flagMg)
		}
	} else if msg.IsGoodbyeGreet(content) {
		if flagEg {
			s.ChannelMessageSend(m.ChannelID, msg.GoodbyeGreet())
			flagEg = false
			go setFlag(&flagEg)
		}
	}
}

// This function will be called (due to AddHandler above) every time a new guild is joined.
func guildCreate(s *discordgo.Session, event *discordgo.GuildCreate) {
	if event.Guild.Unavailable {
		return
	}

	for _, channel := range event.Guild.Channels {
		if channel.ID == event.Guild.ID {
			_, _ = s.ChannelMessageSend(channel.ID, "Eelbot is ready! Type /help to see the list of commands.")
			return
		}
	}
}
