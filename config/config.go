package config

// priority:
// CLI Flags
// Env Vars
// File Config

type MigConfig struct {
	connection string // DB connection string
	directory  string // migrations directory, e.g. ./migrations
}

func GetConfig() {
	Flags()
	Environment()
	File()
}
