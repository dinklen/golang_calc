package expr_handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"golang_calc/internal/calc_libs/expressions"
	"golang_calc/internal/database"

	"github.com/gorilla/mux"
)

type Tasks struct {
	Exprs []*expressions.ExpressionInfo `json:"tasks"`
}

func CurrentExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	totalID, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal("failed to get id")
		return

		// error_output
	}

	expr, err := database.DataBase.UnloadCurrentTask(totalID)
	if err != nil {
		return // error_output
	}

	err = json.NewEncoder(w).Encode(expr)
	if err != nil {
		log.Printf("[ERROR] failed to encode expression info: %v", err)
		return
	}

	log.Printf("[INFO] success to send message with the expression info")
	w.WriteHeader(http.StatusOK)
}

func AllExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	exprs, err := database.DataBase.UnloadAllTasks()
	if err != nil {
		return // error_output
	}

	encExprs := &Tasks{Exprs: exprs}

	err = json.NewEncoder(w).Encode(encExprs)
	if err != nil {
		// error_output
		log.Printf("[ERROR] failed to encode expressions info: %v", err)
		return
	}

	log.Printf("[INFO] success to send message with the expressions info")
	w.WriteHeader(http.StatusOK)
}
