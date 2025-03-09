package orchestrator

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type outputData interface {
	GetData() string
}

type successOutputData struct {
	Id uint32 `json:"id"`
}

type failureOutputData struct {
	Error string `json:"error"`
}

type inputData struct {
	Expression string `json:"expression"`
}

// successOutputData methods
func (sod successOutputData) GetData() string {
	return fmt.Sprintf("%v", sod.Id)
}

// failureOutputData methods
func (fod failureOutputData) GetData() string {
	return fod.Error
}

func ErrorOutput(w http.ResponseWriter, errText string, errCode int, errEvent error) {
	log.Printf("[ERROR] %v", errEvent)

	w.WriteHeader(errCode)
	err := json.NewEncoder(w).Encode(
		failureOutputData{
			Error: errText,
		},
	)

	if err != nil {
		log.Printf("[ERROR] %v", err)
		w.WriteHeader(500)
		return
	}
}
