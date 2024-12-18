package application

import (
	"net/http"

	"github.com/dinklen08/golang_calc/pkg/golang_calc"
)

type Application struct {}

func New() *Application {
	return &Application{}
}

func (app *Application) Run() error {
	http.HandleFunc("/api/v1/calculate", golang_calc.CalcHandler)

	err := http.ListenAndServe(":8080", nil)

	return err
}
