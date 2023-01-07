package commands

import (
	"fmt"
	"os"
	"regexp"
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

var unsafes = regexp.MustCompile(`[^a-z-_]`)
var repeaters = regexp.MustCompile(`_+`)

func CommandCreate(cfg config.MigConfig, name string) error {
	name = SanitizeName(name)
	now := time.Now()

	filename := fmt.Sprintf("%04d%02d%02d%02d%02d%02d_%s.sql",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(),
		name)

	filePath := cfg.Migrations + "/" + filename

	file, err := os.Create(filePath)
	if err != nil {
		color.Red("Unable to create migration file!")
		return err
	}

	defer file.Close()

	_, err = file.WriteString(TEMPLATE)
	if err != nil {
		color.Red("Unable to write to migration file!")
		return err
	}

	color.Green("created migration: " + filePath)

	return nil
}

func SanitizeName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "_")
	name = unsafes.ReplaceAllString(name, "")
	name = repeaters.ReplaceAllString(name, "_")

	return name
}
