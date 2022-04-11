package main

import (
	"gopkg.in/ini.v1"
)

type config struct {
	Database databaseConfig
	Commands commandsConfig
	Replies  repliesConfig
}

type databaseConfig struct {
	Name string `ini:"name"`
}

type commandsConfig struct {
	BadjokeEnable bool `ini:"badjoke_enable"`
	BadjokeDelay  int  `ini:"badjoke_delay"`
	EelEnable     bool `ini:"eel_enable"`
	TauntEnable   bool `ini:"taunt_enable"`
	ChannelEnable bool `ini:"channel_enable"`
	FlipEnable    bool `ini:"flip_enable"`
	ListenEnable  bool `ini:"listen_enable"`
	PlayEnable    bool `ini:"play_enable"`
	PingEnable    bool `ini:"ping_enable"`
	SayEnable     bool `ini:"say_enable"`
	SaychanEnable bool `ini:"saychan_enable"`
}

type repliesConfig struct {
	CapsEnable      bool `ini:"caps_enable"`
	CapsMinLen      int  `ini:"caps_min_len"`
	CapsPercent     int  `ini:"caps_percent"`
	CapsTimeout     int  `ini:"caps_timeout"`
	HelloEnable     bool `ini:"hello_enable"`
	HelloPercent    int  `ini:"hello_percent"`
	HelloTimeout    int  `ini:"hello_timeout"`
	GoodbyeEnable   bool `ini:"goodbye_enable"`
	GoodbyePercent  int  `ini:"goodbye_percent"`
	GoodbyeTimeout  int  `ini:"goodbye_timeout"`
	LaughEnable     bool `ini:"laugh_enable"`
	LaughPercent    int  `ini:"laugh_percent"`
	LaughTimeout    int  `ini:"laugh_timeout"`
	QuestionEnable  bool `ini:"question_enable"`
	QuestionPercent int  `ini:"question_percent"`
	QuestionTimeout int  `ini:"question_timeout"`
}

func parseConfig(path string) *config {
	file, err := ini.InsensitiveLoad(path)
	if err != nil {
		panic(err)
	}
	opts := new(config)
	if err = file.StrictMapTo(opts); err != nil {
		panic(err)
	}
	return opts
}
