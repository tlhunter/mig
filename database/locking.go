package database

import "database/sql"

const (
	OBTAIN  = `UPDATE migrations_lock SET is_locked = 1 WHERE index = 1 AND is_locked = 0;`
	RELEASE = `UPDATE migrations_lock SET is_locked = 0 WHERE index = 1 AND is_locked = 1;`
)

func ObtainLock(db *sql.DB) (bool, error) {
	result, err := db.Exec(OBTAIN)

	if err != nil {
		return false, err
	}

	affected, err := result.RowsAffected()

	if err != nil {
		return false, err
	}

	return affected == 1, nil
}

func ReleaseLock(db *sql.DB) (bool, error) {
	result, err := db.Exec(RELEASE)

	if err != nil {
		return false, err
	}

	affected, err := result.RowsAffected()

	if err != nil {
		return false, err
	}

	return affected == 1, nil
}
