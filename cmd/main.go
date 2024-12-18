package main

import (
	"log"

	"github.com/dinklen08/golang_calc/internal/application"
)

func main() {
	app := application.New()

	err := app.Run()
	if err != nil {
		log.Printf("[FATAL] %v", err)
	}
}
