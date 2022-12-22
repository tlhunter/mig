package config

import "errors"

const (
	DEF_MIG_DIR = "./migrations"
)

type MigConfig struct {
	Connection string // DB connection string
	Migrations string // migrations directory, e.g. ./migrations
}

func GetConfig() (MigConfig, error) {
	config := MigConfig{}

	err := SetEnvFromConfigFile() // reads .env and sets env vars but does not override

	if err != nil {
		return config, err
	}

	flagConfig, _ := GetConfigFromProcessFlags()
	envConfig, _ := GetConfigFromEnvVars()

	if flagConfig.Connection != "" {
		config.Connection = flagConfig.Connection
	} else if envConfig.Connection != "" {
		config.Connection = envConfig.Connection
	} else {
		return config, errors.New("unable to determinte server connection")
	}

	if flagConfig.Migrations != "" {
		config.Migrations = flagConfig.Migrations
	} else if envConfig.Migrations != "" {
		config.Migrations = envConfig.Migrations
	} else {
		config.Migrations = DEF_MIG_DIR // TODO: combine with CWD for absolute path
	}

	return config, nil
}
