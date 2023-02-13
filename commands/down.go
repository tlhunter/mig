package commands

import (
	"fmt"

	"github.com/tlhunter/mig/config"
	"github.com/tlhunter/mig/database"
	"github.com/tlhunter/mig/migrations"
	"github.com/tlhunter/mig/result"
)

func CommandDown(cfg config.MigConfig) result.Response {
	dbox, err := database.Connect(cfg.Connection)

	if err != nil {
		return *result.NewErrorWithDetails("database connection error", "db_conn", err)
	}

	defer dbox.Db.Close()

	status, err := migrations.GetStatus(cfg, dbox)

	if err != nil {
		return *result.NewErrorWithDetails("Encountered an error trying to get migrations status!", "retrieve_status", err)
	}

	last := status.Last

	if last == nil {
		return *result.NewError("There are no migrations to revert.", "nothing_to_revert")
	}

	filename := cfg.Migrations + "/" + last.Name

	queries, err := migrations.GetQueriesFromFile(filename)

	if err != nil {
		res := *result.NewErrorWithDetails("Error attempting to read last migration file!", "unable_read_migration_file", err)

		res.AddErrorLn("Normally a missing migration file isn't a big deal but it's a no go for migrating down.")
		res.AddErrorLn("This file is required before continuing. Perhaps it can be pulled from version control?")

		return res
	}

	locked, err := database.ObtainLock(dbox)

	if err != nil {
		return *result.NewErrorWithDetails("Error obtaining lock for migration down!", "obtain_lock", err)
	}

	if !locked {
		return *result.NewError("Unable to obtain lock for migrating down!", "obtain_lock")
	}

	err = dbox.ExecMaybeTx(queries.Down, queries.DownTx)

	if err != nil {
		return *result.NewErrorWithDetails("Encountered an error while running down migration!", "migration_failed", err)
	}

	res := result.NewSuccess(fmt.Sprintf("Down migration for %s was successfully applied!", last.Name))

	err = migrations.RemoveMigration(dbox, last.Name, last.Id)

	if err != nil {
		res.SetError("The migration down query executed but unable to track it in the migrations table!", "untracked_migration")
		res.SetErrorDetails(err)
		res.AddErrorLn("You may want to manually remove it and investigate the error.")
		return *res
	}

	released, err := database.ReleaseLock(dbox)

	if err != nil {
		res.SetError("Error obtaining lock for down migration!", "release_lock")
		return *res
	}

	if !released {
		res.SetError("Unable to obtain lock for down migration!", "release_lock")
	}

	return *res
}
