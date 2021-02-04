package config

import (
	"fmt"
	"time"

	"gopkg.in/ini.v1"
)

// LoadConfig loads the given config file and returns a Reader to access various config settings.
func LoadConfig(file string) (cfg Reader, err error) {
	iniFile, err := ini.Load(file)
	if err != nil {
		return
	}

	cfg = &reader{
		file: iniFile,
	}
	return
}

func (cfg *reader) GetGeneralSettings() (settings GeneralSettings, err error) {
	msgTimeoutSecs, err := cfg.file.Section(sectionGeneral).Key(keyMsgTimeout).Uint()
	if err != nil {
		return
	}

	settings = GeneralSettings{
		MsgTimeout: time.Duration(uint(time.Second) * msgTimeoutSecs),
	}
	return
}

func (cfg *reader) GetJokeSettings() (settings JokeSettings, err error) {
	multiLineJokeDelaySecs, err := cfg.file.Section(sectionJokes).Key(keyMultiLineJokeDelay).Uint()
	if err != nil {
		return
	}

	settings = JokeSettings{
		MultiLineJokeDelay: time.Duration(uint(time.Second) * multiLineJokeDelaySecs),
	}
	return
}

func (cfg *reader) GetTauntSettings() (settings TauntSettings, err error) {
	tauntsFolder := cfg.file.Section(sectionTaunts).Key(keyTauntsFolder).String()
	if tauntsFolder == "" {
		err = fmt.Errorf("%s setting not found", keyTauntsFolder)
		return
	}

	settings = TauntSettings{
		TauntsFolder: tauntsFolder,
	}
	return
}

func (cfg *reader) GetDatabaseSettings() (settings DatabaseSettings, err error) {
	sqliteDbName := cfg.file.Section(sectionDatabase).Key(keyDbName).String()
	if sqliteDbName == "" {
		err = fmt.Errorf("%s setting not found", keyDbName)
		return
	}

	settings = DatabaseSettings{
		DBName: sqliteDbName,
	}
	return
}
