package migrations

import (
	"time"

	"github.com/tlhunter/mig/database"
)

type MigrationRow struct {
	Id    int        `json:"id,omitempty"`
	Name  string     `json:"name"`
	Batch int        `json:"batch,omitempty"`
	Time  *time.Time `json:"time,omitempty"`
}

func ListRows(dbox database.DbBox) ([]MigrationRow, error) {
	var migRows []MigrationRow

	if !dbox.IsMysql && !dbox.IsPostgres {
		panic("unknown database: " + dbox.Type)
	}

	rows, err := dbox.Db.Query("SELECT id, name, batch, migration_time FROM migrations ORDER BY id ASC;") // query is the same for MySQL and Postgres
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
			Time:  &time,
		})
	}

	return migRows, nil
}
