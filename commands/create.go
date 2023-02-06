package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/tlhunter/mig/config"
	"github.com/tlhunter/mig/result"
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

type CommandCreateResult struct {
	Filename string `json:"filename"`
}

func CommandCreate(cfg config.MigConfig, name string) result.Response {
	name = SanitizeName(name)
	now := time.Now()

	filename := fmt.Sprintf("%04d%02d%02d%02d%02d%02d_%s.sql",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(),
		name)

	filePath := filepath.Join(cfg.Migrations, filename)

	file, err := os.Create(filePath)
	if err != nil {
		return *result.NewErrorWithDetails("Unable to create migration file!", "unable_create_migration", err)
	}

	defer file.Close()

	_, err = file.WriteString(TEMPLATE)
	if err != nil {
		return *result.NewErrorWithDetails("Unable to write to migration file!", "unable_write_migration", err)
	}

	return *result.NewSerializable("created migration: "+filePath, CommandCreateResult{
		Filename: filePath,
	})
}

func SanitizeName(name string) string {
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "_")
	name = unsafes.ReplaceAllString(name, "")
	name = repeaters.ReplaceAllString(name, "_")

	return name
}
