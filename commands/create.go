package commands

import (
	"fmt"
	"os"
	"strings"
	"time"
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

func CommandCreate(rawName string) error {
	name := strings.ToLower(rawName)
	name = strings.Replace(name, " ", "_", -1)
	now := time.Now()

	// TODO: Lots of cleanup, basically anything not [a-z]

	filename := fmt.Sprintf("%04d%02d%02d%02d%02d%02d-%s.sql",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(),
		name)

	file, err := os.Create("migrations/" + filename)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.WriteString(TEMPLATE)
	if err != nil {
		return err
	}

	return nil
	// TODO: Write file
}
