package commands

import (
	"fmt"

	"github.com/tlhunter/mig/config"
	"github.com/tlhunter/mig/database"
	"github.com/tlhunter/mig/migrations"
	"github.com/tlhunter/mig/result"
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

type MigrationName struct {
	Migration string `json:"migration"`
}

type CommandUpResult struct {
	MigrationBatch int                      `json:"batch"`
	Migration      *migrations.MigrationRow `json:"migration"`
}

func CommandUp(cfg config.MigConfig) result.Response {
	dbox, err := database.Connect(cfg.Connection)

	if err != nil {
		return *result.NewErrorWithDetails("database connection error", "db_conn", err)
	}

	defer dbox.Db.Close()

	status, err := migrations.GetStatus(cfg, dbox)

	if err != nil {
		return *result.NewErrorWithDetails("Encountered an error trying to get migrations status!", "retrieve_status", err)
	}

	if status.Skipped > 0 {
		return *result.NewError("Refusing to run with skipped migrations! Run `mig status` for details.", "abort_skipped_migrations")
	}

	next := status.Next

	if next == "" {
		return *result.NewError("There are no migrations to run.", "no_migrations")
	}

	filename := cfg.Migrations + "/" + next

	queries, err := migrations.GetQueriesFromFile(filename)

	if err != nil {
		return *result.NewErrorWithDetails("Error attempting to read next migration file!", "read_next_migration", err)
	}

	locked, err := database.ObtainLock(dbox)

	if err != nil {
		return *result.NewErrorWithDetails("Error obtaining lock for migration!", "obtain_lock", err)
	}

	if !locked {
		return *result.NewError("Unable to obtain lock for migration!", "obtain_lock")
	}

	var query string

	if queries.UpTx {
		query = BEGIN.For(dbox.Type) + queries.Up + END.For(dbox.Type)
	} else {
		query = queries.Up
	}

	_, err = dbox.Db.Exec(query)

	if err != nil {
		return *result.NewErrorWithDetails("Encountered an error while running migration!", "migration_failed", err)
	}

	migration, err := migrations.AddMigration(dbox, next)

	if err != nil {
		res := result.NewErrorWithDetails("The migration query executed but unable to track it in the migrations table!", "untracked_migration", err)
		res.AddErrorLn("You may want to manually add it and investigate the error.")
		return *res
	}

	res := result.NewSerializable(fmt.Sprintf("Migration %s was successfully applied!", next), CommandUpResult{
		MigrationBatch: migration.Batch,
		Migration:      &migration,
	})

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
