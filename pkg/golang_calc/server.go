package golang_calc 

import (
	"net/http"
	"encoding/json"
	"log"
	"fmt"
	"errors"
	"io"
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

	w.WriteHeader(errCode)
	err := json.NewEncoder(w).Encode(
		failureOutputData {
			Error: errText,
		},
	)

	if err != nil {
		log.Printf("[ERROR] %v", err)
		w.WriteHeader(501)
		return
	}
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	var (
		decryptData inputData
		result float64

		err error
	)

	if r.Method != "POST" {
		errorOutput(w, "Access denied", 405, errors.New("try to use method GET"))
		return
	}

	defer r.Body.Close()

	data, err := io.ReadAll(r.Body)
	log.Printf("[info] %s;%v", string(data), err)
	if err != nil {
		errorOutput(w, "Internal server error", 500, err)
		return
	}
	
	json.Unmarshal(data, &decryptData)

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
