package commands

import (
	"fmt"
	"os"
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

	db, dbType := database.Connect(cfg.Connection)

	// Check if migration tables exist

	existMigrations := false

	err := db.QueryRow(EXIST_MIGRATIONS.For(dbType)).Scan(&existMigrations)

	if err != nil {
		color.Red("unable to tell if 'migrations' table exists!")
		os.Stderr.WriteString(err.Error() + "\n")
		return err
	}

	existLock := false

	err = db.QueryRow(EXIST_LOCK.For(dbType)).Scan(&existLock)

	if err != nil {
		color.Red("unable to tell if 'migrations_lock' table exists!\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return err
	}

	if !existMigrations && !existLock {
		color.Red("The tables used for tracking migrations are missing.")
		color.White("This likely means that mig hasn't yet been initialized.")
		color.White("This can be solved by running the following command:")
		color.White("$ mig init")
		return nil
	}

	if !existMigrations {
		color.Red("The 'migrations' table is missing but the 'migrations_lock' table is present!\n")
		color.White("This might mean that data has been corrupted and that migration status is missing.\n")
		color.White("Consider looking into the root cause of the problem.\n")
		color.White("The quickest fix is to delete the lock table and initialize again:\n")
		color.White("> DROP TABLE migrations_lock;\n")
		color.White("$ mig init\n")
		return nil
	}

	if !existLock {
		color.Red("The 'migrations' table is present but the 'migrations_lock' table is missing!\n")
		color.White("This might mean that data has been corrupted.\n")
		color.White("Consider looking into the cause of the problem.\n")
		color.White("The quickest fix is backup the migrations table data, initialize again, then restore the data:\n")
		color.White("> ALTER TABLE migrations RENAME TO migrations_backup;\n")
		color.White("$ mig init\n")
		color.White("> DROP TABLE migrations;\n")
		color.White("> ALTER TABLE migrations_backup RENAME TO migrations;\n")
		return nil
	}

	if dbType == "mysql" {
		// The following gnarly comparison checks need to be rebuilt first
		color.Yellow("migration table description check is currently unimplemented for mysql.\n")
	} else {
		rows, err := db.Query(DESCRIBE.For(dbType))

		if err != nil {
			color.Red("unable to describe the migration tables!\n")
			os.Stderr.WriteString(err.Error() + "\n")
			return err
		}

		// Check if tables have correct columns

		var table, column, data string

		// struggling to find the Go-way to do this...
		rows.Next()
		rows.Scan(&table, &column, &data)
		if table != "migrations" || column != "batch" || data != "integer" {
			color.Red("expected migrations.batch of type integer\n")
			return nil
		}

		rows.Next()
		rows.Scan(&table, &column, &data)
		if table != "migrations" || column != "id" || data != "integer" {
			color.Red("expected migrations.id of type integer\n")
			return nil
		}

		rows.Next()
		rows.Scan(&table, &column, &data)
		if table != "migrations" || column != "migration_time" || data != "timestamp with time zone" {
			color.Red("expected migrations.migration_time of type timestamp with time zone\n")
			return nil
		}

		rows.Next()
		rows.Scan(&table, &column, &data)
		if table != "migrations" || column != "name" || data != "character varying" {
			color.Red("expected migrations.name of type character varying\n")
			return nil
		}

		rows.Next()
		rows.Scan(&table, &column, &data)
		if table != "migrations_lock" || column != "index" || data != "integer" {
			color.Red("expected migrations_lock.index of type integer\n")
			return nil
		}

		rows.Next()
		rows.Scan(&table, &column, &data)
		if table != "migrations_lock" || column != "is_locked" || data != "integer" {
			color.Red("expected migrations_lock.is_locked of type integer\n")
			return nil
		}
	}

	// Check if locked

	locked := false
	err = db.QueryRow(LOCK_STATUS.For(dbType)).Scan(&locked)

	if err != nil {
		color.Red("unable to determine lock status!\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return err
	}

	if locked {
		color.Red("Migrations are currently locked!")
		color.White("It could be that a migration is in progress. However it could also mean that a migration failed.\n")
		color.White("If migrations remain locked then someone will want to investigate the failed migration.\n")
		color.White("Once that's over you can unlock migrations by running the following:\n")
		color.White("$ mig unlock\n")
		color.White("\n")
		// Note: Don't need to return at this point
	}

	// Display the name of the last run migration

	status, err := migrations.GetStatus(cfg, db, false)

	if err != nil {
		color.Red("unable to determine migration status!")
		os.Stderr.WriteString(err.Error() + "\n")
		return err
	}

	if status.Last.Name != "" {
		color.HiWhite("Last Migration: %s (id=%d,batch=%d) on %s", status.Last.Name, status.Last.Id, status.Last.Batch, status.Last.Time.Format(time.RFC3339))
		fmt.Println()
	}

	color.White("Applied: %d, Unapplied: %d, Skipped: %d, Missing: %d\n", status.Applied, status.Unapplied, status.Skipped, status.Missing)
	fmt.Println()

	if status.Skipped > 0 {
		color.Red("There are at least one skipped migrations! Mig will not be able to run migrations until this is fixed.")
		color.White("A skipped migration happens when a local migration file is older than the most recently run migration.\n")
		color.White("To fix this, rename any skipped migrations so that their timestamps are newer.\n")
		color.White("Run this command to list skipped migrations:\n")
		color.White("$ mig list\n")
		return nil
	}

	if status.Next != "" {
		color.HiWhite("Next Migration: %s", status.Next)
		color.White("To run this migration, execute the following command:\n")
		color.White("$ mig up\n")
	}

	return nil
}
