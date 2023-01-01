package migrations

import (
	"database/sql"
	"errors"

	"github.com/tlhunter/mig/database"
)

// TODO: QueryBox
var (
	HIGHEST = database.QueryBox{
		Postgres: `SELECT (SELECT batch FROM migrations ORDER BY batch DESC LIMIT 1) AS highest_batch, (SELECT id FROM migrations ORDER BY id DESC LIMIT 1) AS highest_id;`,
		Mysql:    `SELECT (SELECT batch FROM migrations ORDER BY batch DESC LIMIT 1) AS highest_batch, (SELECT id FROM migrations ORDER BY id DESC LIMIT 1) AS highest_id;`,
	}
	ADD = database.QueryBox{
		Postgres: `INSERT INTO migrations (id, name, batch, migration_time) VALUES ($1, $2, $3, NOW());`,
		Mysql:    `INSERT INTO migrations (id, name, batch, migration_time) VALUES (?, ?, ?, NOW());`,
	}
	ULTIMATE = database.QueryBox{
		Postgres: `SELECT id, name FROM migrations ORDER BY id DESC LIMIT 1;`,
		Mysql:    `SELECT id, name FROM migrations ORDER BY id DESC LIMIT 1;`,
	}
	DELETE = database.QueryBox{
		Postgres: `DELETE FROM migrations WHERE id = $1 AND name = $2;`,
		Mysql:    `DELETE FROM migrations WHERE id = ? AND name = ?;`,
	}
	COUNT = database.QueryBox{
		Postgres: `SELECT COUNT(*) AS count FROM migrations;`,
		Mysql:    `SELECT COUNT(*) AS count FROM migrations;`,
	}
)

type BatchAndId struct {
	Batch int
	Id    int
}

// up
func AddMigration(db *sql.DB, migration string, dbType string) error {
	Highest, err := GetHighestValues(db, dbType)

	if err != nil {
		return err
	}

	_, err = db.Exec(ADD.For(dbType), Highest.Id, migration, Highest.Batch)

	if err != nil {
		return err
	}

	return nil
}

// upto, all
func AddMigrationWithBatch(db *sql.DB, migration string, group int, dbType string) error {
	Highest, err := GetHighestValues(db, dbType)

	if err != nil {
		return err
	}

	_, err = db.Exec(ADD.For(dbType), Highest.Id, migration, group)

	if err != nil {
		return err
	}

	return nil
}

// down
func RemoveMigration(db *sql.DB, migration string, id int, dbType string) error {
	// Ensure that the provided migration is the final migration
	// If it's not then fail
	var lastId int
	var lastName string

	err := db.QueryRow(ULTIMATE.For(dbType)).Scan(&lastId, &lastName)

	if err != nil {
		return err
	}

	if lastId != id || migration != lastName {
		return errors.New("Tried to delete the non-final migration")
	}

	result, err := db.Exec(DELETE.For(dbType), id, migration)

	affected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if affected != 1 {
		return errors.New("Unable to remove migration from migrations table")
	}

	return nil
}

func GetHighestValues(db *sql.DB, dbType string) (BatchAndId, error) {
	var Highest BatchAndId

	var count int

	err := db.QueryRow(COUNT.For(dbType)).Scan(&count)

	if err != nil {
		return Highest, err
	}

	if count == 0 {
		// First migration
		Highest.Id = 1
		Highest.Batch = 1

		return Highest, nil
	}

	err = db.QueryRow(HIGHEST.For(dbType)).Scan(&Highest.Batch, &Highest.Id)

	if err != nil {
		return Highest, err
	}

	Highest.Id++
	Highest.Batch++

	return Highest, nil
}
