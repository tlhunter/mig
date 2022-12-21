package commands

import (
	"fmt"
	"os"
)

// TODO: switch to flag.NewFlagSet

func Dispatch() {
	if len(os.Args) == 1 {
		fmt.Println("usage: mig <command>")
		os.Exit(9)
	}

	switch os.Args[1] {
	case "create":
		if len(os.Args) >= 3 {
			CommandCreate(os.Args[2])
			return
		}
		fmt.Println("usage: mig create '<migration name>'")
		os.Exit(10)
	default:
		fmt.Println("unsupported")
		os.Exit(10)
	}
}
