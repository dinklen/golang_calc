package main

import (
	"log"

	"golang_calc/internal/application"
)

func main() {
	app := application.New()

	if err := app.RunApp(); err != nil {
		log.Fatal("failed to start app: ", err)
	}
}
