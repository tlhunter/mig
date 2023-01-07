package commands

import (
	"errors"

	"github.com/fatih/color"
	"github.com/tlhunter/mig/config"
	"github.com/tlhunter/mig/database"
	"github.com/tlhunter/mig/migrations"
)

func CommandDown(cfg config.MigConfig) error {
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

	last := status.Last

	filename := cfg.Migrations + "/" + last.Name

	queries, err := migrations.GetQueriesFromFile(filename)

	if err != nil {
		color.Red("Error attempting to read last migration file!")
		color.White("Normally a missing migration file isn't a big deal but it's a no go for migrating down.")
		color.White("This file is required before continuing. Perhaps it can be pulled from version control?")
		return err
	}

	locked, err := database.ObtainLock(dbox)

	if err != nil {
		color.Red("Error obtaining lock for migrating down!")
		return err
	}

	if !locked {
		return errors.New("Unable to obtain lock for migrating down!")
	}

	var query string

	if queries.DownTx {
		query = BEGIN.For(dbox.Type) + queries.Down + END.For(dbox.Type)
	} else {
		query = queries.Down
	}

	_, err = dbox.Db.Exec(query)

	if err != nil {
		color.Red("Encountered an error while running down migration!")
		return err
	}

	color.Green("Migration down %s was successfully applied!", last.Name)

	err = migrations.RemoveMigration(dbox, last.Name, last.Id)

	if err != nil {
		color.Red("The migration down query executed but unable to track it in the migrations table!")
		color.White("You may want to manually remove it and investigate the error.")
		return err
	}

	released, err := database.ReleaseLock(dbox)

	if err != nil {
		color.Red("Error obtaining lock for down migration!")
		return err
	}

	if !released {
		return errors.New("Unable to obtain lock for down migration!")
	}

	return nil
}
