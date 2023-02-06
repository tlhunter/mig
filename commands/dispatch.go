package commands

import (
	"fmt"

	"github.com/tlhunter/mig/config"
	"github.com/tlhunter/mig/result"
)

func Dispatch(cfg config.MigConfig, subcommands []string) result.Response {
	var res result.Response
	if len(subcommands) == 0 {
		res.SetError("usage: mig <command>", "command_usage")
	}

	switch subcommands[0] {
	case "create":
		if len(subcommands) >= 2 {
			res = CommandCreate(cfg, subcommands[1])
		} else {
			res.SetError("usage: mig create \"<migration name>\"", "command_usage")
		}

	case "init":
		res = CommandInit(cfg)

	case "lock":
		res = CommandLock(cfg)

	case "unlock":
		res = CommandUnlock(cfg)

	case "list":
		fallthrough
	case "ls":
		res = CommandList(cfg)

	case "status":
		res = CommandStatus(cfg)

	case "up":
		res = CommandUp(cfg)

	case "down":
		res = CommandDown(cfg)

	case "all":
		res = CommandAll(cfg)

	case "upto":
		if len(subcommands) >= 2 {
			res = CommandUpto(cfg, subcommands[1])
		} else {
			res.SetError("usage: mig upto \"<migration name>\"", "command_usage")
		}

	case "version":
		res = CommandVersion(cfg)

	default:
		res.SetError(fmt.Sprintf("unsupported command %s", subcommands[0]), "command_unknown")
	}

	return res

}
