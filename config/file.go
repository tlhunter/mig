package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

const (
	MIGRC = ".migrc"
)

// Check the current working directory of the process for a .migrc file.
// If the file is found then read the contents and set it as environment variables.
// Existing environment variables are not overwritten in this manner.
// If a file isn't found in the current directory then check the parent directory.
// This repeats until reaching the root directory.
func SetEnvFromConfigFile(migRcPath string) error {
	if migRcPath != "" {
		return godotenv.Load(migRcPath)
	}

	dir, err := os.Getwd()

	if err != nil {
		// something bad is happening
		return err
	}

	for {
		candidate := filepath.Join(dir, MIGRC)

		if _, err := os.Stat(candidate); err == nil {
			err = godotenv.Load(candidate)

			if err != nil {
				// bad file or bad perms
				return err
			}

			return nil
		}

		if dir == "/" {
			// reached root w/ no file, giving up
			return nil
		}

		// visit parent
		dir = filepath.Dir(dir)
	}
}
