package migrations

import (
	"errors"

	"github.com/tlhunter/mig/database"
)

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
func AddMigration(dbox database.DbBox, migration string) error {
	Highest, err := GetHighestValues(dbox)

	if err != nil {
		return err
	}

	_, err = dbox.Exec(ADD, Highest.Id, migration, Highest.Batch)

	if err != nil {
		return err
	}

	return nil
}

// upto, all
func AddMigrationWithBatch(dbox database.DbBox, migration string, group int) error {
	Highest, err := GetHighestValues(dbox)

	if err != nil {
		return err
	}

	_, err = dbox.Exec(ADD, Highest.Id, migration, group)

	if err != nil {
		return err
	}

	return nil
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
	var Highest BatchAndId

	var count int

	err := dbox.QueryRow(COUNT).Scan(&count)

	if err != nil {
		return Highest, err
	}

	if count == 0 {
		// First migration
		Highest.Id = 1
		Highest.Batch = 1

		return Highest, nil
	}

	err = dbox.QueryRow(HIGHEST).Scan(&Highest.Batch, &Highest.Id)

	if err != nil {
		return Highest, err
	}

	Highest.Id++
	Highest.Batch++

	return Highest, nil
}
