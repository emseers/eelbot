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

	opts := parseConfig(configFile)
	bot, err := eelbot.New(token)
	if err != nil {
		panic(err)
	}

	dbEnabled := opts.Database.Name != ""
	var db *sql.DB
	if dbEnabled {
		if db, err = sql.Open("sqlite3", opts.Database.Name); err != nil {
			panic(err)
		}
	}

	// Commands.
	if opts.Commands.BadjokeEnable {
		if !dbEnabled {
			panic("/badjoke command requires a database")
		}
		if opts.Commands.BadjokeDelay == 0 {
			opts.Commands.BadjokeDelay = 3
		}
		bot.RegisterCommand("badjoke", *commands.JokeCommand(db, time.Second*time.Duration(opts.Commands.BadjokeDelay)))
	}
	if opts.Commands.EelEnable {
		if !dbEnabled {
			panic("/eel command requires a database")
		}
		bot.RegisterCommand("eel", *commands.ImageCommand(db))
	}
	if opts.Commands.TauntEnable {
		if !dbEnabled {
			panic("/taunt command requires a database")
		}
		bot.RegisterCommand("taunt", *commands.TauntCommand(db))
	}
	if opts.Commands.ChannelEnable {
		bot.RegisterCommand("channel", *commands.ChannelCommand())
	}
	if opts.Commands.FlipEnable {
		bot.RegisterCommand("flip", *commands.FlipCommand())
	}
	if opts.Commands.ListenEnable {
		bot.RegisterCommand("listen", *commands.ListenCommand())
	}
	if opts.Commands.PlayEnable {
		bot.RegisterCommand("play", *commands.PlayCommand())
	}
	if opts.Commands.PingEnable {
		bot.RegisterCommand("ping", *commands.PingCommand())
	}
	if opts.Commands.SayEnable {
		bot.RegisterCommand("say", *commands.SayCommand())
	}
	if opts.Commands.SaychanEnable {
		bot.RegisterCommand("saychan", *commands.SayChanCommand())
	}

	// Replies.
	if opts.Replies.CapsEnable {
		if opts.Replies.CapsMinLen == 0 {
			opts.Replies.CapsMinLen = 5
		}
		if opts.Replies.CapsPercent == 0 {
			opts.Replies.CapsPercent = 100
		}
		r := replies.CapsReply(opts.Replies.CapsMinLen, opts.Replies.CapsPercent)
		r.Timeout = time.Second * time.Duration(opts.Replies.CapsTimeout)
		bot.RegisterReply(*r)
	}
	if opts.Replies.HelloEnable {
		if opts.Replies.HelloPercent == 0 {
			opts.Replies.HelloPercent = 100
		}
		r := replies.HelloReply(opts.Replies.HelloPercent)
		r.Timeout = time.Second * time.Duration(opts.Replies.HelloTimeout)
		bot.RegisterReply(*r)
	}
	if opts.Replies.GoodbyeEnable {
		if opts.Replies.GoodbyePercent == 0 {
			opts.Replies.GoodbyePercent = 100
		}
		r := replies.GoodbyeReply(opts.Replies.GoodbyePercent)
		r.Timeout = time.Second * time.Duration(opts.Replies.GoodbyeTimeout)
		bot.RegisterReply(*r)
	}
	if opts.Replies.LaughEnable {
		if opts.Replies.LaughPercent == 0 {
			opts.Replies.LaughPercent = 100
		}
		r := replies.LaughReply(opts.Replies.LaughPercent)
		r.Timeout = time.Second * time.Duration(opts.Replies.LaughTimeout)
		bot.RegisterReply(*r)
	}
	if opts.Replies.QuestionEnable {
		if opts.Replies.QuestionPercent == 0 {
			opts.Replies.QuestionPercent = 100
		}
		r := replies.QuestionReply(opts.Replies.QuestionPercent)
		r.Timeout = time.Second * time.Duration(opts.Replies.QuestionTimeout)
		bot.RegisterReply(*r)
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
