package main

import (
	"log"

	"golang_calc/internal/application"
)

func main() {
	if err := application.App.RunApp(); err != nil {
		log.Fatal("failed to start app: ", err)
	}
}
