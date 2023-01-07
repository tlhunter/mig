package commands

import (
	"github.com/fatih/color"
	"github.com/tlhunter/mig/config"
	"github.com/tlhunter/mig/database"
)

var (
	LOCK = database.QueryBox{
		Postgres: `UPDATE migrations_lock SET is_locked = 1 WHERE index = 1 RETURNING ( SELECT is_locked AS was_locked FROM migrations_lock WHERE index = 1);`,
		Mysql: `START TRANSACTION;
	SELECT is_locked AS was_locked FROM migrations_lock WHERE ` + "`index`" + ` = 1;
	UPDATE migrations_lock SET is_locked = 1 WHERE ` + "`index`" + ` = 1;
COMMIT;`,
	}
	UNLOCK = database.QueryBox{
		Postgres: `UPDATE migrations_lock SET is_locked = 0 WHERE index = 1 RETURNING ( SELECT is_locked AS was_locked FROM migrations_lock WHERE index = 1);`,
		Mysql: `START TRANSACTION;
	SELECT is_locked AS was_locked FROM migrations_lock WHERE ` + "`index`" + ` = 1;
	UPDATE migrations_lock SET is_locked = 0 WHERE ` + "`index`" + ` = 1;
COMMIT;`,
	}
)

func CommandLock(cfg config.MigConfig) error {
	dbox, err := database.Connect(cfg.Connection)

	if err != nil {
		return err
	}

	defer dbox.Db.Close()

	var was_locked int
	err = dbox.QueryRow(LOCK).Scan(&was_locked)

	if err != nil {
		color.Red("mig: unable to lock!")
		return err
	}

	if was_locked == 0 {
		color.Green("mig: successfully locked.")
		return nil
	}

	color.Yellow("mig: already locked!")

	return nil
}

func CommandUnlock(cfg config.MigConfig) error {
	dbox, err := database.Connect(cfg.Connection)

	if err != nil {
		return err
	}

	defer dbox.Db.Close()

	var was_locked int
	err = dbox.QueryRow(UNLOCK).Scan(&was_locked)

	if err != nil {
		color.Red("mig: unable to unlock!")
		return err
	}

	if was_locked == 1 {
		color.Green("mig: successfully unlocked.")
		return nil
	}

	color.Yellow("mig: already unlocked!")

	return nil
}
