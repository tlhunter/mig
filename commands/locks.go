package commands

import (
	"fmt"

	"github.com/tlhunter/mig/config"
	"github.com/tlhunter/mig/database"
)

const (
	LOCK   = `UPDATE migrations_lock SET is_locked = 1 WHERE index = 1 RETURNING ( SELECT is_locked AS was_locked FROM migrations_lock WHERE index = 1);`
	UNLOCK = `UPDATE migrations_lock SET is_locked = 0 WHERE index = 1 RETURNING ( SELECT is_locked AS was_locked FROM migrations_lock WHERE index = 1);`
)

func CommandLock(cfg config.MigConfig) error {
	db := database.Connect(cfg.Connection)
	defer db.Close()

	var was_locked int
	err := db.QueryRow(LOCK).Scan(&was_locked)

	if err != nil {
		fmt.Println("unable to lock!", err)
		return err
	}

	if was_locked == 0 {
		fmt.Println("successfully locked.")
		return nil
	}

	fmt.Println("already locked!")

	return nil
}

func CommandUnlock(cfg config.MigConfig) error {
	db := database.Connect(cfg.Connection)

	defer db.Close()

	var was_locked int
	err := db.QueryRow(UNLOCK).Scan(&was_locked)

	if err != nil {
		fmt.Println("unable to unlock!", err)
		return err
	}

	if was_locked == 1 {
		fmt.Println("successfully unlocked.")
		return nil
	}

	fmt.Println("already unlocked!")

	return nil
}
