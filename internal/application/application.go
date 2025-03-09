package application

import (
	"log"
	"net/http"

	"golang_calc/internal/agent"
	"golang_calc/internal/calc_libs/expr_handler"
	"golang_calc/internal/config"
	"golang_calc/internal/database"
	"golang_calc/internal/orchestrator"

	"github.com/gorilla/mux"
)

type Application struct {
	Configuration *config.Config
}

func New() *Application {
	db, err := database.NewDatabase()
	if err != nil {
		log.Fatal("failed to create config")
		panic(err)
	}

	database.DataBase = db

	config.Conf = config.ConfigFill()

	return &Application{
		Configuration: config.Conf,
	}
}

func (app *Application) RunApp() error {
	defer database.DataBase.DB.Close()

	err := database.DataBase.Clean()
	if err != nil {
		panic(err)
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/v1/calculate", orchestrator.CalcHandler).Methods("POST")
	router.HandleFunc("/api/v1/expressions", expr_handler.AllExpressionsHandler).Methods("GET")
	router.HandleFunc("/api/v1/expressions/{id}", expr_handler.CurrentExpressionsHandler).Methods("GET")

	log.Printf("[INFO] app:listening... (port=%s)", app.Configuration.AppPort)
	err = http.ListenAndServe(":"+app.Configuration.AppPort, router)

	return err
}

func (app *Application) RunAgent() error {
	router := mux.NewRouter()

	router.HandleFunc("/internal/task", agent.AgentHandler).Methods("POST")

	log.Printf("[INFO] agent:listening... (port=%s)", app.Configuration.AgentPort)
	err := http.ListenAndServe(":"+app.Configuration.AgentPort, router)

	return err
}
