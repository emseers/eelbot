package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/emseers/eelbot/internal/config"
	"github.com/emseers/eelbot/internal/eelbot"
)

func main() {
	var (
		configFile string
		token      string
	)

	flag.StringVar(&configFile, "c", "config.ini", "Config file")
	flag.StringVar(&token, "t", "", "Bot token")
	flag.Parse()

	if token == "" {
		fmt.Fprintln(os.Stderr, "no token provided; please use -t to provide bot token")
		os.Exit(1)
	}

	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	eelbot, err := eelbot.NewBot(eelbot.NewBotCtx{
		Token: token,
		Cfg:   cfg,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = eelbot.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("Eelbot started up successfully.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	<-sc
	err = eelbot.Stop()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
