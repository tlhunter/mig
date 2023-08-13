package main

import (
	"os"

	"github.com/tlhunter/mig/commands"
	"github.com/tlhunter/mig/config"
	"github.com/tlhunter/mig/result"
)

func main() {
	cfg, subcommands, bail := config.GetConfig()

	var res result.Response

	if bail != nil && len(subcommands) == 0 {
		res.SetError("usage: mig <command>", "command_usage")
	} else if bail != nil && len(subcommands) == 1 && subcommands[0] == "version" {
		res = commands.CommandVersion(cfg)
	} else if bail != nil {
		res = *bail
	} else if bail == nil {
		res = commands.Dispatch(cfg, subcommands)
	}

	res.Display(cfg.OutputJson)
	os.Exit(int(res.ExitStatus))
}
