package orchestrator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"golang_calc/internal/calc_libs/calculator"
	"golang_calc/internal/calc_libs/errors"
	"golang_calc/internal/config"
	"golang_calc/internal/database"
)

type idMap struct {
	Map map[string]uint32 `json:"map"`
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	var (
		decryptData inputData

		err error
	)

	if r.Method != "POST" {
		ErrorOutput(w, fmt.Sprintf("Incorrect method: %v", errors.ErrIncorrectMethod), 405, errors.ErrIncorrectMethod)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	defer r.Body.Close()

	database.DataBase.Clean()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		ErrorOutput(w, fmt.Sprintf("Internal server error: %v", err), 500, err)
		return
	}

	err = json.Unmarshal(data, &decryptData)
	if err != nil {
		ErrorOutput(w, fmt.Sprintf("Internal server error: %v", err), 500, err)
		return
	}

	if string(data) != "{\"expression\":\"\"}" && decryptData.Expression == "" {
		ErrorOutput(w, fmt.Sprintf("Incorrect query: %v", errors.ErrIncorrectQuery), 500, errors.ErrIncorrectQuery)
		return
	}

	_, lastID, err := calculator.Parser(decryptData.Expression)
	if err != nil {
		ErrorOutput(w, fmt.Sprintf("Expression is not valid: %v", err), 422, err)
		return
	}

	for {
		log.Printf("Stage 1")
		tasks, err := database.DataBase.UnloadTasks()
		if err != nil {
			ErrorOutput(w, fmt.Sprintf("Internal server error: %v", err), 500, err)
			return
		}

		if len(tasks.Exprs) == 0 {
			break
		}

		log.Printf("Stage 2")
		jsonData, err := json.Marshal(tasks)
		if err != nil {
			ErrorOutput(w, fmt.Sprintf("Internal server error: %v", err), 500, err)
			return
		}

		log.Printf("%s", string(jsonData))

		log.Printf("Stage 3")

		resp, err := http.Post("http://localhost:"+config.Conf.AgentPort+"/internal/task", "application/json", bytes.NewBuffer(jsonData))
		if resp.StatusCode == 422 {
			database.DataBase.Clean()

			err = database.DataBase.InsertError(lastID)
			if err != nil {
				ErrorOutput(w, fmt.Sprintf("Internal server error: %v", err), 500, err)
				return
			}

			goto out
		} else if err != nil || resp.StatusCode != 200 {
			log.Printf("sc:%d", resp.StatusCode)
			ErrorOutput(w, fmt.Sprintf("Internal server error: %v", err), 500, err)
			return
		}

		defer func() {
			resp.Body.Close()
			r.Body.Close()
		}()

		log.Printf("Stage 4")
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			ErrorOutput(w, fmt.Sprintf("Internal server error: %v", err), 500, err)
			return
		}

		mapId := idMap{}

		log.Printf("Stage 5")
		err = json.Unmarshal(body, &mapId)
		if err != nil {
			ErrorOutput(w, fmt.Sprintf("Internal server error: %v", err), 500, err)
			return
		}

		log.Printf("Stage 6")
		dbMap := make(map[float64]uint32)

		for key, value := range mapId.Map {
			num, err := strconv.ParseFloat(key, 64)
			if err != nil {
				log.Printf("[ERROR] failed to convert key of agent's map")
				ErrorOutput(w, fmt.Sprintf("Internal server error: %v", err), 500, err)
				return
			}

			dbMap[num] = value
		}

		if err = database.DataBase.UpdateUsedStatus(dbMap); err != nil {
			ErrorOutput(w, fmt.Sprintf("Internal server error: %v", err), 500, err)
			return
		}

		log.Printf("Stage 7")
		if err = database.DataBase.UpdateValues(); err != nil {
			ErrorOutput(w, fmt.Sprintf("Internal server error: %v", err), 500, err)
			return
		}
	}

	log.Printf("Stage 8")

out:
	err = json.NewEncoder(w).Encode(
		successOutputData{
			Id: lastID,
		},
	)

	if err != nil {
		ErrorOutput(w, fmt.Sprintf("Internal server error: %v", err), 500, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("[INFO] success")
}
