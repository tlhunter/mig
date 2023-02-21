package migrations

import (
	"errors"
	"fmt"
	"time"

	"github.com/tlhunter/mig/database"
)

var (
	HIGHEST = database.QueryBox{
		Postgres: `SELECT (SELECT batch FROM migrations ORDER BY batch DESC LIMIT 1) AS highest_batch, (SELECT id FROM migrations ORDER BY id DESC LIMIT 1) AS highest_id;`,
		Mysql:    `SELECT (SELECT batch FROM migrations ORDER BY batch DESC LIMIT 1) AS highest_batch, (SELECT id FROM migrations ORDER BY id DESC LIMIT 1) AS highest_id;`,
		Sqlite:   `SELECT (SELECT batch FROM migrations ORDER BY batch DESC LIMIT 1) AS highest_batch, (SELECT id FROM migrations ORDER BY id DESC LIMIT 1) AS highest_id;`, // TODO: This returns two nils instead of an empty row?
	}
	ADD = database.QueryBox{
		Postgres: `INSERT INTO migrations (id, name, batch, migration_time) VALUES ($1, $2, $3, NOW()) RETURNING id, name, batch, migration_time;`,
		Mysql:    `INSERT INTO migrations (id, name, batch, migration_time) VALUES (?, ?, ?, NOW());`, // TODO: Transaction
		Sqlite:   `INSERT INTO migrations (id, name, batch, migration_time) VALUES (?, ?, ?, CURRENT_TIMESTAMP) RETURNING id, name, batch, migration_time;`,
	}
	ULTIMATE = database.QueryBox{
		Postgres: `SELECT id, name FROM migrations ORDER BY id DESC LIMIT 1;`,
		Mysql:    `SELECT id, name FROM migrations ORDER BY id DESC LIMIT 1;`,
		Sqlite:   `SELECT id, name FROM migrations ORDER BY id DESC LIMIT 1;`,
	}
	DELETE = database.QueryBox{
		Postgres: `DELETE FROM migrations WHERE id = $1 AND name = $2;`,
		Mysql:    `DELETE FROM migrations WHERE id = ? AND name = ?;`,
		Sqlite:   `DELETE FROM migrations WHERE id = ? AND name = ?;`,
	}
	COUNT = database.QueryBox{
		Postgres: `SELECT COUNT(*) AS count FROM migrations;`,
		Mysql:    `SELECT COUNT(*) AS count FROM migrations;`,
		Sqlite:   `SELECT COUNT(*) AS count FROM migrations;`,
	}
)

type BatchAndId struct {
	Batch int
	Id    int
}

// up
func AddMigration(dbox database.DbBox, migrationName string) (MigrationRow, error) {
	var migration MigrationRow

	highest, err := GetHighestValues(dbox)

	if err != nil {
		return migration, err
	}

	if dbox.IsMysql { // TODO: Transaction
		// MySQL provides no easy RETURNING equivalent, so we'll fake it and omit the calculated timestamp
		_, err = dbox.Exec(ADD, highest.Id, migrationName, highest.Batch)
		migration.Id = highest.Id
		migration.Name = migrationName
		migration.Batch = highest.Batch
		migration.Time = nil
	} else if dbox.IsSqlite {
		// This branch is needed as sqlite doesn't provide column type information when using a RETURNING clause.
		// This means that we need to manually convert the returned timestamp from a string to a time.Time.
		// @see https://github.com/mattn/go-sqlite3/issues/951
		var tempTime string
		err = dbox.QueryRow(ADD, highest.Id, migrationName, highest.Batch).Scan(&migration.Id, &migration.Name, &migration.Batch, &tempTime)
		fmt.Println(tempTime)
		parsed, err := time.Parse(time.DateTime, tempTime)
		if err != nil {

		}
		migration.Time = &parsed
	} else {
		err = dbox.QueryRow(ADD, highest.Id, migrationName, highest.Batch).Scan(&migration.Id, &migration.Name, &migration.Batch, &migration.Time)
	}

	if err != nil {
		return migration, err
	}

	return migration, nil
}

// upto, all
func AddMigrationWithBatch(dbox database.DbBox, migrationName string, batch int) (MigrationRow, error) {
	var migration MigrationRow

	highest, err := GetHighestValues(dbox)

	if err != nil {
		return migration, err
	}

	if dbox.IsMysql { // TODO: Transaction
		// MySQL provides no easy RETURNING equivalent, so we'll fake it and omit the calculated timestamp
		_, err = dbox.Exec(ADD, highest.Id, migrationName, highest.Batch)
		migration.Id = highest.Id
		migration.Name = migrationName
		migration.Batch = batch
		migration.Time = nil
	} else if dbox.IsSqlite {
		// This branch is needed as sqlite doesn't provide column type information when using a RETURNING clause.
		// This means that we need to manually convert the returned timestamp from a string to a time.Time.
		// @see https://github.com/mattn/go-sqlite3/issues/951
		var tempTime string
		err = dbox.QueryRow(ADD, highest.Id, migrationName, highest.Batch).Scan(&migration.Id, &migration.Name, &migration.Batch, &tempTime)
		fmt.Println(tempTime)
		parsed, err := time.Parse(time.DateTime, tempTime)
		if err != nil {

		}
		migration.Time = &parsed
	} else {
		err = dbox.QueryRow(ADD, highest.Id, migrationName, batch).Scan(&migration.Id, &migration.Name, &migration.Batch, &migration.Time)
	}
	if err != nil {
		return migration, err
	}

	return migration, nil
}

// down
func RemoveMigration(dbox database.DbBox, migration string, id int) error {
	// Ensure that the provided migration is the final migration
	// If it's not then fail
	var lastId int
	var lastName string

	err := dbox.QueryRow(ULTIMATE).Scan(&lastId, &lastName)

	if err != nil {
		return err
	}

	if lastId != id || migration != lastName {
		return errors.New("Tried to delete the non-final migration")
	}

	result, err := dbox.Exec(DELETE, id, migration)

	affected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if affected != 1 {
		return errors.New("Unable to remove migration from migrations table")
	}

	return nil
}

func GetHighestValues(dbox database.DbBox) (BatchAndId, error) {
	var highest BatchAndId

	var count int

	err := dbox.QueryRow(COUNT).Scan(&count)

	if err != nil {
		return highest, err
	}

	if count == 0 {
		// First migration
		highest.Id = 1
		highest.Batch = 1

		return highest, nil
	}

	err = dbox.QueryRow(HIGHEST).Scan(&highest.Batch, &highest.Id)

	if err != nil {
		return highest, err
	}

	highest.Id++
	highest.Batch++

	return highest, nil
}
