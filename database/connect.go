package database

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/lib/pq"
)

func Connect(connection string) *sql.DB {
	parsed, err := pq.ParseURL(connection)

	if err != nil {
		os.Stderr.WriteString("unable to parse connection string")
		os.Exit(2)
	}

	// fmt.Printf("parsed: %v\n", parsed)

	db, err := sql.Open("postgres", parsed)

	// defer db.Close()

	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}

	// rows, err := db.Query(`SELECT 1 AS foo`)

	// defer rows.Close()

	// if err != nil {
	// os.Stderr.WriteString("unable to query database")
	// os.Exit(4)
	// }

	err = db.Ping()

	if err != nil {
		os.Stderr.WriteString("unable to ping database")
		os.Exit(4)
	}

	return db
}
