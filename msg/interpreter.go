package msg

import (
	"database/sql"
	"io"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// UpdateStatusFunc is a callback function that will be used by the Interpreter to update the bot's status.
type UpdateStatusFunc func(game string)

// SendMsgFunc is a callback function that will be used by the Interpreter to send messages.
type SendMsgFunc func(channelID, msg string)

// SendFileFunc is a callback function that will be used by the Interpreter to send files.
type SendFileFunc func(channelID, filename string, file io.Reader)

// DeleteMsgFunc is a callback function that will be used by the Interpreter to delete messages.
type DeleteMsgFunc func(channelID, msgID string)

// Callbacks is a container to hold all callback functions for the message interpreter.
type Callbacks struct {
	UpdateStatus UpdateStatusFunc
	SendMsg      SendMsgFunc
	SendFile     SendFileFunc
	DeleteMsg    DeleteMsgFunc
}

// An Interpreter can parse and respond to certain messages.
type Interpreter struct {
	MsgTimeout         time.Duration
	MultiLineJokeDelay time.Duration

	db           *sql.DB
	flagAllCaps  *flagMap
	flagQuestion *flagMap
	flagHello    *flagMap
	flagGoodbye  *flagMap
}

// NewInterpreter creates a new message interpretor.
func NewInterpreter(dbName string) (interpreter *Interpreter, err error) {
	var db *sql.DB
	if db, err = sql.Open("sqlite3", dbName); err != nil {
		return
	}

	interpreter = &Interpreter{
		MsgTimeout:         5 * time.Second,
		MultiLineJokeDelay: 5 * time.Second,
		db:                 db,
		flagAllCaps:        newFlagMap(),
		flagQuestion:       newFlagMap(),
		flagHello:          newFlagMap(),
		flagGoodbye:        newFlagMap(),
	}
	return
}

// GetWelcomeMsg returns a suitable welcome message to print on initial joins to guilds.
func (interpreter *Interpreter) GetWelcomeMsg() (msg string) {
	msg = "I am now ready! Type /help to see what I can do."
	return
}

// ParseAndReply parses an incoming message from a channel and replies if necessary.
func (interpreter *Interpreter) ParseAndReply(channelID, msgID, msg string, callbacks Callbacks) {
	// Check if entire message is uppercase.
	msgCleaned := toValidASCII(msg)
	msgCleanedAlphabetsOnly := toAlphabetOnly(msgCleaned)
	if len(msgCleanedAlphabetsOnly) >= 5 && msgCleanedAlphabetsOnly == strings.ToUpper(msgCleanedAlphabetsOnly) {
		if !interpreter.flagAllCaps.hasChannel(channelID) {
			interpreter.flagAllCaps.addChannelWithTimedReset(channelID, interpreter.MsgTimeout)
			callbacks.SendMsg(channelID, yellResponse())
			return
		}
	}

	// Convert message to lowercase for parsing.
	msgLowerCase := strings.ToLower(msg)

	// Check if the message is a command.
	if strings.HasPrefix(msgLowerCase, "/") {
		cmd := msgLowerCase[1:]
		cmdSlices := strings.Split(cmd, " ")

		switch cmdSlices[0] {
		case "badjoke":
			if len(cmdSlices) > 1 {
				if cmdSlices[1] == "me" {
					partOne, partTwo, err := interpreter.getJoke()
					if err != nil {
						callbacks.SendMsg(channelID, err.Error())
						return
					}

					callbacks.SendMsg(channelID, partOne)
					if partTwo.Valid {
						time.Sleep(interpreter.MultiLineJokeDelay)
						callbacks.SendMsg(channelID, partTwo.String)
					}
					return
				} else if num, err := strconv.ParseUint(cmdSlices[1], 10, 0); err == nil {
					partOne, partTwo, err := interpreter.getSpecificJoke(num)
					if err != nil {
						callbacks.SendMsg(channelID, err.Error())
						return
					}

					callbacks.SendMsg(channelID, partOne)
					if partTwo.Valid {
						time.Sleep(interpreter.MultiLineJokeDelay)
						callbacks.SendMsg(channelID, partTwo.String)
					}
					return
				}
			}
		case "channel":
			callbacks.SendMsg(channelID, "Channel ID: "+channelID)
		case "eel":
			if len(cmdSlices) > 1 {
				if cmdSlices[1] == "me" {
					image, filename, err := interpreter.getImage()
					if err != nil {
						callbacks.SendMsg(channelID, err.Error())
						return
					}

					callbacks.SendFile(channelID, filename, image)
					return
				} else if num, err := strconv.ParseUint(cmdSlices[1], 10, 0); err == nil {
					image, filename, err := interpreter.getSpecificImage(num)
					if err != nil {
						callbacks.SendMsg(channelID, err.Error())
						return
					}

					callbacks.SendFile(channelID, filename, image)
					return
				}
			}
		case "flip":
			callbacks.SendMsg(channelID, flip())
			return
		case "help":
			callbacks.SendMsg(channelID, "The following commands are available for use:\n"+
				"```\n"+
				"/badjoke me             : Tell a joke\n"+
				"/badjoke <id>           : Tell a specific joke\n"+
				"/channel                : Get channel info\n"+
				"/eel me                 : Post an image\n"+
				"/eel <id>               : Post a specific image\n"+
				"/flip                   : Flip a coin\n"+
				"/help                   : Display this message\n"+
				"/ping                   : Pong\n"+
				"/play <game>            : Play the given game\n"+
				"/say <msg>              : Say the given message\n"+
				"/saychan <chan> <msg>   : Say the given message in the given channel\n"+
				"/taunt me               : Post a taunt\n"+
				"/taunt <id>             : Post a specific taunt\n"+
				"```")
			return
		case "ping":
			callbacks.SendMsg(channelID, "Pong")
			return
		case "play":
			callbacks.DeleteMsg(channelID, msgID)

			if len(msg) <= 6 {
				callbacks.UpdateStatus("")
				return
			}

			callbacks.UpdateStatus(msg[6:])
			return
		case "say":
			if len(cmdSlices) > 1 {
				if len(cmdSlices) != 2 || cmdSlices[1] != "lol" {
					callbacks.DeleteMsg(channelID, msgID)
					callbacks.SendMsg(channelID, msg[len("/say "):])
					return
				}
			}
		case "saychan":
			if len(cmdSlices) > 2 {
				if len(cmdSlices) != 3 || cmdSlices[2] != "lol" {
					callbacks.DeleteMsg(channelID, msgID)
					callbacks.SendMsg(cmdSlices[1], msg[len("/saychan "+cmdSlices[1]+" "):])
					return
				}
			}
		case "taunt":
			if len(cmdSlices) > 1 {
				if cmdSlices[1] == "me" {
					taunt, filename, err := interpreter.getTaunt()
					if err != nil {
						callbacks.SendMsg(channelID, err.Error())
						return
					}

					callbacks.SendFile(channelID, filename, taunt)
					return
				} else if num, err := strconv.ParseUint(cmdSlices[1], 10, 0); err == nil {
					taunt, filename, err := interpreter.getSpecificTaunt(num)
					if err != nil {
						callbacks.SendMsg(channelID, err.Error())
						return
					}

					callbacks.SendFile(channelID, filename, taunt)
					return
				}
			}
		}
	}

	// Check if the message is a questionmark.
	if strings.HasPrefix(msgLowerCase, "?") {
		if !interpreter.flagQuestion.hasChannel(channelID) {
			interpreter.flagQuestion.addChannelWithTimedReset(channelID, interpreter.MsgTimeout)
			callbacks.SendMsg(channelID, "Don't you questionmark me")
			return
		}
	}

	// Check if message is a hello greet.
	if isHelloGreet(msgLowerCase) {
		if !interpreter.flagHello.hasChannel(channelID) {
			interpreter.flagHello.addChannelWithTimedReset(channelID, interpreter.MsgTimeout)
			callbacks.SendMsg(channelID, helloGreet())
			return
		}
	}

	// Check if message is a goodbye greet.
	if isGoodbyeGreet(msgLowerCase) {
		if !interpreter.flagGoodbye.hasChannel(channelID) {
			interpreter.flagGoodbye.addChannelWithTimedReset(channelID, interpreter.MsgTimeout)
			callbacks.SendMsg(channelID, goodbyeGreet())
			return
		}
	}

	// Check if message is a laugh.
	if isLaugh(msgLowerCase) {
		callbacks.SendMsg(channelID, "lol")
		return
	}
}

// Stop stops the interpreter.
func (interpreter *Interpreter) Stop() (err error) {
	err = interpreter.db.Close()
	return
}
