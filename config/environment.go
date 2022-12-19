package config

import "os"

const (
	CONNECTION = "MIG_CONNECTION"
	MIGRATIONS = "MIG_MIGRATIONS"
)

func GetConfigFromEnvVars() (MigConfig, error) {
	connection := os.Getenv(CONNECTION)
	migrations := os.Getenv(MIGRATIONS)

	config := MigConfig{
		Connection: connection,
		Migrations: migrations,
	}

	return config, nil
}
