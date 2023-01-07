package commands

import (
	"errors"
	"fmt"

	"github.com/tlhunter/mig/config"
)

func Dispatch(cfg config.MigConfig, subcommands []string) error {
	var err error
	if len(subcommands) == 0 {
		err = errors.New("usage: mig <command>")
	}

	switch subcommands[0] {
	case "create":
		if len(subcommands) >= 2 {
			err = CommandCreate(cfg, subcommands[1])
		} else {
			err = errors.New("usage: mig create \"<migration name>\"")
		}

	case "init":
		err = CommandInit(cfg)

	case "lock":
		err = CommandLock(cfg)

	case "unlock":
		err = CommandUnlock(cfg)

	case "list":
		fallthrough
	case "ls":
		err = CommandList(cfg)

	case "status":
		err = CommandStatus(cfg)

	case "up":
		err = CommandUp(cfg)

	case "down":
		err = CommandDown(cfg)

	case "all":
		err = CommandAll(cfg)

	case "upto":
		if len(subcommands) >= 2 {
			err = CommandUpto(cfg, subcommands[1])
		} else {
			err = errors.New("usage: mig upto \"<migration name>\"")
		}

	case "version":
		err = CommandVersion(cfg)

	default:
		err = errors.New(fmt.Sprintf("unsupported command %s", subcommands[0]))
	}

	return err

}
