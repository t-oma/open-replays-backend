package main

import (
	"log"

	"open-replays/internal/api/run"
)

func main() {
	if err := run.Run(); err != nil {
		log.Fatal(err)
	}
}
