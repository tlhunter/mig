package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/lib/pq"
	"github.com/tlhunter/mig/config"
	"github.com/tlhunter/mig/migrations"
)

func main() {
	cfg, err := config.GetConfig()

	if err != nil {
		os.Stderr.WriteString("unable to parse configuration") // TODO: print error
		os.Exit(1)
	}

	fmt.Printf("cfg: %v\n", cfg)

	queries, err := migrations.GetQueriesFromFile("./tests/20221219184300-example.sql")

	fmt.Println("UP")
	fmt.Println(queries.Up)
	fmt.Println("DOWN")
	fmt.Println(queries.Down)

	parsed, err := pq.ParseURL(cfg.Connection)

	if err != nil {
		os.Stderr.WriteString("unable to parse connection string")
		os.Exit(2)
	}

	fmt.Printf("parsed: %v\n", parsed)

	db, err := sql.Open("postgres", parsed)

	defer db.Close()

	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}

	rows, err := db.Query(`SELECT 1 AS foo`)

	defer rows.Close()

	if err != nil {
		os.Stderr.WriteString("unable to query database")
		os.Exit(4)
	}

	err = db.Ping()

	if err != nil {
		os.Stderr.WriteString("unable to ping database")
		os.Exit(4)
	}
}
