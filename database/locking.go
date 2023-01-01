package database

import "database/sql"

var (
	OBTAIN = QueryBox{
		Postgres: `UPDATE migrations_lock SET is_locked = 1 WHERE index = 1 AND is_locked = 0;`,
		Mysql:    `UPDATE migrations_lock SET is_locked = 1 WHERE ` + "`index`" + ` = 1 AND is_locked = 0;`,
	}
	RELEASE = QueryBox{
		Postgres: `UPDATE migrations_lock SET is_locked = 0 WHERE index = 1 AND is_locked = 1;`,
		Mysql:    `UPDATE migrations_lock SET is_locked = 0 WHERE ` + "`index`" + ` = 1 AND is_locked = 1;`,
	}
)

func ObtainLock(db *sql.DB, dbType string) (bool, error) {
	result, err := db.Exec(OBTAIN.For(dbType))

	if err != nil {
		return false, err
	}

	affected, err := result.RowsAffected()

	if err != nil {
		return false, err
	}

	return affected == 1, nil
}

func ReleaseLock(db *sql.DB, dbType string) (bool, error) {
	result, err := db.Exec(RELEASE.For(dbType))

	if err != nil {
		return false, err
	}

	affected, err := result.RowsAffected()

	if err != nil {
		return false, err
	}

	return affected == 1, nil
}
