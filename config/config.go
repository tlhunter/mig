package config

import "errors"

// priority:
// CLI Flags
// Env Vars
// File Config

type MigConfig struct {
	connection string // DB connection string
	migrations string // migrations directory, e.g. ./migrations
}

func GetConfig() (MigConfig, error) {
	flagConfig, _ := Flags()
	envConfig, _ := Environment()
	fileConfig, _ := File()

	config := MigConfig{}

	if flagConfig.connection != "" {
		config.connection = flagConfig.connection
	} else if envConfig.connection != "" {
		config.connection = envConfig.connection
	} else if fileConfig.connection != "" {
		config.connection = fileConfig.connection
	} else {
		return config, errors.New("unable to determinte server connection")
	}

	if flagConfig.migrations != "" {
		config.migrations = flagConfig.migrations
	} else if envConfig.migrations != "" {
		config.migrations = envConfig.migrations
	} else if fileConfig.migrations != "" {
		config.migrations = fileConfig.migrations
	} else {
		return config, errors.New("unable to determinte migrations directory")
	}

	return config, nil
}
