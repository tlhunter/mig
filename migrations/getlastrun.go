package migrations

import (
	"database/sql"
	"time"
)

// TODO: This is redundant with status.Last

const QUERY = `SELECT id, name, batch, migration_time FROM migrations ORDER BY id DESC LIMIT 1;`

func GetLastRun(db *sql.DB) (MigrationRow, error) {

	var migration MigrationRow

	var id int
	var name string
	var batch int
	var time time.Time

	err := db.QueryRow(QUERY).Scan(&id, &name, &batch, &time)

	if err != nil {
		return migration, err
	}

	migration.Id = id
	migration.Name = name
	migration.Batch = batch
	migration.Time = time

	return migration, nil
}
