package main

import (
	"fmt"
	"os"

	"github.com/tlhunter/mig/commands"
	"github.com/tlhunter/mig/config"
)

func main() {
	cfg, err := config.GetConfig()

	if err != nil {
		os.Stderr.WriteString("unable to parse configuration") // TODO: print err
		os.Exit(1)
	}

	fmt.Printf("cfg: %v\n", cfg)

	commands.Dispatch(cfg)

	// queries, err := migrations.GetQueriesFromFile("./tests/20221219184300-example.sql")

	// fmt.Println("UP")
	// fmt.Println(queries.Up)
	// fmt.Println("DOWN")
	// fmt.Println(queries.Down)

}
