package main

import (
	"log"

	"slam-cli/internal/app"
)

func main() {
	if err := app.Execute(); err != nil {
		log.Fatal(err)
	}
}
