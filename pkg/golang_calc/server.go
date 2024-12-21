package golang_calc 

import (
	"net/http"
	"encoding/json"
	"log"
	"fmt"
	"errors"
	"io"
)

var (
	ErrIncorrectMethod error = errors.New("incorrect method")
	ErrIncorrectQuery error = errors.New("incorrect query")
)

type outputData interface {
	GetData() string
}

type successOutputData struct {
	Result string `json:"result"`
}

type failureOutputData struct {
	Error string `json:"error"`
}

type inputData struct {
	Expression string `json:"expression"`
}

//successOutputData methods
func (sod successOutputData) GetData() string {
	return sod.Result
}

//failureOutputData methods
func (fod failureOutputData) GetData() string {
	return fod.Error
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
		w.WriteHeader(500)
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
		errorOutput(w, fmt.Sprintf("Internal server error: %v", ErrIncorrectMethod), 500, ErrIncorrectMethod)
		return
	}

	defer r.Body.Close()

	data, err := io.ReadAll(r.Body)	
	if err != nil {
		errorOutput(w, fmt.Sprintf("Internal server error: %v", err), 500, err)
		return
	}
	
	err = json.Unmarshal(data, &decryptData)
	if err != nil {
		errorOutput(w, fmt.Sprintf("Internal server error: %v", err), 500, err)
		return
	}

	if string(data) != "{\"expression\":\"\"}" && decryptData.Expression == "" {
		errorOutput(w, fmt.Sprintf("Internal server error: %v", ErrIncorrectQuery), 500, ErrIncorrectQuery)
		return
	}

	result, err = Calc(decryptData.Expression)
	if err != nil {
		errorOutput(w, fmt.Sprintf("Expression is not valid: %v", err), 422, err)
		return
	}

	err = json.NewEncoder(w).Encode(
		successOutputData {
			Result: fmt.Sprint(result),
		},
	)

	if err != nil {
		errorOutput(w, fmt.Sprintf("Internal server error: %v", err), 500, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("[INFO] success")
}
