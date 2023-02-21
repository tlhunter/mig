package commands

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/tlhunter/mig/config"
	"github.com/tlhunter/mig/database"
	"github.com/tlhunter/mig/migrations"
	"github.com/tlhunter/mig/result"
)

var EXIST_MIGRATIONS = database.QueryBox{
	Postgres: `SELECT EXISTS (
SELECT FROM
	pg_tables
WHERE
	schemaname = 'public' AND
	tablename  = 'migrations'
) AS table_exists;`,
	Mysql:  `CALL sys.table_exists(DATABASE(), 'migrations', @table_type); SELECT @table_type LIKE 'BASE TABLE';`,
	Sqlite: `SELECT COUNT(name) >= 1 AS table_is_present FROM sqlite_master WHERE type='table' AND name='migrations';`,
}

var EXIST_LOCK = database.QueryBox{
	Postgres: `SELECT EXISTS (
	SELECT FROM
		pg_tables
	WHERE
		schemaname = 'public' AND
		tablename  = 'migrations_lock'
	) AS table_exists;`,
	Mysql:  `CALL sys.table_exists(DATABASE(), 'migrations_lock', @table_type); SELECT @table_type LIKE 'BASE TABLE';`,
	Sqlite: `SELECT COUNT(name) >= 1 AS table_is_present FROM sqlite_master WHERE type='table' AND name='migrations_lock';`,
}

var DESCRIBE = database.QueryBox{
	Postgres: `SELECT
	table_name,
	column_name,
	data_type
FROM
	information_schema.columns
WHERE
	table_name = 'migrations' OR table_name = 'migrations_lock'
ORDER BY
	table_name, column_name;`,
	Mysql:  `DESC migrations;`,                 // unused
	Sqlite: `pragma table_info('migrations');`, // unused
}

var LOCK_STATUS = database.QueryBox{
	Postgres: `SELECT is_locked FROM migrations_lock WHERE index = 1;`,
	Mysql:    `SELECT is_locked FROM migrations_lock WHERE ` + "`index`" + ` = 1;`,
	Sqlite:   `SELECT is_locked FROM migrations_lock WHERE "index" = 1;`,
}

type StatusResponse struct {
	Locked bool `json:"locked"`
	Status any  `json:"status"`
}

