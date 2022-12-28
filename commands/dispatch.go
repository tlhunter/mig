package commands

import (
	"fmt"
	"os"

	"github.com/tlhunter/mig/config"
)

// TODO: switch to flag.NewFlagSet

func Dispatch(cfg config.MigConfig) {
	if len(os.Args) == 1 {
		fmt.Println("usage: mig <command>")
		os.Exit(9)
	}

	switch os.Args[1] {
	case "create":
		if len(os.Args) >= 3 {
			CommandCreate(cfg, os.Args[2])
			return
		}
		fmt.Println("usage: mig create '<migration name>'")
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
		CommandUp(cfg)

	default:
		fmt.Println("unsupported")
		os.Exit(10)
	}
}
