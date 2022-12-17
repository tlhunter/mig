package config

import (
	"flag"
	"fmt"
)

func Flags() {
	var connection string
	flag.StringVar(&connection, "connection", "", "SQL Connection String")

	flag.Parse()
	// TODO: Fails when env var is missing. Will need to also check env vars and config files.
	// TODO: Error message lists -connection, prefer the more common --connection text.

	fmt.Println("Connection Flag:", connection)
}
