package database

import (
	"database/sql"
	"fmt"
	"net/url"
	"os"

	"github.com/fatih/color"

	_ "github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
)

func Connect(connection string) (*sql.DB, string) {
	u, err := url.Parse(connection)

	// TODO: ?tls=verify|insecure|disable
	// defaults to disable

	dbType := u.Scheme

	if err != nil {
		color.Red("unable to parse connection url!\n")
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(2)
	}

	var db *sql.DB

	if u.Scheme == "postgresql" {
		parsed, err := pq.ParseURL(connection)

		if err != nil {
			color.Red("unable to parse postgresql connection string!\n")
			os.Stderr.WriteString(err.Error() + "\n")
			os.Exit(2)
		}

		db, err = sql.Open("postgres", parsed)

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

		// multiStatements=true required to run multiple queries in a single call, basically all migrations
		mysqlConnString := fmt.Sprintf("%s@tcp(%s:%s)%s?tls=%s&multiStatements=true&parseTime=true", u.User, u.Host, port, u.Path, "skip-verify")

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
