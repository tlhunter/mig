package migrations

import (
	"errors"

	"github.com/tlhunter/mig/database"
)

var (
	HIGHEST = database.QueryBox{
		Postgres: `SELECT (SELECT batch FROM migrations ORDER BY batch DESC LIMIT 1) AS highest_batch, (SELECT id FROM migrations ORDER BY id DESC LIMIT 1) AS highest_id;`,
		Mysql:    `SELECT (SELECT batch FROM migrations ORDER BY batch DESC LIMIT 1) AS highest_batch, (SELECT id FROM migrations ORDER BY id DESC LIMIT 1) AS highest_id;`,
		Sqlite:   `SELECT (SELECT batch FROM migrations ORDER BY batch DESC LIMIT 1) AS highest_batch, (SELECT id FROM migrations ORDER BY id DESC LIMIT 1) AS highest_id;`, // TODO: This returns two nils instead of an empty row?
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

	if dbox.IsPostgres {
		migration, err = postgresAddMigration(dbox, highest.Id, migrationName, highest.Batch)
	} else if dbox.IsMysql {
		migration, err = mysqlAddMigration(dbox, highest.Id, migrationName, highest.Batch)
	} else if dbox.IsSqlite {
		migration, err = sqliteAddMigration(dbox, highest.Id, migrationName, highest.Batch)
	} else {
		panic("unknown database: " + dbox.Type)
	}

	if err != nil {
		return migration, err
	}

	return migration, nil
}

func postgresAddMigration(dbox database.DbBox, id int, name string, batchId int) (MigrationRow, error) {
	var migration MigrationRow

	err := dbox.Db.
		QueryRow(`INSERT INTO migrations (id, name, batch, migration_time) VALUES ($1, $2, $3, NOW()) RETURNING id, name, batch, migration_time;`, id, name, batchId).
		Scan(&migration.Id, &migration.Name, &migration.Batch, &migration.Time)

	return migration, err
}

func mysqlAddMigration(dbox database.DbBox, id int, name string, batchId int) (MigrationRow, error) {
	var migration MigrationRow

	tx, err := dbox.Db.Begin()
	if err != nil {
		return migration, err
	}

	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO migrations (id, name, batch, migration_time) VALUES (?, ?, ?, NOW());", id, name, batchId)
	if err != nil {
		return migration, err
	}

	// Technically, the only value we need is the time, since that's the only value that the database determins
	err = tx.
		QueryRow("SELECT id, name, batch, migration_time FROM migrations WHERE id = ?;", id).
		Scan(&migration.Id, &migration.Name, &migration.Batch, &migration.Time)
	if err != nil {
		return migration, err
	}

	if err = tx.Commit(); err != nil {
		return migration, err
	}

	return migration, nil
}

func sqliteAddMigration(dbox database.DbBox, id int, name string, batchId int) (MigrationRow, error) {
	var migration MigrationRow

	err := dbox.Db.
		QueryRow(`INSERT INTO migrations (id, name, batch, migration_time) VALUES (?, ?, ?, CURRENT_TIMESTAMP) RETURNING id, name, batch, migration_time;`, id, name, batchId).
		Scan(&migration.Id, &migration.Name, &migration.Batch, &migration.Time)

	return migration, err
}

// upto, all
func AddMigrationWithBatch(dbox database.DbBox, migrationName string, batch int) (MigrationRow, error) {
	var migration MigrationRow

	highest, err := GetHighestValues(dbox)
	if err != nil {
		return migration, err
	}

	if dbox.IsPostgres {
		migration, err = postgresAddMigration(dbox, highest.Id, migrationName, batch)
	} else if dbox.IsMysql {
		migration, err = mysqlAddMigration(dbox, highest.Id, migrationName, batch)
	} else if dbox.IsSqlite {
		migration, err = sqliteAddMigration(dbox, highest.Id, migrationName, batch)
	} else {
		panic("unknown database: " + dbox.Type)
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
	if err != nil {
		return err
	}

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
