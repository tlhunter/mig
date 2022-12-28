package migrations

import "database/sql"

const (
	HIGHEST = `SELECT (SELECT batch FROM migrations ORDER BY batch DESC LIMIT 1) AS highest_batch, (SELECT id FROM migrations ORDER BY id DESC LIMIT 1) AS highest_id;`
	ADD     = `INSERT INTO migrations (id, name, batch, migration_time) VALUES ($1, $2, $3, NOW()) RETURNING id, name, batch, migration_time;`
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
func AddMigrations(db *sql.DB, migrations []string) error {
	// TODO: multiple inserts with different IDs and shared Batch

	return nil
}

// down
func RemoveMigration(db *sql.DB, migration string) error {
	// Ensure that the provided migration is the final migration
	// If it's not then fail

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
