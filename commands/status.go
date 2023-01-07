package commands

import (
	"errors"
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/tlhunter/mig/config"
	"github.com/tlhunter/mig/database"
	"github.com/tlhunter/mig/migrations"
)

var EXIST_MIGRATIONS = database.QueryBox{
	Postgres: `SELECT EXISTS (
SELECT FROM
	pg_tables
WHERE
	schemaname = 'public' AND
	tablename  = 'migrations'
) AS table_exists;`,
	Mysql: `CALL sys.table_exists(DATABASE(), 'migrations', @table_type); SELECT @table_type LIKE 'BASE TABLE';`,
}

var EXIST_LOCK = database.QueryBox{
	Postgres: `SELECT EXISTS (
	SELECT FROM
		pg_tables
	WHERE
		schemaname = 'public' AND
		tablename  = 'migrations_lock'
	) AS table_exists;`,
	Mysql: `CALL sys.table_exists(DATABASE(), 'migrations_lock', @table_type); SELECT @table_type LIKE 'BASE TABLE';`,
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
	Mysql: `DESC migrations;`,
}

var LOCK_STATUS = database.QueryBox{
	Postgres: `SELECT is_locked FROM migrations_lock WHERE index = 1;`,
	Mysql:    `SELECT is_locked FROM migrations_lock WHERE ` + "`index`" + ` = 1;`,
}

// Provide a narrative to the user about the current status of mig
// Inspired by `git status` and `brew doctor`
func CommandStatus(cfg config.MigConfig) error {

	// Attempt to connect to database

	dbox, err := database.Connect(cfg.Connection)

	if err != nil {
		return err
	}

	// Check if migration tables exist

	existMigrations := false

	err = dbox.QueryRow(EXIST_MIGRATIONS).Scan(&existMigrations)

	if err != nil {
		color.Red("unable to tell if 'migrations' table exists!")
		return err
	}

	existLock := false

	err = dbox.QueryRow(EXIST_LOCK).Scan(&existLock)

	if err != nil {
		color.Red("unable to tell if 'migrations_lock' table exists!")
		return err
	}

	if !existMigrations && !existLock {
		color.Red("The tables used for tracking migrations are missing.")
		color.White("This likely means that mig hasn't yet been initialized.")
		color.White("This can be solved by running the following command:")
		color.White("$ mig init")
		return errors.New("missing migration tables")
	}

	if !existMigrations {
		color.Red("The 'migrations' table is missing but the 'migrations_lock' table is present!")
		color.White("This might mean that data has been corrupted and that migration status is missing.")
		color.White("Consider looking into the root cause of the problem.")
		color.White("The quickest fix is to delete the lock table and initialize again:")
		color.White("> DROP TABLE migrations_lock;")
		color.White("$ mig init")
		return errors.New("missing migrations table")
	}

	if !existLock {
		color.Red("The 'migrations' table is present but the 'migrations_lock' table is missing!")
		color.White("This might mean that data has been corrupted.")
		color.White("Consider looking into the cause of the problem.")
		color.White("The quickest fix is backup the migrations table data, initialize again, then restore the data:")
		color.White("> ALTER TABLE migrations RENAME TO migrations_backup;")
		color.White("$ mig init")
		color.White("> DROP TABLE migrations;")
		color.White("> ALTER TABLE migrations_backup RENAME TO migrations;")
		return errors.New("missing migrations_lock table")
	}

	if dbox.Type == "mysql" {
		// The following gnarly comparison checks need to be rebuilt first
		color.Yellow("migration table description check is currently unimplemented for mysql.")
	} else {
		rows, err := dbox.Query(DESCRIBE)

		if err != nil {
			color.Red("unable to describe the migration tables!")
			return err
		}

		// Check if tables have correct columns

		var table, column, data string

		// struggling to find the Go-way to do this...
		rows.Next()
		rows.Scan(&table, &column, &data)
		if table != "migrations" || column != "batch" || data != "integer" {
			return errors.New("expected migrations.batch of type integer")
		}

		rows.Next()
		rows.Scan(&table, &column, &data)
		if table != "migrations" || column != "id" || data != "integer" {
			return errors.New("expected migrations.id of type integer")
		}

		rows.Next()
		rows.Scan(&table, &column, &data)
		if table != "migrations" || column != "migration_time" || data != "timestamp with time zone" {
			return errors.New("expected migrations.migration_time of type timestamp with time zone")
		}

		rows.Next()
		rows.Scan(&table, &column, &data)
		if table != "migrations" || column != "name" || data != "character varying" {
			return errors.New("expected migrations.name of type character varying")
		}

		rows.Next()
		rows.Scan(&table, &column, &data)
		if table != "migrations_lock" || column != "index" || data != "integer" {
			return errors.New("expected migrations_lock.index of type integer")
		}

		rows.Next()
		rows.Scan(&table, &column, &data)
		if table != "migrations_lock" || column != "is_locked" || data != "integer" {
			return errors.New("expected migrations_lock.is_locked of type integer")
		}
	}

	// Check if locked

	locked := false
	err = dbox.QueryRow(LOCK_STATUS).Scan(&locked)

	if err != nil {
		color.Red("unable to determine lock status!")
		return err
	}

	if locked {
		color.Red("Migrations are currently locked!")
		color.White("It could be that a migration is in progress. However it could also mean that a migration failed.")
		color.White("If migrations remain locked then someone will want to investigate the failed migration.")
		color.White("Once that's over you can unlock migrations by running the following:")
		color.White("$ mig unlock")
		color.White("")
		// Note: Don't need to return at this point
	}

	// Display the name of the last run migration

	status, err := migrations.GetStatus(cfg, dbox)

	if err != nil {
		color.Red("unable to determine migration status!")
		return err
	}

	if status.Last.Name != "" {
		color.HiWhite("Last Migration: %s (id=%d,batch=%d) on %s", status.Last.Name, status.Last.Id, status.Last.Batch, status.Last.Time.Format(time.RFC3339))
		fmt.Println()
	}

	color.White("Applied: %d, Unapplied: %d, Skipped: %d, Missing: %d", status.Applied, status.Unapplied, status.Skipped, status.Missing)
	fmt.Println()

	if status.Skipped > 0 {
		color.Red("There are at least one skipped migrations! Mig will not be able to run migrations until this is fixed.")
		color.White("A skipped migration happens when a local migration file is older than the most recently run migration.")
		color.White("To fix this, rename any skipped migrations so that their timestamps are newer.")
		color.White("Run this command to list skipped migrations:")
		color.White("$ mig list")
		return errors.New("encountered skipped migrations")
	}

	if status.Next != "" {
		color.HiWhite("Next Migration: %s", status.Next)
		color.White("To run this migration, execute the following command:")
		color.White("$ mig up")
	}

	return nil
}
