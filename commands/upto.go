package commands

import (
	"fmt"

	"github.com/fatih/color"

	"github.com/tlhunter/mig/config"
	"github.com/tlhunter/mig/database"
	"github.com/tlhunter/mig/migrations"
	"github.com/tlhunter/mig/result"
)

type CommandUptoResult struct {
	History []migrations.MigrationRowStatus `json:"history"`
}

// TODO: DRY up the common code between up, upto, all, and down
// TODO: This requires exact match "TIME_foo.sql"
//       would be nice to support "TIME_foo" or "foo" if unambiguous

func CommandUpto(cfg config.MigConfig, target string) result.Response {
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

	// ensure that the target is unapplied and ahead of current
	reachedCandidates := false
	if status.Last == nil || status.Last.Name == "" {
		// if no migrations were run then any migration is fair game
		reachedCandidates = true
	}
	targetIsCandidate := false

	for _, entry := range status.History {
		if !reachedCandidates {
			if entry.Migration.Name == status.Last.Name {
				reachedCandidates = true
			}
		}
		if reachedCandidates {
			if entry.Migration.Name == target {
				targetIsCandidate = true
				break
			}
		}
	}

	if !targetIsCandidate {
		// TODO: The name of the migration should be a JSON field
		return *result.NewError(fmt.Sprintf("Unable to find an unexecuted upcoming migration named %s", target), "cannot_find_migration")
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

		next := status.Next

		if next == "" {
			break
		}

		filename := cfg.Migrations + "/" + next

		queries, err := migrations.GetQueriesFromFile(filename)

		if err != nil {
			// TODO: Should tell user the `filename`
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
			res.AddErrorLn("You may want to manually add it and investigate the error.")
			res.AddErrorLn("Any remaining migrations will not be executed!")
			return *res
		}

		executedMigrations = append(executedMigrations, migration)

		if next == target {
			break
		}
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
