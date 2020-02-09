package msg

import (
	"gopkg.in/ini.v1"
)

var sqliteDBString string
var tauntsFolder string

func init() {
	cfg, err := ini.Load("config.ini")
	if err == nil {
		sqliteDBString = cfg.Section("Database").Key("db_name").String()
		tauntsFolder = cfg.Section("Taunts").Key("taunts_folder").String()
	}

	// Fallback to default database name if config is invalid or is missing config setting
	if sqliteDBString == "" {
		sqliteDBString = "EelbotDB.db"
	}

	if tauntsFolder == "" {
		tauntsFolder = "taunts"
	}
}
