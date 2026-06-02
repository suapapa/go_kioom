package main

import (
	"log"
	"os"

	"github.com/suapapa/go_kioom/internal/kioomenv"
)

func main() {
	// Load environment variables from .env if present
	if err := kioomenv.LoadEnvFile(".env"); err != nil {
		log.Printf("warning: failed to load .env file: %v", err)
	}

	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr, os.Getenv))
}
