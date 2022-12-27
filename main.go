package main

import (
	"os"

	"github.com/tlhunter/mig/commands"
	"github.com/tlhunter/mig/config"
)

func main() {
	cfg, err := config.GetConfig()

	if err != nil {
		os.Stderr.WriteString("unable to parse configuration\n")
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}

	commands.Dispatch(cfg)
}
