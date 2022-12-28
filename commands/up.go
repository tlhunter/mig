package commands

import (
	"fmt"
	"os"

	"github.com/tlhunter/mig/config"
	"github.com/tlhunter/mig/database"
	"github.com/tlhunter/mig/migrations"
)

const (
	BEGIN = "BEGIN TRANSACTION;\n"
	END   = "COMMIT TRANSACTION;\n"
)

func CommandUp(cfg config.MigConfig) error {
	db := database.Connect(cfg.Connection)

	defer db.Close()

	status, err := migrations.GetStatus(cfg, db, false)

	if err != nil {
		os.Stderr.WriteString("Encountered an error trying to get migrations status!\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return err
	}

	if status.Skipped > 0 {
		os.Stderr.WriteString("Refusing to run with skipped migrations! Run `mig status` for details.\n")
		return nil
	}

	next := status.Next

	filename := cfg.Migrations + "/" + next

	queries, err := migrations.GetQueriesFromFile(filename)

	if err != nil {
		os.Stderr.WriteString("Error attempting to read next migration file!\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return err
	}

	locked, err := database.ObtainLock(db)

	if err != nil {
		os.Stderr.WriteString("Error obtaining lock for migration!\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return err
	}

	if !locked {
		os.Stderr.WriteString("Unable to obtain lock for migration!\n")
		return nil
	}

	_, err = db.Exec(BEGIN + queries.Up + END)

	if err != nil {
		os.Stderr.WriteString("Encountered an error while running migration!\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return err
	}

	fmt.Printf("Migration %s was successfully applied!\n", next)

	err = migrations.RecordMigration(db, next)

	if err != nil {
		os.Stderr.WriteString("The migration query executed but unable to track it in the migrations table!\n")
		os.Stderr.WriteString("You may want to manually add it and investigate the error.\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return err
	}

	released, err := database.ReleaseLock(db)

	if err != nil {
		os.Stderr.WriteString("Error obtaining lock for migration!\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return err
	}

	if !released {
		os.Stderr.WriteString("Unable to obtain lock for migration!\n")
		return nil
	}

	return nil
}
