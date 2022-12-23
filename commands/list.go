package commands

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/tlhunter/mig/config"
	"github.com/tlhunter/mig/database"
)

const (
	LIST = `SELECT * FROM migrations ORDER BY id ASC;`
)

func CommandList(cfg config.MigConfig) error {
	db := database.Connect(cfg.Connection)
	defer db.Close()

	migDir := cfg.Migrations

	migrations, err := os.ReadDir(migDir)

	if err != nil {
		fmt.Println("unable to read migrations direcotry")
		return err
	}

	// TODO: iterate through files and rows together, looking for missing values
	// applied migrations are green, skipped are red, not yet applied migrations are yellow
	// and if that's too hard, applied=green and unapplied = red
	// finding a skipped migration is a big deal, need to recommend a rollback

	for _, f := range migrations {
		filename := f.Name()

		if f.IsDir() || strings.HasPrefix(filename, ".") || !strings.HasSuffix(filename, ".sql") {
			continue
		}

		fmt.Println(f.Name())
	}

	rows, err := db.Query(LIST)

	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var batch int
		var migration_time time.Time
		err = rows.Scan(&id, &name, &batch, &migration_time)

		if err != nil {
			return err
		}

		fmt.Println(name, "\t", id, "\t", batch, "\t", migration_time)
	}

	return nil
}
