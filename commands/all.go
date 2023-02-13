package commands

import (
	"github.com/fatih/color"

	"github.com/tlhunter/mig/config"
	"github.com/tlhunter/mig/database"
	"github.com/tlhunter/mig/migrations"
	"github.com/tlhunter/mig/result"
)

type CommandUpFamilyResult struct {
	MigrationBatch int                        `json:"batch"`
	Migrations     *[]migrations.MigrationRow `json:"migrations"`
}

func CommandAll(cfg config.MigConfig) result.Response {
	dbox, err := database.Connect(cfg.Connection)
	if err != nil {
		return *result.NewErrorWithDetails("database connection error", "db_conn", err)
	}

	defer dbox.Db.Close()

	// First call to GetStatus is mostly unused. if it fails then don't continue.
	status, err := migrations.GetStatus(cfg, dbox)
	if err != nil {
		return *result.NewErrorWithDetails("Encountered an error trying to get migrations status!", "retrieve_status", err)
	}
	if status.Skipped > 0 {
		return *result.NewError("Refusing to run with skipped migrations! Run `mig status` for details.", "abort_skipped_migrations")
	}
	if status.Next == "" {
		return *result.NewError("There are no migrations to run.", "no_migrations")
	}

	locked, err := database.ObtainLock(dbox)
	if err != nil {
		return *result.NewErrorWithDetails("Error obtaining lock for migration!", "obtain_lock", err)
	}
	if !locked {
		return *result.NewError("Unable to obtain lock for migration!", "obtain_lock")
	}

	highest, err := migrations.GetHighestValues(dbox)
	batchId := highest.Batch

	var executedMigrations []migrations.MigrationRow

	res := result.NewSerializable(color.HiWhiteString("Running migrations for batch %d...", batchId), CommandUpFamilyResult{
		MigrationBatch: batchId,
		Migrations:     &executedMigrations,
	})

	for {
		status, err := migrations.GetStatus(cfg, dbox)
		if err != nil {
			return *result.NewErrorWithDetails("Encountered an error trying to get migrations status!", "retrieve_status", err)
		}

		next := status.Next
		if next == "" {
			break
		}

		filename := cfg.Migrations + "/" + next

		queries, err := migrations.GetQueriesFromFile(filename)
		if err != nil {
			return *result.NewErrorWithDetails("Error attempting to read next migration file!", "read_next_migration", err)
		}

		err = dbox.ExecMaybeTx(queries.Up, queries.UpTx)
		if err != nil {
			return *result.NewErrorWithDetails("Encountered an error while running migration!", "migration_failed", err)
		}

		res.AddSuccessLn(color.GreenString("Migration %s was successfully applied!", next))

		migration, err := migrations.AddMigrationWithBatch(dbox, next, batchId)
		if err != nil {
			res.SetError("The migration query executed but unable to track it in the migrations table!", "untracked_migration")
			res.SetErrorDetails(err)
			res.AddErrorLn("You may want to manually add it and investigate the error.")
			res.AddErrorLn("Any remaining migrations will not be executed!")
			return *res
		}

		executedMigrations = append(executedMigrations, migration)
	}

	released, err := database.ReleaseLock(dbox)
	if err != nil {
		res.SetError("Error releasing lock after running migration!", "release_lock")
		return *res
	}
	if !released {
		res.SetError("Unable to release lock after running migration!", "release_lock")
	}

	return *res
}
