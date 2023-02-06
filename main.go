package main

import (
	"os"

	"github.com/tlhunter/mig/commands"
	"github.com/tlhunter/mig/config"
	"github.com/tlhunter/mig/result"
)

func main() {
	cfg, subcommands, err := config.GetConfig()

	var res result.Response

	if err != nil && len(subcommands) == 1 && subcommands[0] == "version" {
		res = commands.CommandVersion(cfg)
	} else if err == nil {
		res = commands.Dispatch(cfg, subcommands)
	}

	res.Display(cfg.OutputJson)
	os.Exit(int(res.ExitStatus))
}
