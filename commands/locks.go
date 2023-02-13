package commands

import (
	"github.com/tlhunter/mig/config"
	"github.com/tlhunter/mig/database"
	"github.com/tlhunter/mig/result"
)

func CommandLock(cfg config.MigConfig) result.Response {
	dbox, err := database.Connect(cfg.Connection)

	if err != nil {
		return *result.NewErrorWithDetails("database connection error", "db_conn", err)
	}

	defer dbox.Db.Close()

	var wasLocked bool

	if dbox.IsPostgres {
		wasLocked, err = postgresLock(dbox)
	} else if dbox.IsMysql {
		wasLocked, err = mysqlLock(dbox)
	} else {
		panic("unknown database: " + dbox.Type)
	}

	if err != nil {
		return *result.NewErrorWithDetails("unable to lock!", "unable_lock", err)
	}

	if !wasLocked {
		return *result.NewSuccess("successfully locked.")
	}

	return *result.NewSuccess("already locked!") // TODO: yellow
}

func postgresLock(dbox database.DbBox) (bool, error) {
	var wasLocked int
	err := dbox.Db.QueryRow("UPDATE migrations_lock SET is_locked = 1 WHERE index = 1 RETURNING ( SELECT is_locked AS was_locked FROM migrations_lock WHERE index = 1);").Scan(&wasLocked)
	if err != nil {
		return false, err
	}

	return wasLocked > 0, nil
}

func mysqlLock(dbox database.DbBox) (bool, error) {
	tx, err := dbox.Db.Begin()
	if err != nil {
		return false, err
	}

	defer tx.Rollback()

	var wasLocked int

	err = tx.QueryRow("SELECT is_locked AS was_locked FROM migrations_lock WHERE `index` = 1;").Scan(&wasLocked)
	if err != nil {
		return false, err
	}

	_, err = tx.Exec("UPDATE migrations_lock SET is_locked = 1 WHERE `index` = 1;")

	if err = tx.Commit(); err != nil {
		return false, err
	}

	return wasLocked > 0, nil
}

func CommandUnlock(cfg config.MigConfig) result.Response {
	dbox, err := database.Connect(cfg.Connection)

	if err != nil {
		return *result.NewErrorWithDetails("database connection error", "db_conn", err)
	}

	defer dbox.Db.Close()

	var wasLocked bool

	if dbox.IsPostgres {
		wasLocked, err = postgresUnlock(dbox)
	} else if dbox.IsMysql {
		wasLocked, err = mysqlUnlock(dbox)
	} else {
		panic("unknown database: " + dbox.Type)
	}

	if err != nil {
		return *result.NewErrorWithDetails("unable to unlock!", "unable_unlock", err)
	}

	if wasLocked {
		return *result.NewSuccess("successfully unlocked.")
	}

	return *result.NewSuccess("already unlocked!") // TODO: yellow
}

func postgresUnlock(dbox database.DbBox) (bool, error) {
	var wasLocked int

	err := dbox.Db.QueryRow("UPDATE migrations_lock SET is_locked = 0 WHERE index = 1 RETURNING ( SELECT is_locked AS was_locked FROM migrations_lock WHERE index = 1);").Scan(&wasLocked)
	if err != nil {
		return false, err
	}

	return wasLocked > 0, nil
}

func mysqlUnlock(dbox database.DbBox) (bool, error) {
	tx, err := dbox.Db.Begin()
	if err != nil {
		return false, err
	}

	defer tx.Rollback()

	var wasLocked int

	err = tx.QueryRow("SELECT is_locked AS was_locked FROM migrations_lock WHERE `index` = 1;").Scan(&wasLocked)
	if err != nil {
		return false, err
	}

	_, err = tx.Exec("UPDATE migrations_lock SET is_locked = 0 WHERE `index` = 1;")
	if err = tx.Commit(); err != nil {
		return false, err
	}

	return wasLocked > 0, nil
}
