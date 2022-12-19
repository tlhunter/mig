package config

import (
	"flag"
)

func GetConfigFromProcessFlags() (MigConfig, error) {
	connection := flag.String("connection", "", "SQL connection string")
	migrations := flag.String("migrations", "", "Migrations directory")

	flag.Parse()

	config := MigConfig{
		Connection: *connection,
		Migrations: *migrations,
	}

	return config, nil
}
