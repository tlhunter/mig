package database

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type DbBox struct {
	Db   *sql.DB
	Type string
}

func (dbox DbBox) GetQuery(qb QueryBox) string {
	return qb.For(dbox.Type)
}

func (dbox DbBox) Exec(qb QueryBox, args ...any) (sql.Result, error) {
	return dbox.Db.Exec(qb.For(dbox.Type), args...)
}

func (dbox DbBox) Query(qb QueryBox, args ...any) (*sql.Rows, error) {
	return dbox.Db.Query(qb.For(dbox.Type), args...)
}

func (dbox DbBox) QueryRow(qb QueryBox, args ...any) *sql.Row {
	return dbox.Db.QueryRow(qb.For(dbox.Type), args...)
}

// mig needs a common TLS flag mapping across all RDBMS
// Postgres
//   verify -> verify-full
//   insecure -> require
//   disable -> disable
// MySQL
//   verify -> true
//   insecure -> skip-verify
//   disable -> false

func Connect(connection string) (DbBox, error) {
	var dbox DbBox
	u, err := url.Parse(connection)

	dbox.Type = u.Scheme

	qs, err := url.ParseQuery(u.RawQuery)
	tls_in := qs.Get("tls")

	if err != nil {
		return dbox, errors.New("unable to parse connection url!")
	}

	if u.Scheme == "postgresql" {
		tls := "disable"

		if tls_in == "verify" {
			tls = "verify-full"
		} else if tls_in == "insecure" {
			tls = "require"
		}

		port := "5432"
		if u.Port() != "" {
			port = u.Port()
		}

		dbox.Db, err = sql.Open("postgres", fmt.Sprintf("postgresql://%s@%s:%s%s?sslmode=%s", u.User, u.Hostname(), port, u.Path, tls))

		if err != nil {
			return dbox, errors.New("unable to connect to postgresql database!")
		}

		err = dbox.Db.Ping()

		if err != nil {
			return dbox, errors.New("unable to connect to postgresql database!")
		}
	} else if u.Scheme == "mysql" {
		port := "3306"
		if u.Port() != "" {
			port = u.Port()
		}

		tls := "false"

		if tls_in == "verify" {
			tls = "true"
		} else if tls_in == "insecure" {
			tls = "skip-verify"
		}

		// multiStatements=true required to run multiple queries in a single call, basically all migrations
		mysqlConnString := fmt.Sprintf("%s@tcp(%s:%s)%s?tls=%s&multiStatements=true&parseTime=true", u.User, u.Hostname(), port, u.Path, tls)

		dbox.Db, err = sql.Open("mysql", mysqlConnString)

		if err != nil {
			return dbox, errors.New("unable to connect to mysql database!")
		}

		err := dbox.Db.Ping()

		if err != nil {
			return dbox, errors.New("unable to connect to mysql database!")
		}
	} else if u.Scheme == "sqlite" { // or sqlite3?
		// TODO: What should the connection URL look like?
		// sqlite://user:pass@watever/file.db
		// sqlite://user:pass@watever//tmp/file.db
		// sqlite3://user:pass@watever/./file.db
		dbox.Db, err = sql.Open("sqlite3", ":memory:")
		if err != nil {
			return dbox, errors.New("unable to connect to sqlite database!")
		}
	} else {
		return dbox, errors.New(fmt.Sprintf("mig doesn't support the '%s' database", u.Scheme))
	}

	return dbox, nil
}
