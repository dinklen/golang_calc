package orchestrator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"golang_calc/internal/calc_libs/calculator"
	"golang_calc/internal/calc_libs/errors"
	"golang_calc/internal/config"
	"golang_calc/internal/database"
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

		err error
	)

	if r.Method != "POST" {
		errorOutput(w, fmt.Sprintf("Incorrect method: %v", errors.ErrIncorrectMethod), 405, errors.ErrIncorrectMethod)
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
		errorOutput(w, fmt.Sprintf("Incorrect query: %v", errors.ErrIncorrectQuery), 500, errors.ErrIncorrectQuery)
		return
	}

	_, err = calculator.Parser(decryptData.Expression)
	if err != nil {
		errorOutput(w, fmt.Sprintf("Expression is not valid: %v", err), 422, err)
		return
	}

	var body []byte

	for {
		tasks, err := database.DataBase.UnloadTasks()
		if err != nil {
			return
		}

		if len(tasks.Exprs) == 0 {
			break
		}

		if database.DataBase.UpdateValues() != nil {
			return
		}

		jsonData, err := json.Marshal(tasks)

		req, err := http.NewRequest("POST", "http://localhost:"+config.Conf.AgentPort+"/internal/task", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("[ERROR] failed to create request: %v", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("[ERROR] failed to send request: %v", err)
			return
		}

		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("[ERROR] failed to read answer: %v", err)
			return
		}

		mapId := idMap{}

		json.Unmarshal(body, &mapId)

		database.DataBase.UpdateUsedStatus(mapId.Map)
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
