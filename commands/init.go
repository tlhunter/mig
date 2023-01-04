package commands

import (
	"github.com/fatih/color"
	"github.com/tlhunter/mig/config"
	"github.com/tlhunter/mig/database"
)

var (
	INIT = database.QueryBox{
		Postgres: `CREATE TABLE migrations (
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
INSERT INTO migrations_lock ("index", is_locked) VALUES(1, 0);`,
		Mysql: `CREATE TABLE migrations (
	id serial NOT NULL PRIMARY KEY,
	name varchar(255) NULL,
	batch int4 NULL,
	migration_time TIMESTAMP NULL
);
CREATE TABLE migrations_lock (
	` + "`index`" + ` serial NOT NULL PRIMARY KEY,
	is_locked int4 NULL
);
INSERT INTO migrations_lock SET ` + "`index`" + ` = 1, is_locked = 0;`,
	}
)

func CommandInit(cfg config.MigConfig) error {
	dbox := database.Connect(cfg.Connection)

	defer dbox.Db.Close()

	_, err := dbox.Exec(INIT)

	if err != nil {
		color.Red("error initializing mig!", err)
		return err
	}

	color.Green("successfully initialized mig.")

	return nil
}
