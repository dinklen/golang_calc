package application

import (
	"net/http"
	"os"

	"golang_calc/internal/calc_libs/expressions"
	golang_calc "golang_calc/internal/database"
	"golang_calc/internal/orchestrator"

	"github.com/gorilla/mux"
)

var App *Application

func init() {
	App = New()
}

type Config struct {
	AppPort   string
	AgentPort string

	PlusTime     string
	MinusTime    string
	MultipTime   string
	DivisionTime string

	ComputingPower string

	Database *database.Database
}

func ConfigFill() *Config {
	config := new(Config)

	// env
	config.AppPort = os.Getenv("APP_PORT")
	config.AgentPort = os.Getenv("AGENT_PORT")

	config.PlusTime = os.Getenv("TIME_ADDITION_MS")
	config.MinusTime = os.Getenv("TIME_SUBSTRACTION_MS")
	config.MultipTime = os.Getenv("TIME_MULTIPLICATION_MS")
	config.DivisionTime = os.Getenv("TIME_DIVISION_MS")

	config.ComputingPower = os.Getenv("COMPUTING_POWER")

	if config.AppPort == "" {
		config.AppPort = "8080"
	}

	if config.AgentPort == "" {
		config.AgentPort = "8081"
	} // dangerous

	// default
	db, err := database.NewDatabase()
	if err != nil {
		panic(err)
	}

	config.Database = db

	return config
}

type Application struct {
	Configuration *Config
}

func New() *Application {
	return &Application{
		Configuration: ConfigFill(),
	}
}

func (app *Application) RunApp() error {
	defer App.Configuration.Database.DB.Close()

	router := mux.NewRouter()

	router.HandleFunc("/api/v1/calculate", orchestrator.CalcHandler).Methods("POST")
	router.HandleFunc("/api/v1/expressions", expressions.CurrentExpressionHandler).Methods("GET")
	router.HandleFunc("/api/v1/expressions/{id}", expressions.AllExpressionsHandler).Methods("GET")

	err := http.ListenAndServe(":"+app.Configuration.AppPort, nil)

	return err
}

func (app *Application) RunAgent() error {
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/internal/task", golang_calc.AgentHandler).Methods("GET", "POST")

	err := http.ListenAndServe(":"+app.Configuration.AgentPort, nil) // cange port

	return err
}
