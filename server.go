package main

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
	"log"
	"fmt"
)

type SuccessOutputData struct {
	Result string `json:"result"`
}

type FailureOutputData struct {
	Error string `json:"error"`
}

type InputData struct {
	Expression string `json:"expression"`
}

func Calc(abc string) (float64, error) {
	return 1, nil
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		defer r.Body.Close()

		encrypt_data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("[ERROR] %v", err)
			json.NewEncoder(w).Encode( //!
				FailureOutputData {
					Error: "Internal server error",
				},
			)

			w.WriteHeader(500)
			return
		}

		var decrypt_data InputData
		if err = json.Unmarshal(encrypt_data, &decrypt_data); err != nil {
			log.Printf("[ERROR] %v", err)
			json.NewEncoder(w).Encode(
				FailureOutputData {
					Error: "Internal server error",
				},
			)

			w.WriteHeader(500)
			return
		}
		
		var result float64
		result, err = Calc(decrypt_data.Expression)

		if err != nil {
			log.Printf("[ERROR] %v", err)

			json.NewEncoder(w).Encode(
				FailureOutputData {
					Error: "Expression is not valid",
				},
			)

			w.WriteHeader(422)
			return
		}

		json.NewEncoder(w).Encode(
			SuccessOutputData {
				Result: fmt.Sprintf("%.5f", result),
			},
		)

		w.WriteHeader(200)
		log.Printf("[INFO] success")
	}
}

func StartServer() {
	http.HandleFunc("/api/v1/calculate", CalcHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Printf("[ERROR] %v", err)
	}
}

func main() {
	StartServer()
}
