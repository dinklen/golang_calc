package golang_calc 

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"log"
	"fmt"

	"github.com/dinklen08/golang_calc/pkg/golang_calc"
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
		FailureOutputData {
			Error: errText,
		},
	)

	if err != nil {
		log.Printf("[FATAL] %v")
		w.WriteHeader(501)
		return
	}

	w.WriteHeader(errCode)
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	var (
		decryptData inputData
		result float64
	)

	if r.Method != http.MethodPost {
		errorOutput(w, "Access denied", 405, errors.New("try to use method GET"))
		return
	}

	defer r.Body.Close()

	encryptData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errorOutput(w, "Internal server error", 500, err)
		return
	}

	if err = json.NewDecoder(encryptData).Decode(decryptData); err != nil {
		errorOutput(w, "Internal server error", 500, err)
		return
	}
		
	result, err = calc.Calc(decryptData.Expression)
	if err != nil {
		errorOutput(w, "Expression is not valid", 422, err)
		return
	}

	json.NewEncoder(w).Encode(
		successOutputData {
			Result: fmt.Sprintf("%.f", result),
		},
	)

	w.WriteHeader(200)
	log.Printf("[INFO] success")
}
