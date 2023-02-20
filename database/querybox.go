package database

import "fmt"

type QueryBox struct {
	Postgres string
	Mysql    string
	Sqlite   string
}

func (qb QueryBox) For(driver string) string {
	if driver == "mysql" {
		return qb.Mysql
	} else if driver == "postgresql" {
		return qb.Postgres
	} else if driver == "sqlite" {
		return qb.Sqlite
	}

	panic(fmt.Sprintf("requested query for unknown driver %s", driver))
}
