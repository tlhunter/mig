package database

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"

	"github.com/fatih/color"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

// mig needs a common TLS flag mapping across all RDBMS
// Postgres
//   verify -> verify-full
//   insecure -> require
//   disable -> disable
// MySQL
//   verify -> true
//   insecure -> skip-verify
//   disable -> false

func Connect(connection string) (*sql.DB, string) {
	u, err := url.Parse(connection)

	dbType := u.Scheme

	qs, err := url.ParseQuery(u.RawQuery)
	tls_in := qs.Get("tls")

	if err != nil {
		color.Red("unable to parse connection url!\n")
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(2)
	}

	var db *sql.DB

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

		db, err = sql.Open("postgres", fmt.Sprintf("postgresql://%s@%s:%s%s?sslmode=%s", u.User, u.Hostname(), port, u.Path, tls))

		if err != nil {
			color.Red("unable to connect to postgresql database!\n")
			os.Stderr.WriteString(err.Error() + "\n")
			os.Exit(3)
		}

		err = db.Ping()

		if err != nil {
			color.Red("unable to connect to postgresql database!\n")
			os.Stderr.WriteString(err.Error() + "\n")
			os.Exit(4)
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

		db, err = sql.Open("mysql", mysqlConnString)

		if err != nil {
			color.Red("unable to connect to mysql database!\n")
			os.Stderr.WriteString(err.Error() + "\n")
			os.Exit(3)
		}

		err := db.Ping()

		if err != nil {
			color.Red("unable to connect to mysql database!\n")
			os.Stderr.WriteString(err.Error() + "\n")
			os.Exit(4)
		}
	} else {
		color.Red("mig doesn't support the '%s' database", u.Scheme)
		os.Exit(5)
	}

	return db, dbType
}
