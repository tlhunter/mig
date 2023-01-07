package commands

import (
	"errors"

	"github.com/fatih/color"
	"github.com/tlhunter/mig/config"
	"github.com/tlhunter/mig/database"
	"github.com/tlhunter/mig/migrations"
)

var (
	BEGIN = database.QueryBox{
		Postgres: "BEGIN TRANSACTION;\n",
		Mysql:    "START TRANSACTION;\n",
	}
	END = database.QueryBox{
		Postgres: "COMMIT TRANSACTION;\n",
		Mysql:    "COMMIT;\n",
	}
)

func CommandUp(cfg config.MigConfig) error {
	dbox, err := database.Connect(cfg.Connection)

	if err != nil {
		return err
	}

	defer dbox.Db.Close()

	status, err := migrations.GetStatus(cfg, dbox)

	if err != nil {
		color.Red("Encountered an error trying to get migrations status!")
		return err
	}

	if status.Skipped > 0 {
		return errors.New("Refusing to run with skipped migrations! Run `mig status` for details.")
	}

	next := status.Next

	if next == "" {
		return errors.New("There are no migrations to run.")
	}

	filename := cfg.Migrations + "/" + next

	queries, err := migrations.GetQueriesFromFile(filename)

	if err != nil {
		color.Red("Error attempting to read next migration file!")
		return err
	}

	locked, err := database.ObtainLock(dbox)

	if err != nil {
		color.Red("Error obtaining lock for migration!")
		return err
	}

	if !locked {
		return errors.New("Unable to obtain lock for migration!")
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

	err = migrations.AddMigration(dbox, next)

	if err != nil {
		color.Red("The migration query executed but unable to track it in the migrations table!")
		color.White("You may want to manually add it and investigate the error.")
		return err
	}

	released, err := database.ReleaseLock(dbox)

	if err != nil {
		color.Red("Error obtaining lock for migration!")
		return err
	}

	if !released {
		return errors.New("Unable to obtain lock for migration!")
	}

	return nil
}
