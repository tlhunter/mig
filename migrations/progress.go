package migrations

import (
	"database/sql"
	"errors"
)

const (
	HIGHEST  = `SELECT (SELECT batch FROM migrations ORDER BY batch DESC LIMIT 1) AS highest_batch, (SELECT id FROM migrations ORDER BY id DESC LIMIT 1) AS highest_id;`
	ADD      = `INSERT INTO migrations (id, name, batch, migration_time) VALUES ($1, $2, $3, NOW()) RETURNING id, name, batch, migration_time;`
	ULTIMATE = `SELECT id, name FROM migrations ORDER BY id DESC LIMIT 1;`
	DELETE   = `DELETE FROM migrations WHERE id = $1 AND name = $2;`
)

type BatchAndId struct {
	Batch int
	Id    int
}

// up
func AddMigration(db *sql.DB, migration string) error {
	Highest, err := GetHighestValues(db)

	if err != nil {
		return err
	}

	_, err = db.Exec(ADD, Highest.Id, migration, Highest.Batch)

	if err != nil {
		return err
	}

	return nil
}

// upto, all
func AddMigrationWithBatch(db *sql.DB, migration string, group int) error {
	Highest, err := GetHighestValues(db)

	if err != nil {
		return err
	}

	_, err = db.Exec(ADD, Highest.Id, migration, group)

	if err != nil {
		return err
	}

	return nil
}

// down
func RemoveMigration(db *sql.DB, migration string, id int) error {
	// Ensure that the provided migration is the final migration
	// If it's not then fail
	var lastId int
	var lastName string

	err := db.QueryRow(ULTIMATE).Scan(&lastId, &lastName)

	if err != nil {
		return err
	}

	if lastId != id || migration != lastName {
		return errors.New("Tried to delete the non-final migration")
	}

	result, err := db.Exec(DELETE, id, migration)

	affected, err := result.RowsAffected()

	if err != nil {
		return err
	}

	if affected != 1 {
		return errors.New("Unable to remove migration from migrations table")
	}

	return nil
}

func GetHighestValues(db *sql.DB) (BatchAndId, error) {
	var Highest BatchAndId

	err := db.QueryRow(HIGHEST).Scan(&Highest.Batch, &Highest.Id)

	if err != nil {
		return Highest, err
	}

	Highest.Id++
	Highest.Batch++

	return Highest, nil
}
