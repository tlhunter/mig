package config

import "errors"

const (
	DEF_MIG_DIR = "./migrations"
)

type MigConfig struct {
	Connection string // DB connection string
	Migrations string // migrations directory, e.g. ./migrations
	MigRcPath  string // override path to config file
}

func GetConfig() (MigConfig, []string, error) {
	config := MigConfig{}

	flagConfig, subcommands, _ := GetConfigFromProcessFlags()

	err := SetEnvFromConfigFile(flagConfig.MigRcPath) // reads .env and sets env vars but does not override

	if err != nil {
		return config, []string{}, err
	}

	envConfig, _ := GetConfigFromEnvVars()

	if flagConfig.Connection != "" {
		config.Connection = flagConfig.Connection
	} else if envConfig.Connection != "" {
		config.Connection = envConfig.Connection
	} else {
		return config, subcommands, errors.New("unable to determinte server connection")
	}

	if flagConfig.Migrations != "" {
		config.Migrations = flagConfig.Migrations
	} else if envConfig.Migrations != "" {
		config.Migrations = envConfig.Migrations
	} else {
		config.Migrations = DEF_MIG_DIR
	}

	return config, subcommands, nil
}
