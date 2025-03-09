package expr_handler

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"strconv"

	"golang_calc/internal/calc_libs/expressions"
	"golang_calc/internal/database"

	"github.com/gorilla/mux"
)

type Tasks struct {
	Exprs []*expressions.ExpressionInfo `json:"expressions"`
}

func CurrentExpressionsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	totalID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		log.Fatal("[FATAL] failed to get id")
		w.WriteHeader(422)
		return
	}

	expr, err := database.DataBase.UnloadCurrentTask(uint32(totalID))
	if err != nil {
		w.Write([]byte(`{"error":"invalid id"}`))
		w.WriteHeader(404)
		return
	}

	if expr.Result == math.MaxFloat64 {
		expr.Status = "error"
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
		w.WriteHeader(500)
		return
	}

	encExprs := &Tasks{Exprs: exprs}
	if len(encExprs.Exprs) == 1 && encExprs.Exprs[0].Result == math.MaxFloat64 {
		encExprs.Exprs[0].Status = "error"
	}

	err = json.NewEncoder(w).Encode(encExprs)
	if err != nil {
		log.Printf("[ERROR] failed to encode expressions info: %v", err)
		w.WriteHeader(500)
		return
	}

	log.Printf("[INFO] success to send message with the expressions info")
	w.WriteHeader(http.StatusOK)
}
