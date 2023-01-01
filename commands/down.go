package commands

import (
	"os"

	"github.com/fatih/color"
	"github.com/tlhunter/mig/config"
	"github.com/tlhunter/mig/database"
	"github.com/tlhunter/mig/migrations"
)

func CommandDown(cfg config.MigConfig) error {
	db, dbType := database.Connect(cfg.Connection)

	defer db.Close()

	status, err := migrations.GetStatus(cfg, db, false)

	if err != nil {
		color.Red("Encountered an error trying to get migrations status!\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return err
	}

	last := status.Last

	filename := cfg.Migrations + "/" + last.Name

	queries, err := migrations.GetQueriesFromFile(filename)

	if err != nil {
		color.Red("Error attempting to read last migration file!\n")
		color.White("Normally a missing migration file isn't a big deal but it's a no go for migrating down.\n")
		color.White("This file is required before continuing. Perhaps it can be pulled from version control?\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return err
	}

	locked, err := database.ObtainLock(db)

	if err != nil {
		color.Red("Error obtaining lock for migrating down!\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return err
	}

	if !locked {
		color.Red("Unable to obtain lock for migrating down!\n")
		return nil
	}

	var query string

	if queries.DownTx {
		query = BEGIN.For(dbType) + queries.Down + END.For(dbType)
	} else {
		query = queries.Down
	}

	_, err = db.Exec(query)

	if err != nil {
		color.Red("Encountered an error while running down migration!\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return err
	}

	color.Green("Migration down %s was successfully applied!\n", last.Name)

	err = migrations.RemoveMigration(db, last.Name, last.Id)

	if err != nil {
		color.Red("The migration down query executed but unable to track it in the migrations table!\n")
		color.White("You may want to manually remove it and investigate the error.\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return err
	}

	released, err := database.ReleaseLock(db)

	if err != nil {
		color.Red("Error obtaining lock for down migration!\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return err
	}

	if !released {
		color.Red("Unable to obtain lock for down migration!\n")
		return nil
	}

	return nil
}
