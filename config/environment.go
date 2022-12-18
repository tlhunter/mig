package config

import "os"

const (
	CONNECTION = "MIG_CONNECTION"
	MIGRATIONS = "MIG_MIGRATIONS"
)

func Environment() (MigConfig, error) {
	connection := os.Getenv(CONNECTION)
	migrations := os.Getenv(MIGRATIONS)

	println("conn env", connection)
	println("mig env", migrations)

	config := MigConfig{
		Connection: connection,
		Migrations: migrations,
	}

	return config, nil
}
