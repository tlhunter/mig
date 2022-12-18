package main

import (
	"fmt"

	"github.com/lib/pq"
	"github.com/tlhunter/mig/config"
)

func main() {
	fmt.Println("Welcome to mig.")

	cfg, err := config.GetConfig()
	fmt.Printf("cfg: %v\n", cfg)

	if err != nil {
		fmt.Println("shit")
	}

	parsed, err := pq.ParseURL(cfg.Connection)

	if err != nil {
		fmt.Println("unable to parse connection string")
	}

	fmt.Printf("parsed: %v\n", parsed)
}
