package main

import (
	"os"

	"github.com/fatih/color"
	"github.com/tlhunter/mig/commands"
	"github.com/tlhunter/mig/config"
)

func main() {
	cfg, subcommands, err := config.GetConfig()

	if err != nil && len(subcommands) == 1 && subcommands[0] == "version" {
		err = commands.CommandVersion()
	} else if err != nil {
		color.Red("unable to parse configuration")
	} else {
		err = commands.Dispatch(cfg, subcommands)
	}

	if err != nil {
		os.Stderr.WriteString(color.RedString(err.Error()) + "\n")
		os.Exit(1)
	}
}
