package database

import (
	"database/sql"
	"os"

	"github.com/lib/pq"
)

func Connect(connection string) *sql.DB {
	parsed, err := pq.ParseURL(connection)

	if err != nil {
		os.Stderr.WriteString("unable to parse connection string!\n")
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(2)
	}

	db, err := sql.Open("postgres", parsed)

	if err != nil {
		os.Stderr.WriteString("unable to connect to database!\n")
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(3)
	}

	err = db.Ping()

	if err != nil {
		os.Stderr.WriteString("unable to connect to database!\n")
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(4)
	}

	return db
}
