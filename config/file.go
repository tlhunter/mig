package config

// look for a file named .migrc in the current directory up to the root

import (
	"fmt"
	"io/ioutil"
	"os"
)

const (
	MIGRC = ".migrc" // TODO: allow override via --file
)

func File() (MigConfig, error) {
	config := MigConfig{}

	cwd, err := os.Getwd()

	if err != nil {
		return config, err
	}

	//var path = cwd

	// TODO

	fmt.Println("Current Directory:", cwd)

	contents, err := ioutil.ReadFile(cwd + "/" + MIGRC)

	if err != nil {
		fmt.Println(err)
	}

	// path, _ := filepath.Split(path)

	fmt.Println(".migrc:", string(contents))

	return config, nil
}
