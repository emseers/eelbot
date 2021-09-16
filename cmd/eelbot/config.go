package main

import (
	"fmt"
	"time"

	"github.com/emseers/eelbot"
	"gopkg.in/ini.v1"
)

const (
	sectionGeneral  = "General"
	sectionJokes    = "Jokes"
	sectionDatabase = "Database"

	keyMsgTimeout         = "msg_timeout"
	keyMultiLineJokeDelay = "multi_line_joke_delay"
	keyDbName             = "db_name"
)

func parseConfig(filename string) (opts *eelbot.Options, err error) {
	var file *ini.File
	if file, err = ini.Load(filename); err != nil {
		return
	}

	var msgTimeoutSecs uint
	if msgTimeoutSecs, err = file.Section(sectionGeneral).Key(keyMsgTimeout).Uint(); err != nil {
		return
	}

	var multiLineJokeDelaySecs uint
	if multiLineJokeDelaySecs, err = file.Section(sectionJokes).Key(keyMultiLineJokeDelay).Uint(); err != nil {
		return
	}

	dbName := file.Section(sectionDatabase).Key(keyDbName).String()
	if dbName == "" {
		err = fmt.Errorf("%s setting not found", keyDbName)
		return
	}

	opts = &eelbot.Options{
		MsgTimeout:         time.Second * time.Duration(msgTimeoutSecs),
		MultiLineJokeDelay: time.Second * time.Duration(multiLineJokeDelaySecs),
		DBName:             dbName,
	}
	return
}
