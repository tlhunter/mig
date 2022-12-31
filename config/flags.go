package config

import (
	"os"

	"github.com/DavidGamba/go-getoptions"
)

func GetConfigFromProcessFlags() (MigConfig, []string, error) {
	opt := getoptions.New()
	connection := opt.String("connection", "")
	migrations := opt.String("migrations", "")
	migRcPath := opt.String("file", "")

	subcommand, err := opt.Parse(os.Args[1:])

	config := MigConfig{
		Connection: *connection,
		Migrations: *migrations,
		MigRcPath:  *migRcPath,
	}

	if err != nil {
		return config, subcommand, err
	}

	return config, subcommand, nil
}
