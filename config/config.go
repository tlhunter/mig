package config

import "errors"

// priority:
// CLI Flags
// Env Vars
// File Config

const (
	DEF_MIG_DIR = "./migrations"
)

type MigConfig struct {
	Connection string // DB connection string
	Migrations string // migrations directory, e.g. ./migrations
}

func GetConfig() (MigConfig, error) {
	flagConfig, _ := Flags()
	envConfig, _ := Environment()
	fileConfig, _ := File()

	config := MigConfig{}

	if flagConfig.Connection != "" {
		config.Connection = flagConfig.Connection
	} else if envConfig.Connection != "" {
		config.Connection = envConfig.Connection
	} else if fileConfig.Connection != "" {
		config.Connection = fileConfig.Connection
	} else {
		return config, errors.New("unable to determinte server connection")
	}

	if flagConfig.Migrations != "" {
		config.Migrations = flagConfig.Migrations
	} else if envConfig.Migrations != "" {
		config.Migrations = envConfig.Migrations
	} else if fileConfig.Migrations != "" {
		config.Migrations = fileConfig.Migrations
	} else {
		config.Migrations = DEF_MIG_DIR // TODO: combine with CWD for absolute path
	}

	return config, nil
}
