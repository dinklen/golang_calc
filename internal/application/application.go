package application

import (
	"log"
	"net/http"

	"github.com/dinklen08/golang_calc/pkg/golang_calc"
)

type Application struct {}

func New() *Application {
	return &Application{}
}

func (app *Application) Run() error {
	http.HandleFunc("/api/v1/calculate", server.CalcHandler)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Printf("[ERROR] %v", err)
	}
}
