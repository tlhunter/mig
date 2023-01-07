package commands

import (
	"errors"

	"github.com/fatih/color"

	"github.com/tlhunter/mig/config"
	"github.com/tlhunter/mig/database"
	"github.com/tlhunter/mig/migrations"
)

func CommandAll(cfg config.MigConfig) error {
	dbox, err := database.Connect(cfg.Connection)

	if err != nil {
		return err
	}

	defer dbox.Db.Close()

	// First call to GetStatus is mostly unused. if it fails then don't continue.
	status, err := migrations.GetStatus(cfg, dbox)

	if err != nil {
		return errors.New("Encountered an error trying to get migrations status!")
	}

	if status.Skipped > 0 {
		return errors.New("Refusing to run with skipped migrations! Run `mig status` for details.")
	}

	if status.Next == "" {
		return errors.New("There are no migrations to run.")
	}

	locked, err := database.ObtainLock(dbox)

	if err != nil {
		color.Red("Error obtaining lock for running migrations!")
		return err
	}

	if !locked {
		return errors.New("Unable to obtain lock for running migrations!")
	}

	highest, err := migrations.GetHighestValues(dbox)
	batchId := highest.Batch

	color.HiWhite("Running migrations for batch %d...", batchId)

	for {
		status, err := migrations.GetStatus(cfg, dbox)

		next := status.Next

		if next == "" {
			break
		}

		filename := cfg.Migrations + "/" + next

		queries, err := migrations.GetQueriesFromFile(filename)

		if err != nil {
			color.Red("Error attempting to read next migration file!")
			return err
		}

		var query string

		if queries.UpTx {
			query = BEGIN.For(dbox.Type) + queries.Up + END.For(dbox.Type)
		} else {
			query = queries.Up
		}

		_, err = dbox.Db.Exec(query)

		if err != nil {
			color.Red("Encountered an error while running migration!")
			return err
		}

		color.Green("Migration %s was successfully applied!", next)

		err = migrations.AddMigrationWithBatch(dbox, next, batchId)

		if err != nil {
			color.Red("The migration query executed but unable to track it in the migrations table!")
			color.White("You may want to manually add it and investigate the error.")
			color.White("Any remaining migrations will not be executed!")
			return err
		}
	}

	released, err := database.ReleaseLock(dbox)

	if err != nil {
		color.Red("Error releasing lock for migration!")
		return err
	}

	if !released {
		return errors.New("Unable to release lock for migration!")
	}

	return nil
}
