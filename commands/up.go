package commands

import (
	"os"

	"github.com/fatih/color"
	"github.com/tlhunter/mig/config"
	"github.com/tlhunter/mig/database"
	"github.com/tlhunter/mig/migrations"
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

func CommandUp(cfg config.MigConfig) error {
	db, dbType := database.Connect(cfg.Connection)

	defer db.Close()

	status, err := migrations.GetStatus(cfg, db, false)

	if err != nil {
		color.Red("Encountered an error trying to get migrations status!\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return err
	}

	if status.Skipped > 0 {
		color.Red("Refusing to run with skipped migrations! Run `mig status` for details.\n")
		return nil
	}

	next := status.Next

	if next == "" {
		color.Red("There are no migrations to run.")
		return nil
	}

	filename := cfg.Migrations + "/" + next

	queries, err := migrations.GetQueriesFromFile(filename)

	if err != nil {
		color.Red("Error attempting to read next migration file!\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return err
	}

	locked, err := database.ObtainLock(db)

	if err != nil {
		color.Red("Error obtaining lock for migration!\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return err
	}

	if !locked {
		color.Red("Unable to obtain lock for migration!\n")
		return nil
	}

	var query string

	if queries.UpTx {
		query = BEGIN.For(dbType) + queries.Up + END.For(dbType)
	} else {
		query = queries.Up
	}

	_, err = db.Exec(query)

	if err != nil {
		color.Red("Encountered an error while running migration!\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return err
	}

	color.Green("Migration %s was successfully applied!\n", next)

	err = migrations.AddMigration(db, next)

	if err != nil {
		color.Red("The migration query executed but unable to track it in the migrations table!\n")
		color.White("You may want to manually add it and investigate the error.\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return err
	}

	released, err := database.ReleaseLock(db)

	if err != nil {
		color.Red("Error obtaining lock for migration!\n")
		os.Stderr.WriteString(err.Error() + "\n")
		return err
	}

	if !released {
		color.Red("Unable to obtain lock for migration!\n")
		return nil
	}

	return nil
}
