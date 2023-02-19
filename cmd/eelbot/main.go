// Package main is the entrypoint of the application.
package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/emseers/eelbot"
	"github.com/emseers/eelbot/commands"
	"github.com/emseers/eelbot/replies"
	_ "github.com/jackc/pgx/v5/stdlib"
	"gopkg.in/yaml.v3"
)

func mustDo(err error) {
	if err != nil {
		panic(err)
	}
}

func must[val any](v val, err error) val {
	mustDo(err)
	return v
}

var (
	configFile string
	token      string
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

	if configFile == "" {
		for _, configFile = range []string{"config.json", "config.yml", "config.yaml"} {
			if _, err := os.Stat(configFile); err == nil {
				break
			}
		}
	}

	s := must(eelbot.NewSession(token))
	bot := eelbot.New(s)
	config := must(os.ReadFile(configFile))

	var opts map[string]any
	switch strings.ToLower(path.Ext(configFile)) {
	case ".json":
		mustDo(json.Unmarshal(config, &opts))
	case ".yml", ".yaml":
		yOpts := map[any]any{}
		mustDo(yaml.Unmarshal(config, &yOpts))

		// Unlike json, yaml allows for non-string keys. Therefore, convert the map[any]any to a map[string]any to be
		// json compatible.
		opts = toStrKeys(yOpts).(map[string]any)

		// Since the yaml and json packages unmarshal numeric types differently, stick to the json way (always float64)
		// by marshalling to json and unmarshalling back.
		config = must(json.Marshal(opts))
		mustDo(json.Unmarshal(config, &opts))
	}

	var (
		db        *sql.DB
		dbTimeout = 5 * time.Second
	)

	if dbOpts, ok := opts["database"].(map[string]any); ok {
		db = must(sql.Open("pgx", connectionString(dbOpts)))
		defer func() { mustDo(db.Close()) }()

		if timeoutSecs, ok2 := opts["timeout"].(float64); ok2 {
			dbTimeout = time.Second * time.Duration(timeoutSecs)
		}
	}

	if cmdOpts, ok := opts["commands"].(map[string]any); ok {
		mustDo(commands.Register(bot, cmdOpts, db, dbTimeout))
	}

	if replyOpts, ok := opts["replies"].(map[string]any); ok {
		mustDo(replies.Register(bot, replyOpts))
	}

	mustDo(bot.Start())
	defer func() { mustDo(bot.Stop()) }()

	fmt.Println("eelbot started up successfully")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-sc
	fmt.Println()
	fmt.Println("goodbye")
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

// Gets a string param from opts by checking if there exists a key for the prefix itself, <prefix>_env or <prefix>_file
// (in that order) and reading the value appropriately.
func getParam(opts map[string]any, prefix string) (param string) {
	param, _ = opts[prefix].(string)
	param = strings.TrimSpace(param)
	if param == "" {
		if paramEnv, _ := opts[prefix+"_env"].(string); paramEnv != "" {
			param, _ = os.LookupEnv(paramEnv)
			param = strings.TrimSpace(param)
		}
	}
	if param == "" {
		if paramFile, _ := opts[prefix+"_file"].(string); paramFile != "" {
			paramBytes := must(os.ReadFile(paramFile))
			param = strings.TrimSpace(string(paramBytes))
		}
	}
	return
}

// Creates a database connection string based on the given opts.
func connectionString(opts map[string]any) string {
	u := &url.URL{Scheme: "postgresql"}
	u.Host = getParam(opts, "host")
	u.Path = getParam(opts, "database")
	username := getParam(opts, "username")
	password := getParam(opts, "password")
	if username != "" {
		if password != "" {
			u.User = url.UserPassword(username, password)
		} else {
			u.User = url.User(username)
		}
	}
	return u.String()
}
