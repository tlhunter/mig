package config

import (
	"flag"
	"fmt"
)

func Flags() (MigConfig, error) {
	connection := flag.String("connection", "", "SQL connection string")
	migrations := flag.String("migrations", "", "Migrations directory")

	flag.Parse()
	// TODO: Fails when env var is missing. Will need to also check env vars and config files.
	// TODO: Error message lists -connection, prefer the more common --connection text.

	fmt.Println("conn flag", connection)
	fmt.Println("mig flag", migrations)

	config := MigConfig{
		connection: *connection,
		migrations: *migrations,
	}

	return config, nil
}
