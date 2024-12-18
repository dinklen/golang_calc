package 

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

func ErrorOutput(w http.ResponseWriter, errText string, errCode int, errEvent error) {
	log.Printf("[ERROR] %v", errEvent)
	
	json.NewEncoder(w).Encode(
		FailureOutputData {
			Error: errText,
		},
	)

	w.WriteHeader(errCode)
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	var (
		decryptData InputData
		result float64
	)

	if r.Method != http.MethodPost {
		ErrorOutput(w, "Access denied", 405, errors.New("try to use method GET"))
		return
	}

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

	if err = json.NewDecoder(encrypt_data).Decode(&decrypt_data); err != nil {
		ErrOutput(w, "Internal server error", 500, err)
		return
	}
		
	result, err = Calc(decrypt_data.Expression)

	if err != nil {
		ErrorOutput(w, "Expression is not valid", 422, err)
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
