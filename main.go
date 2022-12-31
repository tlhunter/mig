package main

import (
	"os"

	"github.com/fatih/color"
	"github.com/tlhunter/mig/commands"
	"github.com/tlhunter/mig/config"
)

func main() {
	cfg, subcommands, err := config.GetConfig()

	if err != nil {
		color.Red("unable to parse configuration\n")
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}

	commands.Dispatch(cfg, subcommands)
}
