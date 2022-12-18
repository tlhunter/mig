package main

import (
	"fmt"

	"github.com/tlhunter/mig/config"
)

func main() {
	fmt.Println("Welcome to mig.")

	config.GetConfig()

	//pq.ParseURL
}
