package migrations

import (
	"time"

	"github.com/tlhunter/mig/database"
)

var LIST = database.QueryBox{
	Postgres: `SELECT id, name, batch, migration_time FROM migrations ORDER BY id ASC;`,
	Mysql:    `SELECT id, name, batch, migration_time FROM migrations ORDER BY id ASC;`,
}

type MigrationRow struct {
	Id    int
	Name  string
	Batch int
	Time  time.Time
}

func ListRows(dbox database.DbBox) ([]MigrationRow, error) {
	var migRows []MigrationRow

	rows, err := dbox.Query(LIST)

	if err != nil {
		return migRows, err
	}

	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		var batch int
		var time time.Time
		err = rows.Scan(&id, &name, &batch, &time)

		if err != nil {
			return migRows, err
		}

		migRows = append(migRows, MigrationRow{
			Id:    id,
			Name:  name,
			Batch: batch,
			Time:  time,
		})
	}

	return migRows, nil
}
