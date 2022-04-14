package main

import (
	"database/sql"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/emseers/eelbot"
	"github.com/emseers/eelbot/commands"
	"github.com/emseers/eelbot/replies"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/ini.v1"
)

var (
	configFile = "config.ini"
	token      = ""
)

func init() {
	flag.StringVar(&configFile, "c", configFile, "Config file")
	flag.StringVar(&token, "t", token, "Bot token")

	rand.Seed(time.Now().UnixNano())
}

func main() {
	flag.Parse()

	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, r)
			os.Exit(1)
		}
	}()

	if token == "" {
		panic("no token provided; please use -t to provide bot token")
	}

	bot, err := eelbot.New(token)
	if err != nil {
		panic(err)
	}

	var opts *ini.File
	if opts, err = ini.InsensitiveLoad(configFile); err != nil {
		panic(err)
	}

	var db *sql.DB
	if dbName := opts.Section("Database").Key("name").String(); dbName != "" {
		if db, err = sql.Open("sqlite3", dbName); err != nil {
			panic(err)
		}
	}

	if err = commands.Register(bot, opts.Section("Commands"), db); err != nil {
		panic(err)
	}

	if err = replies.Register(bot, opts.Section("Replies")); err != nil {
		panic(err)
	}

	if err = bot.Start(); err != nil {
		panic(err)
	}

	fmt.Println("eelbot started up successfully")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-sc
	fmt.Println()
	fmt.Println("goodbye")
	if err = bot.Stop(); err != nil {
		panic(err)
	}
}
