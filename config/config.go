package config

import "github.com/tlhunter/mig/result"

const (
	DEF_MIG_DIR = "./migrations"
)

type MigConfig struct {
	Connection string // DB connection string
	Migrations string // migrations directory, e.g. ./migrations
	MigRcPath  string // override path to config file
	OutputJson bool   // stdout should be valid JSON
}

func GetConfig() (MigConfig, []string, *result.Response) {
	config := MigConfig{}

	flagConfig, subcommands, _ := GetConfigFromProcessFlags()

	config.OutputJson = flagConfig.OutputJson

	err := SetEnvFromConfigFile(flagConfig.MigRcPath) // reads .env and sets env vars but does not override

	if err != nil {
		return config, []string{}, nil
	}

	envConfig, _ := GetConfigFromEnvVars()

	if flagConfig.Connection != "" {
		config.Connection = flagConfig.Connection
	} else if envConfig.Connection != "" {
		config.Connection = envConfig.Connection
	} else {
		return config, subcommands, result.NewError("unable to determine server connection", "bad_config")
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
