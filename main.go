package main

import (
	"open-replays/api/internal"
)

func main() {
	r := internal.New()

	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
