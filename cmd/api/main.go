package main

import (
	"log"

	"open-replays/api/internal/api/run"
)

func main() {
	if err := run.Run(); err != nil {
		log.Fatal(err)
	}
}