// Provide a narrative to the user about the current status of mig
// Inspired by `git status` and `brew doctor`
func CommandStatus(cfg config.MigConfig) result.Response {

	// Attempt to connect to database

	dbox, err := database.Connect(cfg.Connection)
	if err != nil {
		return *result.NewErrorWithDetails("database connection error", "db_conn", err)
	}

	// Check if migration tables exist

	existMigrations := false

	err = dbox.QueryRow(EXIST_MIGRATIONS).Scan(&existMigrations)
	if err != nil {
		return *result.NewErrorWithDetails("unable to tell if 'migrations' table exists!", "unable_check_migrations", err)
	}

	existLock := false

	err = dbox.QueryRow(EXIST_LOCK).Scan(&existLock)
	if err != nil {
		return *result.NewErrorWithDetails("unable to tell if 'migrations_lock' table exists!", "unable_check_migrations_lock", err)
	}

	if !existMigrations && !existLock {
		res := result.NewError("The tables used for tracking migrations are missing.", "missing_tables")
		res.ExitStatus = 9

		res.AddErrorLn(color.WhiteString("This likely means that mig hasn't yet been initialized."))
		res.AddErrorLn(color.WhiteString("This can be solved by running the following command:"))
		res.AddErrorLn(color.WhiteString("$ mig init"))

		return *res
	}

	if !existMigrations {
		res := *result.NewError("The 'migrations' table is missing but the 'migrations_lock' table is present!", "missing_migrations_table")
		res.ExitStatus = 9
		res.AddErrorLn("This might mean that data has been corrupted and that migration status is missing.")
		res.AddErrorLn("Consider looking into the root cause of the problem.")
		res.AddErrorLn("The quickest fix is to delete the lock table and initialize again:")
		res.AddErrorLn("> DROP TABLE migrations_lock;")
		res.AddErrorLn("$ mig init")
		return res
	}

	if !existLock {
		res := *result.NewError("The 'migrations' table is present but the 'migrations_lock' table is missing!", "missing_lock_table")
		res.ExitStatus = 9
		res.AddErrorLn("This might mean that data has been corrupted.")
		res.AddErrorLn("Consider looking into the cause of the problem.")
		res.AddErrorLn("The quickest fix is backup the migrations table data, initialize again, then restore the data:")
		res.AddErrorLn("> ALTER TABLE migrations RENAME TO migrations_backup;")
		res.AddErrorLn("$ mig init")
		res.AddErrorLn("> DROP TABLE migrations;")
		res.AddErrorLn("> ALTER TABLE migrations_backup RENAME TO migrations;")
		return res
	}

	res := result.NewSerializable("", "")

	if dbox.IsMysql {
		// TODO: rebuild the gnarly checks
		res.AddSuccessLn(color.YellowString("migration table description check is currently unimplemented for mysql."))
	} else if dbox.IsSqlite {
		// TODO: rebuild the gnarly checks
		res.AddSuccessLn(color.YellowString("migration table description check is currently unimplemented for sqlite."))
	} else if dbox.IsPostgres {
		rows, err := dbox.Query(DESCRIBE)
		if err != nil {
			return *result.NewErrorWithDetails("unable to describe the migration tables!", "unable_describe", err)
		}

		// Check if tables have correct columns

		var table, column, data string

		// struggling to find the Go-way to do this...
		rows.Next()
		rows.Scan(&table, &column, &data)
		if table != "migrations" || column != "batch" || data != "integer" {
			return *result.NewError("expected migrations.batch of type integer", "invalid_batch_type")
		}

		rows.Next()
		rows.Scan(&table, &column, &data)
		if table != "migrations" || column != "id" || data != "integer" {
			return *result.NewError("expected migrations.id of type integer", "invalid_id_type")
		}

		rows.Next()
		rows.Scan(&table, &column, &data)
		if table != "migrations" || column != "migration_time" || data != "timestamp with time zone" {
			return *result.NewError("expected migrations.migration_time of type timestamp with time zone", "invalid_time_type")
		}

		rows.Next()
		rows.Scan(&table, &column, &data)
		if table != "migrations" || column != "name" || data != "character varying" {
			return *result.NewError("expected migrations.name of type character varying", "invalid_name_type")
		}

		rows.Next()
		rows.Scan(&table, &column, &data)
		if table != "migrations_lock" || column != "index" || data != "integer" {
			return *result.NewError("expected migrations_lock.index of type integer", "invalid_index_type")
		}

		rows.Next()
		rows.Scan(&table, &column, &data)
		if table != "migrations_lock" || column != "is_locked" || data != "integer" {
			return *result.NewError("expected migrations_lock.is_locked of type integer", "invalid_locked_type")
		}
	} else {
		panic("unknown database: " + dbox.Type)
	}

	// Check if locked

	locked := false

	err = dbox.QueryRow(LOCK_STATUS).Scan(&locked)
	if err != nil {
		return *result.NewErrorWithDetails("unable to determine lock status!", "unable_determine_lock_status", err)
	}

	if locked {
		res.AddSuccessLn(color.RedString("Migrations are currently locked!"))
		res.AddSuccessLn(color.WhiteString("It could be that a migration is in progress. However it could also mean that a migration failed."))
		res.AddSuccessLn(color.WhiteString("If migrations remain locked then someone will want to investigate the failed migration."))
		res.AddSuccessLn(color.WhiteString("Once that's over you can unlock migrations by running the following:"))
		res.AddSuccessLn(color.WhiteString("$ mig unlock"))
		res.AddSuccessLn("")
		// Note: Don't need to return at this point
	}

	// Display the name of the last run migration

	status, err := migrations.GetStatus(cfg, dbox)
	if err != nil {
		return *result.NewErrorWithDetails("Encountered an error trying to get migrations status!", "retrieve_status", err)
	}

	if cfg.OutputJson {
		status.History = nil // omit for status command, it's still present for list command

		res.Serializable = StatusResponse{
			Status: status,
			Locked: locked,
		}

		if status.Skipped > 0 {
			res.ExitStatus = 10
		}

		return *res
	}

	if status.Last != nil && status.Last.Name != "" {
		res.AddSuccessLn(color.HiWhiteString("Last Migration: %s (id=%d,batch=%d) on %s", status.Last.Name, status.Last.Id, status.Last.Batch, status.Last.Time.Format(time.RFC3339)))
		res.AddSuccessLn("")
	} else {
		res.AddSuccessLn("No migrations have yet been executed.")
	}

	res.AddSuccessLn(fmt.Sprintf("Applied: %d, Unapplied: %d, Skipped: %d, Missing: %d", status.Applied, status.Unapplied, status.Skipped, status.Missing))

	// TODO: How to return custon JSON format but also allow later failure?

	if status.Skipped > 0 {
		res := *result.NewError("There are at least one skipped migrations! Mig will not be able to run migrations until this is fixed.", "encounter_skipped_migrations")
		res.AddErrorLn("A skipped migration happens when a local migration file is older than the most recently run migration.")
		res.AddErrorLn("To fix this, rename any skipped migrations so that their timestamps are newer.")
		res.AddErrorLn("Run this command to list skipped migrations:")
		res.AddErrorLn("$ mig list")
	}

	if status.Next != "" {
		res.AddSuccessLn(color.HiWhiteString("Next Migration: %s", status.Next))
		res.AddSuccessLn(color.WhiteString("To run this migration, execute the following command:"))
		res.AddSuccessLn(color.WhiteString("$ mig up"))
	}

	return *res
}
