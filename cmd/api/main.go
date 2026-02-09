// Package main is the entry point for the Open-Replays API server.
package main

import (
	"log"

	"open-replays/internal/api/run"
)

// main starts the API server.
func main() {
	if err := run.Run(); err != nil {
		log.Fatal(err)
	}
}
