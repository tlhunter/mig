package commands

import (
	"github.com/tlhunter/mig/config"
	"github.com/tlhunter/mig/database"
	"github.com/tlhunter/mig/result"
)

func CommandInit(cfg config.MigConfig) result.Response {
	dbox, err := database.Connect(cfg.Connection)

	if err != nil {
		return *result.NewErrorWithDetails("database connection error", "db_conn", err)
	}

	defer dbox.Db.Close()

	if dbox.Type == "postgresql" {
		err = postgresInit(cfg, dbox)
	} else if dbox.Type == "mysql" {
		err = mysqlInit(cfg, dbox)
	}

	if err != nil {
		return *result.NewErrorWithDetails("error initializing mig!", "unable_init", err)
	}

	return *result.NewSuccess("successfully initialized mig")
}

func postgresInit(cfg config.MigConfig, dbox database.DbBox) error {
	_, err := dbox.Db.Exec(`CREATE TABLE migrations (
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
INSERT INTO migrations_lock ("index", is_locked) VALUES(1, 0);`)

	return err
}

func mysqlInit(cfg config.MigConfig, dbox database.DbBox) error {
	_, err := dbox.Db.Exec(`CREATE TABLE migrations (
	id serial NOT NULL PRIMARY KEY,
	name varchar(255) NULL,
	batch int4 NULL,
	migration_time TIMESTAMP NULL
);
CREATE TABLE migrations_lock (
	` + "`index`" + ` serial NOT NULL PRIMARY KEY,
	is_locked int4 NULL
);
INSERT INTO migrations_lock SET ` + "`index`" + ` = 1, is_locked = 0;`)

	return err
}
