package golang_calc 

import (
	"net/http"
	"encoding/json"
	"log"
	"fmt"
	"errors"
)

type successOutputData struct {
	Result string `json:"result"`
}

type failureOutputData struct {
	Error string `json:"error"`
}

type inputData struct {
	Expression string `json:"expression"`
}

func errorOutput(w http.ResponseWriter, errText string, errCode int, errEvent error) {
	log.Printf("[ERROR] %v", errEvent)
	
	err := json.NewEncoder(w).Encode(
		failureOutputData {
			Error: errText,
		},
	)

	if err != nil {
		log.Printf("[ERROR] %v")
		w.WriteHeader(501)
		return
	}

	w.WriteHeader(errCode)
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	var (
		decryptData inputData
		result float64

		err error
	)

	if r.Method != http.MethodPost {
		errorOutput(w, "Access denied", 405, errors.New("try to use method GET"))
		return
	}

	defer r.Body.Close()
	
	if err = json.NewDecoder(r.Body).Decode(&decryptData); err != nil {
		errorOutput(w, "Internal server error", 500, err)
		return
	}
		
	result, err = Calc(decryptData.Expression)
	if err != nil {
		errorOutput(w, "Expression is not valid", 422, err)
		return
	}



	err = json.NewEncoder(w).Encode(
		successOutputData {
			Result: fmt.Sprint(result),
		},
	)

	if err != nil {
		errorOutput(w, "Internal server error", 500, err)
		return
	}

	w.WriteHeader(200)
	log.Printf("[INFO] success")
}
