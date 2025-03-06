package expressions

import (
	"encoding/json"
	"log"
	"net/http"

	"golang_calc/internal/application"

	"github.com/gorilla/mux"
	// include db, application
)

type Tasks struct {
	Exprs []ExpressionInfo `json:"tasks"`
}

func CurrentExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	expr, err := application.App.Configuration.Database.DB.UnloadCurrentTask(id)
	if err != nil {
		return // error_output
	}

	err = json.Encoder(w).NewEncoder(expr)
	if err != nil {
		log.Printf("[ERROR] failed to encode expression info: %v", err)
		return
	}

	log.Printf("[INFO] success to send message with the expression info")
	w.WriteHeader(http.StatusOK)
}

func AllAxpressionsHandler(w http.ResponseWriter, r *http.Request) {
	exprs, err := application.App.Configuration.Database.DB.UnloadAllTasks()
	if err != nil {
		return // error_output
	}

	encExprs := &Tasks{Exprs: exprs}

	err = json.Encoder(w).NewEncoder(encExprs)
	if err != nil {
		// error_output
		log.Printf("[ERROR] failed to encode expressions info: %v", err)
		return
	}

	log.Printf("[INFO] success to send message with the expressions info")
	w.WriteHeader(http.StatusOK)
}
