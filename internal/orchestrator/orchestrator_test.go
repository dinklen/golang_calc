package orchestrator

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCalcHandler(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		expression string
		exAnswer   string //expected answer
		exCode     int    //expected return status code from server
	}{
		//success tests
		{
			name:       "stage 1",
			method:     "POST",
			expression: `{"expression":"1+1"}`,
			exAnswer:   `{"result":"2"}`,
			exCode:     200,
		},
		{
			name:       "stage 2",
			method:     "POST",
			expression: `{"expression":"2+2*2"}`,
			exAnswer:   `{"result":"6"}`,
			exCode:     200,
		},
		{
			name:       "stage 3",
			method:     "POST",
			expression: `{"expression":"56/7-2*4"}`,
			exAnswer:   `{"result":"0"}`,
			exCode:     200,
		},
		{
			name:       "stage 4",
			method:     "POST",
			expression: `{"expression":"1/4+9/9/2"}`,
			exAnswer:   `{"result":"0.75"}`,
			exCode:     200,
		},
		{
			name:       "stage 5",
			method:     "POST",
			expression: `{"expression":"15. +(2)"}`,
			exAnswer:   `{"result":"17"}`,
			exCode:     200,
		},
		{
			name:       "stage 6",
			method:     "POST",
			expression: `{"expression":"(18-9)/(54.6+7.4)*(0+0.1)"}`,
			exAnswer:   `{"result":"0.014516"}`,
			exCode:     200,
		},
		{
			name:       "stage 7",
			method:     "POST",
			expression: `{"expression":"6.36*4/76.947*(((65-0.163698)+5/2)-65)/2.356415"}`,
			exAnswer:   `{"result":"0.327795"}`,
			exCode:     200,
		},
		{
			name:       "stage 8",
			method:     "POST",
			expression: `{"expression":"1984.985-(((985.09835+986.04)/87.32+(12-4)+7/7/9.754)*0.007)-1.00001"}`,
			exAnswer:   `{"result":"1983.770256"}`,
			exCode:     200,
		},
		{
			name:       "stage 9",
			method:     "POST",
			expression: `{"expression":"234.0958-213487.2345987"}`,
			exAnswer:   `{"result":"-213253.138799"}`,
			exCode:     200,
		},
		{
			name:       "stage 10",
			method:     "POST",
			expression: `{"expression":"(0.1-0.1)*(((1234-4+964)-90)/1234)/1234"}`,
			exAnswer:   `{"result":"0"}`,
			exCode:     200,
		},

		//calc error tests
		{
			name:       "calc error: no numbers",
			method:     "POST",
			expression: `{"expression":"+"}`,
			exAnswer:   `{"error":"Expression is not valid: no numbers"}`,
			exCode:     422,
		},
		{
			name:       "calc error: too many dots",
			method:     "POST",
			expression: `{"expression":"1+1.1.1"}`,
			exAnswer:   `{"error":"Expression is not valid: incorrect input"}`,
			exCode:     422,
		},
		{
			name:       "calc error: no closed bracket",
			method:     "POST",
			expression: `{"expression":"78.5*(12.09-14/(35/6)"}`,
			exAnswer:   `{"error":"Expression is not valid: incorrect input"}`,
			exCode:     422,
		},
		{
			name:       "calc error: division by zero",
			method:     "POST",
			expression: `{"expression":"1295.9030003/(49-7*7)"}`,
			exAnswer:   `{"error":"Expression is not valid: division by zero"}`,
			exCode:     422,
		},
		{
			name:       "calc error: expression with letters",
			method:     "POST",
			expression: `{"expression":"123x+5y"}`,
			exAnswer:   `{"error":"Expression is not valid: incorrect input"}`,
			exCode:     422,
		},
		{
			name:       "calc error: empty expression",
			method:     "POST",
			expression: `{"expression":""}`,
			exAnswer:   `{"error":"Expression is not valid: empty expression"}`,
			exCode:     422,
		},

		//server error tests
		{
			name:       "server error: invalid input data",
			method:     "POST",
			expression: `{"expression?":"1+0"}`,
			exAnswer:   `{"error":"Internal server error: incorrect query"}`,
			exCode:     500,
		},
		{
			name:       "server error: invalid type of input data",
			method:     "POST",
			expression: `[{"expression":"2+4"},{"expression":"5-9"}]`,
			exAnswer:   `{"error":"Internal server error: json: cannot unmarshal array into Go value of type golang_calc.inputData"}`,
			exCode:     500,
		},
		{
			name:       "server error: invalid method",
			method:     "GET",
			expression: `{"expression":"1.1+90"}`,
<<<<<<< HEAD:internal/orchestrator/orchestrator_test.go
			exAnswer:   `{"error":"Internal server error: incorrect method"}`,
			exCode:     405,
=======
			exAnswer: `{"error":"Access denied"}`,
			exCode: 405,
>>>>>>> fa3cc82888f0fb07f4b33a3260f8beb76371fca3:pkg/golang_calc/server_test.go
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest(test.method, "localhost:8080/api/v1/calculate", bytes.NewBufferString(test.expression))
			w := httptest.NewRecorder()

			//the stands stood up...
			CalcHandler(w, r)

			answer := w.Result()

			defer answer.Body.Close()

			data, err := io.ReadAll(answer.Body)
			if err != nil {
				t.Errorf("failed to read body: %v", err)
			}

			var stSucData, stSucAnswer successOutputData
			var stFailData, stFailAnswer failureOutputData
			var index int

			if !strings.Contains(test.name, "error") {
				if err = json.Unmarshal([]byte(data), &stSucData); err != nil {
					t.Errorf("failed to decrypt body: %v", err)
				}

				if err = json.Unmarshal([]byte(test.exAnswer), &stSucAnswer); err != nil {
					t.Errorf("failer to decrypt anwer: %v", err)
				}

				index = 0
			} else {
				if err = json.Unmarshal([]byte(data), &stFailData); err != nil {
					t.Errorf("failed to decrypt data: %v", err)
				}

				if err = json.Unmarshal([]byte(test.exAnswer), &stFailAnswer); err != nil {
					t.Errorf("failed to decrypt answer: %v", err)
				}

				index = 2
			}

			stDatas := []outputData{
				stSucData,
				stSucAnswer,
				stFailData,
				stFailAnswer,
			}
<<<<<<< HEAD:internal/orchestrator/orchestrator_test.go

			if answer.StatusCode != test.exCode || stDatas[index] != stDatas[index+1] {
=======
	
			if answer.StatusCode != test.exCode || stDatas[index].GetData() != stDatas[index+1].GetData() {
>>>>>>> fa3cc82888f0fb07f4b33a3260f8beb76371fca3:pkg/golang_calc/server_test.go
				t.Errorf(
					"%s;\n----- DATA -----\nmethod: %s\nexpected status: %d\nstatus: %d\nexpression: %s\nexpected answer: %s\ngot answer: %s\n----------------",
					test.name,
					test.method,
					test.exCode,
					answer.StatusCode,
					string(test.expression),
					string(test.exAnswer),
					string(data),
				)
			}
		})
	}
}
