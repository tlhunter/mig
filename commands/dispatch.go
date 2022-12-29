package commands

import (
	"os"

	"github.com/fatih/color"
	"github.com/tlhunter/mig/config"
)

// TODO: switch to flag.NewFlagSet

func Dispatch(cfg config.MigConfig) {
	if len(os.Args) == 1 {
		color.White("usage: mig <command>")
		os.Exit(9)
	}

	switch os.Args[1] {
	case "create":
		if len(os.Args) >= 3 {
			CommandCreate(cfg, os.Args[2])
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

	default:
		color.White("unsupported command %s", os.Args[1])
		os.Exit(10)
	}
}
