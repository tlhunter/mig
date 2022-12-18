package config

import "os"

const (
	CONNECTION = "MIG_CONNECTION"
	MIGRATIONS = "MIG_MIGRATIONS"
)

func Environment() {
	connection := os.Getenv(CONNECTION)

	println("Environment Variable:", connection)
}
