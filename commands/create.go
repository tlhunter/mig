package commands

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/tlhunter/mig/config"
)

const TEMPLATE = `--BEGIN MIGRATION UP--
CREATE TABLE foo (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL
);
--END MIGRATION UP--
--BEGIN MIGRATION DOWN--
DROP TABLE foo;
--END MIGRATION DOWN--`

// TODO: Allow custom template file path in config

func CommandCreate(cfg config.MigConfig, rawName string) error {
	name := strings.ToLower(rawName)
	name = strings.Replace(name, " ", "_", -1)
	now := time.Now()

	// TODO: Lots of cleanup, basically anything not [a-z]

	filename := fmt.Sprintf("%04d%02d%02d%02d%02d%02d_%s.sql",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(),
		name)

	filePath := cfg.Migrations + "/" + filename

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.WriteString(TEMPLATE)
	if err != nil {
		return err
	}

	color.Green("created migration: " + filePath)

	return nil
}
