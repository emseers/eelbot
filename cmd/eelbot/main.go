package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/emseers/eelbot"
)

var (
	configFile = "config.ini"
	token      = ""
)

func init() {
	flag.StringVar(&configFile, "c", configFile, "Config file")
	flag.StringVar(&token, "t", token, "Bot token")
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

	opts, err := parseConfig(configFile)
	if err != nil {
		panic(err)
	}
	opts.Token = token

	var bot *eelbot.Bot
	if bot, err = eelbot.New(opts); err != nil {
		panic(err)
	}

	if err = bot.Start(); err != nil {
		panic(err)
	}

	fmt.Println("eelbot started up successfully")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-sc
	if err = bot.Stop(); err != nil {
		panic(err)
	}
}
