package msg

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// NewInterpreter creates a new message interpretor.
func NewInterpreter(ctx NewInterpreterCtx) (interpreter *Interpreter, err error) {
	db, err := sql.Open("sqlite3", ctx.SQLiteDB)
	if err != nil {
		return
	}

	interpreter = &Interpreter{
		msgTimeout:         ctx.MsgTimeout,
		multiLineJokeDelay: ctx.MultiLineJokeDelay,
		sqliteDB:           db,
		tauntsFolder:       ctx.TauntsFolder,
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
func (interpreter *Interpreter) ParseAndReply(channelID, msgID, msg string, callbacks CallBackCtx) {
	// Check if entire message is uppercase.
	msgCleaned := toValidASCII(msg)
	msgCleanedAlphabetsOnly := toAlphabetOnly(msgCleaned)
	if len(msgCleanedAlphabetsOnly) >= 5 && msgCleanedAlphabetsOnly == strings.ToUpper(msgCleanedAlphabetsOnly) {
		if !interpreter.flagAllCaps.hasChannel(channelID) {
			interpreter.flagAllCaps.addChannelWithTimedReset(channelID, interpreter.msgTimeout)
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
						time.Sleep(interpreter.multiLineJokeDelay)
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
						time.Sleep(interpreter.multiLineJokeDelay)
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
					eelPic, err := interpreter.getEelPic()
					if err != nil {
						callbacks.SendMsg(channelID, err.Error())
						return
					}

					callbacks.SendFile(channelID, "eel.png", eelPic)
					return
				} else if len(cmdSlices) > 2 {
					if cmdSlices[1] == "bomb" {
						if num, err := strconv.ParseUint(cmdSlices[2], 10, 0); err == nil {
							eelPics, err := interpreter.getEelBomb(num)
							if err != nil {
								callbacks.SendMsg(channelID, err.Error())
								return
							}

							for _, eelPic := range eelPics {
								callbacks.SendFile(channelID, "eel.png", eelPic)
							}
							return
						}
					}
				}
			}
		case "flip":
			callbacks.SendMsg(channelID, flip())
			return
		case "help":
			callbacks.SendMsg(channelID, "The following commands are available for use:\n"+
				"```\n"+
				"/badjoke me             : Tell a joke\n"+
				"/channel                : Get channel info\n"+
				"/eel me                 : Post an eel pic\n"+
				"/eel bomb <n>           : Post n eel pics\n"+
				"/flip                   : Flip a coin\n"+
				"/help                   : Display this message\n"+
				"/ping                   : Pong\n"+
				"/play <game>            : Play the given game\n"+
				"/say <msg>              : Say the given message\n"+
				"/saychan <chan> <msg>   : Say the given message in the given channel\n"+
				"/taunt <tauntID>        : Post a taunt given a taunt ID\n"+
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
				tauntNum, err := strconv.Atoi(cmdSlices[1])
				if err == nil {
					taunt, fileName, err := interpreter.getTaunt(tauntNum)
					if err != nil {
						callbacks.SendMsg(channelID, err.Error())
						return
					}

					callbacks.SendFile(channelID, fileName, taunt)
				}
			}
		}
	}

	// Check if the message is a questionmark.
	if strings.HasPrefix(msgLowerCase, "?") {
		if !interpreter.flagQuestion.hasChannel(channelID) {
			interpreter.flagQuestion.addChannelWithTimedReset(channelID, interpreter.msgTimeout)
			callbacks.SendMsg(channelID, "Don't you questionmark me")
			return
		}
	}

	// Check is message is a hello greet.
	if isHelloGreet(msgLowerCase) {
		if !interpreter.flagHello.hasChannel(channelID) {
			interpreter.flagHello.addChannelWithTimedReset(channelID, interpreter.msgTimeout)
			callbacks.SendMsg(channelID, helloGreet())
			return
		}
	}

	// Check is message is a goodbye greet.
	if isGoodbyeGreet(msgLowerCase) {
		if !interpreter.flagGoodbye.hasChannel(channelID) {
			interpreter.flagGoodbye.addChannelWithTimedReset(channelID, interpreter.msgTimeout)
			callbacks.SendMsg(channelID, goodbyeGreet())
			return
		}
	}

	// Check is message is a laugh.
	if isLaugh(msgLowerCase) {
		callbacks.SendMsg(channelID, "lol")
		return
	}
}

// Stop stops the interpreter.
func (interpreter *Interpreter) Stop() (err error) {
	err = interpreter.sqliteDB.Close()
	return
}
