package commands

import (
	"github.com/tlhunter/mig/config"
	"github.com/tlhunter/mig/database"
	"github.com/tlhunter/mig/result"
)

var (
	LOCK = database.QueryBox{
		Postgres: `UPDATE migrations_lock SET is_locked = 1 WHERE index = 1 RETURNING ( SELECT is_locked AS was_locked FROM migrations_lock WHERE index = 1);`,
		Mysql: `START TRANSACTION;
	SELECT is_locked AS was_locked FROM migrations_lock WHERE ` + "`index`" + ` = 1;
	UPDATE migrations_lock SET is_locked = 1 WHERE ` + "`index`" + ` = 1;
COMMIT;`,
		Sqlite: `BEGIN TRANSACTION;
	SELECT is_locked AS was_locked FROM migrations_lock WHERE "index" = 1;
	UPDATE migrations_lock SET is_locked = 1 WHERE "index" = 1;
COMMIT TRANSACTION;`, // TODO: Does this even work?
	}
	UNLOCK = database.QueryBox{
		Postgres: `UPDATE migrations_lock SET is_locked = 0 WHERE index = 1 RETURNING ( SELECT is_locked AS was_locked FROM migrations_lock WHERE index = 1);`,
		Mysql: `START TRANSACTION;
	SELECT is_locked AS was_locked FROM migrations_lock WHERE ` + "`index`" + ` = 1;
	UPDATE migrations_lock SET is_locked = 0 WHERE ` + "`index`" + ` = 1;
COMMIT;`,
		Sqlite: `BEGIN TRANSACTION;
	SELECT is_locked AS was_locked FROM migrations_lock WHERE "index" = 1;
	UPDATE migrations_lock SET is_locked = 0 WHERE "index" = 1;
COMMIT TRANSACTION;`, // TODO: Does this even work?
	}
)

func CommandLock(cfg config.MigConfig) result.Response {
	dbox, err := database.Connect(cfg.Connection)

	if err != nil {
		return *result.NewErrorWithDetails("database connection error", "db_conn", err)
	}

	defer dbox.Db.Close()

	var was_locked int
	err = dbox.QueryRow(LOCK).Scan(&was_locked)

	if err != nil {
		return *result.NewErrorWithDetails("unable to lock!", "unable_lock", err)
	}

	if was_locked == 0 {
		return *result.NewSuccess("successfully locked.")
	}

	return *result.NewSuccess("already locked!") // TODO: yellow
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
