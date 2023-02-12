package commands

import (
	"github.com/tlhunter/mig/config"
	"github.com/tlhunter/mig/database"
	"github.com/tlhunter/mig/result"
)

var (
	UNLOCK = database.QueryBox{
		Postgres: `UPDATE migrations_lock SET is_locked = 0 WHERE index = 1 RETURNING ( SELECT is_locked AS was_locked FROM migrations_lock WHERE index = 1);`,
		Mysql: `START TRANSACTION;
	SELECT is_locked AS was_locked FROM migrations_lock WHERE ` + "`index`" + ` = 1;
	UPDATE migrations_lock SET is_locked = 0 WHERE ` + "`index`" + ` = 1;
COMMIT;`,
	}
)

func CommandLock(cfg config.MigConfig) result.Response {
	dbox, err := database.Connect(cfg.Connection)

	if err != nil {
		return *result.NewErrorWithDetails("database connection error", "db_conn", err)
	}

	defer dbox.Db.Close()

	var was_locked bool

	if dbox.IsPostgres {
		was_locked, err = postgresLock(dbox)
	} else if dbox.IsMysql {
		was_locked, err = mysqlLock(dbox)
	} else {
		panic("unknown database: " + dbox.Type)
	}

	if err != nil {
		return *result.NewErrorWithDetails("unable to lock!", "unable_lock", err)
	}

	if !was_locked {
		return *result.NewSuccess("successfully locked.")
	}

	return *result.NewSuccess("already locked!") // TODO: yellow
}

func postgresLock(dbox database.DbBox) (bool, error) {
	var was_locked int
	err := dbox.Db.QueryRow(`UPDATE migrations_lock SET is_locked = 1 WHERE index = 1 RETURNING ( SELECT is_locked AS was_locked FROM migrations_lock WHERE index = 1);`).Scan(&was_locked)
	if err != nil {
		return false, err
	}

	return was_locked > 0, nil
}

func mysqlLock(dbox database.DbBox) (bool, error) {
	tx, err := dbox.Db.Begin()
	if err != nil {
		return false, err
	}

	defer tx.Rollback()

	var was_locked int

	err = tx.QueryRow(`SELECT is_locked AS was_locked FROM migrations_lock WHERE ` + "`index`" + ` = 1;`).Scan(&was_locked)
	if err != nil {
		return false, err
	}

	_, err = tx.Exec(`UPDATE migrations_lock SET is_locked = 1 WHERE ` + "`index`" + ` = 1;`)

	if err = tx.Commit(); err != nil {
		return false, err
	}

	return was_locked > 0, nil
}

func CommandUnlock(cfg config.MigConfig) result.Response {
	dbox, err := database.Connect(cfg.Connection)

	if err != nil {
		return *result.NewErrorWithDetails("database connection error", "db_conn", err)
	}

	defer dbox.Db.Close()

	var was_locked int
	err = dbox.QueryRow(UNLOCK).Scan(&was_locked)

	if err != nil {
		return *result.NewErrorWithDetails("unable to unlock!", "unable_unlock", err)
	}

	if was_locked == 1 {
		return *result.NewSuccess("successfully unlocked.")
	}

	return *result.NewSuccess("already unlocked!") // TODO: yellow
}
