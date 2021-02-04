package config

import (
	"time"

	"gopkg.in/ini.v1"
)

// GeneralSettings is a container to hold all general settings.
type GeneralSettings struct {
	MsgTimeout time.Duration
}

// JokeSettings is a container to hold all joke related settings.
type JokeSettings struct {
	MultiLineJokeDelay time.Duration
}

// TauntSettings is a container to hold all taunt related settings.
type TauntSettings struct {
	TauntsFolder string
}

// DatabaseSettings is a container to hold all database related settings.
type DatabaseSettings struct {
	DBName string
}

// A Reader provides methods to access various config settings.
type Reader interface {
	// GetGeneralSettings returns all general settings.
	GetGeneralSettings() (GeneralSettings, error)

	// GetJokeSettings returns all joke related settings.
	GetJokeSettings() (JokeSettings, error)

	// GetTauntSettings returns all taunt related settings.
	GetTauntSettings() (TauntSettings, error)

	// GetDatabaseSettings returns all database related settings.
	GetDatabaseSettings() (DatabaseSettings, error)
}

type reader struct {
	file *ini.File
}
