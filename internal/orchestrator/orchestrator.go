package orchestrator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"golang_calc/internal/application"
)

type idMap struct {
	Map map[string]int `json:"map"`
}

func errorOutput(w http.ResponseWriter, errText string, errCode int, errEvent error) {
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

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	var (
		decryptData inputData
		result      float64

		err error
	)

	if r.Method != "POST" {
		errorOutput(w, fmt.Sprintf("Incorrect method: %v", ErrIncorrectMethod), 405, ErrIncorrectMethod)
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
		errorOutput(w, fmt.Sprintf("Incorrect query: %v", ErrIncorrectQuery), 500, ErrIncorrectQuery)
		return
	}

	result, err = Parser(decryptData.Expression)
	if err != nil {
		errorOutput(w, fmt.Sprintf("Expression is not valid: %v", err), 422, err)
		return
	}

	var body []byte

	for {
		tasks, err := application.App.Configuration.Database.DB.UnloadTasks()
		if err != nil {
			return
		}

		if len(tasks) == 0 {
			break
		}

		if application.App.Configuration.Database.DB.UpdateValues() != nil {
			return
		}

		jsonData, err := json.Marshal(tasks)

		req, err := http.NewRequest("POST", "http://localhost:"+application.App.Configuration.AgentPort+"/internal/task", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("[ERROR] failed to create request: ", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("[ERROR] failed to send request:", err)
			return
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("[ERROR] failed to read answer:", err)
			return
		}

		mapId := idMap{}

		json.Unmarshal(body, &mapId)

		application.App.Configuration.Database.DB.UpdateStatus(&mapId.Map)
	}

	err = json.NewEncoder(w).Encode(
		successOutputData{
			Result: fmt.Sprint(body),
		},
	)

	if err != nil {
		errorOutput(w, fmt.Sprintf("Internal server error: %v", err), 500, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("[INFO] success")
}
