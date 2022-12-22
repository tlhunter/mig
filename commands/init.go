package commands

import (
	"fmt"

	"github.com/tlhunter/mig/config"
	"github.com/tlhunter/mig/database"
)

const INIT = `CREATE TABLE migrations (
	id serial NOT NULL,
	name varchar(255) NULL,
	batch int4 NULL,
	migration_time timestamptz NULL,
	CONSTRAINT migrations_pkey PRIMARY KEY (id)
);
CREATE TABLE migrations_lock (
	"index" serial NOT NULL,
	is_locked int4 NULL,
	CONSTRAINT migrations_lock_pkey PRIMARY KEY (index)
);
INSERT INTO migrations_lock ("index", is_locked) VALUES(1, 0);`

func CommandInit(cfg config.MigConfig) error {
	db := database.Connect(cfg.Connection)

	defer db.Close()

	_, err := db.Exec(INIT)

	if err != nil {
		fmt.Println("error initializing mig!", err)
		return err
	}

	fmt.Println("successfully initialized mig.")

	return nil
}
