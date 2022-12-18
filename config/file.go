package config

// look for a file named .migrc in the current directory up to the root

import (
	"fmt"
	"io/ioutil"
	"os"
)

const (
	MIGRC = ".migrc"
)

func File() {
	cwd, err := os.Getwd()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Current Directory:", cwd)

	contents, err := ioutil.ReadFile(cwd + "/" + MIGRC)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(".migrc:", string(contents))
}
