package database

import "fmt"

// QueryBox makes it easy to store multiple queries for different databases.
// Use it when an operation can always be completed with a single query for each database type.
type QueryBox struct {
	Postgres string
	Mysql    string
}

func (qb QueryBox) For(driver string) string {
	if driver == "mysql" {
		return qb.Mysql
	} else if driver == "postgresql" {
		return qb.Postgres
	}

	panic(fmt.Sprintf("requested query for unknown driver %s", driver))
}
