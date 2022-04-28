package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/emseers/eelbot"
	"github.com/emseers/eelbot/commands"
	"github.com/emseers/eelbot/replies"
	_ "github.com/mattn/go-sqlite3"
	"gopkg.in/yaml.v3"
)

var (
	configFile string
	token      string
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

	if configFile == "" {
		for _, configFile = range []string{"config.json", "config.yml", "config.yaml"} {
			if _, err := os.Stat(configFile); err == nil {
				break
			}
		}
	}

	s, err := eelbot.NewSession(token)
	if err != nil {
		panic(err)
	}
	bot := eelbot.New(s)

	var config []byte
	if config, err = os.ReadFile(configFile); err != nil {
		panic(err)
	}

	var opts map[string]any
	switch strings.ToLower(path.Ext(configFile)) {
	case ".json":
		if err = json.Unmarshal(config, &opts); err != nil {
			panic(err)
		}
	case ".yml", ".yaml":
		yOpts := map[any]any{}
		if err = yaml.Unmarshal(config, &yOpts); err != nil {
			panic(err)
		}

		// Unlike json, yaml allows for non-string keys. Therefore, convert the map[any]any to a map[string]any to be
		// json compatible.
		opts = toStrKeys(yOpts).(map[string]any)

		// Since the yaml and json packages unmarshal numeric types differently, stick to the json way (always float64)
		// by marshalling to json and unmarshalling back.
		if config, err = json.Marshal(opts); err != nil {
			panic(err)
		}
		if err = json.Unmarshal(config, &opts); err != nil {
			panic(err)
		}
	}

	var db *sql.DB
	if dbOpts, ok := opts["database"].(map[string]any); ok {
		if dbName, ok2 := dbOpts["name"].(string); ok2 {
			if db, err = sql.Open("sqlite3", dbName); err != nil {
				panic(err)
			}
		}
	}

	if cmdOpts, ok := opts["commands"].(map[string]any); ok {
		if err = commands.Register(bot, cmdOpts, db); err != nil {
			panic(err)
		}
	}

	if replyOpts, ok := opts["replies"].(map[string]any); ok {
		if err = replies.Register(bot, replyOpts); err != nil {
			panic(err)
		}
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

// Converts a map[any]any to a map[string]any.
func toStrKeys(a any) any {
	switch ta := a.(type) {
	case map[any]any:
		m := map[string]any{}
		for k, v := range ta {
			key, ok := k.(string)
			if !ok {
				panic(fmt.Errorf("yaml: invalid map key: %#v", k))
			}
			m[key] = toStrKeys(v)
		}
		return m
	case map[string]any:
		m := map[string]any{}
		for k, v := range ta {
			m[k] = toStrKeys(v)
		}
		return m
	case []any:
		var s = make([]any, len(ta))
		for i, v := range ta {
			s[i] = toStrKeys(v)
		}
		return s
	default:
		return ta
	}
}
