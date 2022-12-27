package migrations

import (
	"fmt"
	"os"
	"strings"
)

func ListFiles(directory string) ([]string, error) {
	var migFiles []string

	files, err := os.ReadDir(directory)

	if err != nil {
		fmt.Println("unable to read migrations direcotry")
		return migFiles, err
	}

	for _, entry := range files {
		name := entry.Name()

		if entry.IsDir() || strings.HasPrefix(name, ".") || !strings.HasSuffix(name, ".sql") {
			continue
		}

		migFiles = append(migFiles, entry.Name())
	}

	return migFiles, nil
}
