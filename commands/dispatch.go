package commands

import (
	"os"

	"github.com/fatih/color"
	"github.com/tlhunter/mig/config"
)

func Dispatch(cfg config.MigConfig, subcommands []string) {
	if len(subcommands) == 0 {
		color.White("usage: mig <command>")
		os.Exit(10)
	}

	switch subcommands[0] {
	case "create":
		if len(subcommands) >= 2 {
			CommandCreate(cfg, subcommands[1])
			return
		}
		color.White("usage: mig create \"<migration name>\"")
		os.Exit(10)

	case "init":
		CommandInit(cfg)

	case "lock":
		CommandLock(cfg)

	case "unlock":
		CommandUnlock(cfg)

	case "list":
		fallthrough
	case "ls":
		CommandList(cfg)

	case "status":
		CommandStatus(cfg)

	case "up":
		CommandUp(cfg)

	case "down":
		CommandDown(cfg)

	case "all":
		CommandAll(cfg)

	case "upto":
		if len(subcommands) >= 2 {
			CommandUpto(cfg, subcommands[1])
			return
		}
		color.White("usage: mig upto \"<migration name>\"")
		os.Exit(10)

	case "version":
		CommandVersion(cfg)

	default:
		color.White("unsupported command %s", subcommands[0])
		os.Exit(10)
	}
}
