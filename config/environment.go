package config

import "os"

func Environment() {
	connection := os.Getenv("MIG_CONNECTION")

	println("Environment Variable:", connection)
}
