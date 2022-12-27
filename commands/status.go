package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/tlhunter/mig/config"
	"github.com/tlhunter/mig/database"
	"github.com/tlhunter/mig/migrations"
)

const EXIST_MIGRATIONS = `SELECT EXISTS (
SELECT FROM
	pg_tables
WHERE
	schemaname = 'public' AND
	tablename  = 'migrations'
) AS table_exists;`

const EXIST_LOCK = `SELECT EXISTS (
	SELECT FROM
		pg_tables
	WHERE
		schemaname = 'public' AND
		tablename  = 'migrations_lock'
	) AS table_exists;`

const DESCRIBE = `SELECT
	table_name,
	column_name,
	data_type
FROM
	information_schema.columns
WHERE
	table_name = 'migrations' OR table_name = 'migrations_lock'
ORDER BY
	table_name, column_name;`

const LOCK_STATUS = `SELECT is_locked FROM migrations_lock WHERE INDEX = 1;`

// Provide a narrative to the user about the current status of mig. Inspired by `git status` and `brew doctor`.
func CommandStatus(cfg config.MigConfig) error {

	// Attempt to connect to database

	db := database.Connect(cfg.Connection)

	// Check if migration tables exist

	existMigrations := false

	err := db.QueryRow(EXIST_MIGRATIONS).Scan(&existMigrations)

	if err != nil {
		os.Stderr.WriteString("unable to tell if 'migrations' table exists!\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return err
	}

	existLock := false

	err = db.QueryRow(EXIST_LOCK).Scan(&existLock)

	if err != nil {
		os.Stderr.WriteString("unable to tell if 'migrations_lock' table exists!\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return err
	}

	if !existMigrations && !existLock {
		fmt.Println("The tables used for tracking migrations are missing.")
		fmt.Println("This likely means that mig hasn't yet been initialized.")
		fmt.Println("This can be solved by running the following command:")
		fmt.Println("$ mig init")
		return nil
	}

	if !existMigrations {
		os.Stderr.WriteString("The 'migrations' table is missing but the 'migrations_lock' table is present!\n")
		os.Stderr.WriteString("This might mean that data has been corrupted and that migration status is missing.\n")
		os.Stderr.WriteString("Consider looking into the root cause of the problem.\n")
		os.Stderr.WriteString("The quickest fix is to delete the lock table and initialize again:\n")
		os.Stderr.WriteString("> DROP TABLE migrations_lock;\n")
		os.Stderr.WriteString("$ mig init\n")
		return nil
	}

	if !existLock {
		os.Stderr.WriteString("The 'migrations' table is present but the 'migrations_lock' table is missing!\n")
		os.Stderr.WriteString("This might mean that data has been corrupted.\n")
		os.Stderr.WriteString("Consider looking into the cause of the problem.\n")
		os.Stderr.WriteString("The quickest fix is backup the migrations table data, initialize again, then restore the data:\n")
		os.Stderr.WriteString("> ALTER TABLE migrations RENAME TO migrations_backup;\n")
		os.Stderr.WriteString("$ mig init\n")
		os.Stderr.WriteString("> DROP TABLE migrations;\n")
		os.Stderr.WriteString("> ALTER TABLE migrations_backup RENAME TO migrations;\n")
		return nil
	}

	rows, err := db.Query(DESCRIBE)

	if err != nil {
		os.Stderr.WriteString("unable to describe the migration tables!\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return err
	}

	// Check if tables have correct columns

	var table, column, data string

	// struggling to find the Go-way to do this...
	rows.Next()
	rows.Scan(&table, &column, &data)
	if table != "migrations" || column != "batch" || data != "integer" {
		os.Stderr.WriteString("expected migrations.batch of type integer\n")
		return nil
	}

	rows.Next()
	rows.Scan(&table, &column, &data)
	if table != "migrations" || column != "id" || data != "integer" {
		os.Stderr.WriteString("expected migrations.id of type integer\n")
		return nil
	}

	rows.Next()
	rows.Scan(&table, &column, &data)
	if table != "migrations" || column != "migration_time" || data != "timestamp with time zone" {
		os.Stderr.WriteString("expected migrations.migration_time of type timestamp with time zone\n")
		return nil
	}

	rows.Next()
	rows.Scan(&table, &column, &data)
	if table != "migrations" || column != "name" || data != "character varying" {
		os.Stderr.WriteString("expected migrations.name of type character varying\n")
		return nil
	}

	rows.Next()
	rows.Scan(&table, &column, &data)
	if table != "migrations_lock" || column != "index" || data != "integer" {
		os.Stderr.WriteString("expected migrations_lock.index of type integer\n")
		return nil
	}

	rows.Next()
	rows.Scan(&table, &column, &data)
	if table != "migrations_lock" || column != "is_locked" || data != "integer" {
		os.Stderr.WriteString("expected migrations_lock.is_locked of type integer\n")
		return nil
	}

	// Check if locked

	locked := false
	err = db.QueryRow(LOCK_STATUS).Scan(&locked)

	if err != nil {
		os.Stderr.WriteString("unable to determine lock status!\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return err
	}

	if locked {
		os.Stderr.WriteString("Migrations are currently locked! It could be that a migration is in progress.\n")
		os.Stderr.WriteString("However it could also mean that a migration failed.\n")
		os.Stderr.WriteString("If migrations remain locked then someone will want to investigate the failed migration.\n")
		os.Stderr.WriteString("Once that's over you can unlock migrations by running the following:\n")
		os.Stderr.WriteString("$ mig unlock\n")
		os.Stderr.WriteString("\n")
		// Note: Don't need to return at this point
	}

	// Display the name of the last run migration

	status, err := migrations.GetStatus(cfg, db, false)

	if err != nil {
		os.Stderr.WriteString("unable to determine migration status!\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return err
	}

	if status.Last.Name != "" {
		fmt.Printf("Last Migration: %s (id=%d,batch=%d) on %s\n", status.Last.Name, status.Last.Id, status.Last.Batch, status.Last.Time.Format(time.RFC3339))
		fmt.Println()
	}

	fmt.Printf("Applied: %d, Unapplied: %d, Skipped: %d, Missing: %d\n", status.Applied, status.Unapplied, status.Skipped, status.Missing)
	fmt.Println()

	if status.Skipped > 0 {
		os.Stderr.WriteString("There are at least one skipped migrations! Mig will not be able to run migrations until this is fixed.\n")
		os.Stderr.WriteString("A skipped migration happens when a local migration file is older than the most recently run migration.\n")
		os.Stderr.WriteString("To fix this, rename any skipped migrations so that their timestamps are newer.\n")
		os.Stderr.WriteString("Run this command to list skipped migrations:\n")
		os.Stderr.WriteString("$ mig list\n")
		return nil
	}

	if status.Next != "" {
		fmt.Printf("Next Migration: %s\n", status.Next)
		os.Stderr.WriteString("To run this migration, execute the following command:\n")
		os.Stderr.WriteString("$ mig up\n")
		fmt.Println()
	}

	return nil
}
