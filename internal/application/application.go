package application

import (
	"net/http"
	"os"

	"github.com/dinklen08/golang_calc/pkg/golang_calc"
)

type Config struct {
	Addr string
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = os.Getenv("PORT")

	if config.Addr == "" {
		config.Addr = "8080"
	}

	return config
}

type Application struct {
	config *Config
}

func New() *Application {
	return &Application{
		config: ConfigFromEnv(),
	}
}

func (app *Application) Run() error {
	http.HandleFunc("/api/v1/calculate", golang_calc.CalcHandler)

	err := http.ListenAndServe(":"+app.config.Addr, nil)

	return err
}
