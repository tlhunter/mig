package main

import (
	"fmt"

	"github.com/tlhunter/mig/config"
)

func main() {
	fmt.Println("Welcome to mig.")

	config.Flags()
	config.Environment()
	config.File()
}
