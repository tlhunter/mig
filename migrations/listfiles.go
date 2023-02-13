package migrations

import (
	"os"
	"strings"
)

func ListFiles(directory string) ([]string, error) {
	var migFiles []string

	files, err := os.ReadDir(directory)
	if err != nil {
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
