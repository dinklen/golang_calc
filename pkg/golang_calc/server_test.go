
package golang_calc

import (
	"testing"
	"io"
	"encoding/json"
	"net/http/httptest"
	"bytes"
)

func TestCalcHandler(t *testing.T) {
	tests := []struct{
		name string
		method string
		expression []byte
		exAnswer []byte //expected answer
		exCode int //expected return status code from server
	}{
		//success tests
		{
			name: "stage 1",
			method: "POST",
			expression: []byte(`{"expression":"1+1"}`),
			exAnswer: []byte(`{"result":"2"}`),
			exCode: 200,
		},
		{
			name: "stage 2",
			method: "POST",
			expression: []byte(`{"expression":"2+2*2"}`),
			exAnswer: []byte(`{"result":"6"}`),
			exCode: 200,
		},
		{
			name: "stage 3",
			method: "POST",
			expression: []byte(`{"expression":"56/7-2*4"}`),
			exAnswer: []byte(`{"result":"0"}`),
			exCode: 200,
		},
		{
			name: "stage 4",
			method: "POST",
			expression: []byte(`{"expression":"1/4+9/9/2"}`),
			exAnswer: []byte(`{"result":"0.75"}`),
			exCode: 200,
		},
		{
			name: "stage 5",
			method: "POST",
			expression: []byte(`{"expression":"(13-1.4)*17"}`),
			exAnswer: []byte(`{"result":"197.2"}`),
			exCode: 200,
		},
		{
			name: "stage 6",
			method: "POST",
			expression: []byte(`{"expression":"(18-9)/(54.6+7.4)*(0+0.1)"}`),
			exAnswer: []byte(`{"result":"55.8"}`),
			exCode: 200,
		},
		{
			name: "stage 7",
			method: "POST",
			expression: []byte(`{"expression":"12.5*(9.006+(12.4+0.0001)/7.7 - 7)*7.052"}`),
			exAnswer: []byte(`{"result":"318.785889"}`),
			exCode: 200,
		},
		{
			name: "stage 8",
			method: "POST",
			expression: []byte(`{"expression":"1984.985-(((985.09835+986.04)/87.32+(12-4)+7/7/9.754)*0.007)-1.00001"}`),
			exAnswer: []byte(`{"result":"1983.770166"}`),
			exCode: 200,
		},
		{
			name: "stage 9",
			method: "POST",
			expression: []byte(`{"expression":"234.0958-213487.2345987"}`),
			exAnswer: []byte(`{"result":"-213253.138799"}`),
			exCode: 200,
		},
		{
			name: "stage 10",
			method: "POST",
			expression: []byte(`{"expression":"(0.1-0.1)*(((1234-4+964)-90)/1234)/1234"}`),
			exAnswer: []byte(`{"result":"0"}`),
			exCode: 200,
		},
		
		/*
		//calc error tests
		{
			name: "calc error: no numbers",
			method: "POST",
			expression: []byte(`{"expression":"+"}`),
			exAnswer: []byte(`{"error":"Expression is not valid"}`),
			exCode: 422,
		},
		{
			name: "calc error: too many dots",
			method: "POST",
			expression: []byte(`{"expression":"1+1.1.1"}`),
			exAnswer: []byte(`{"error":"Expression is not valid"}`),
			exCode: 422,
		},
		{
			name: "calc error: incorrect expression?",
			method: "POST",
			expression: []byte(`{"expression":"15. +(2)"}`),
			exAnswer: []byte(`{"error":"Expression is not valid"}`),
			exCode: 422,
		},
		{
			name: "calc error: no closed bracket",
			method: "POST",
			expression: []byte(`{"expression":"78.5*(12.09-14/(35/6)"}`),
			exAnswer: []byte(`{"error":"Expression is not valid"}`),
			exCode: 422,
		},
		{
			name: "calc error: division by zero",
			method: "POST",
			expression: []byte(`{"expression":"1295.9030003/(49-7*7)"}`),
			exAnswer: []byte(`{"error":"Expression is not valid"}`),
			exCode: 422,
		},
		{
			name: "calc error: expression with letters",
			method: "POST",
			expression: []byte(`{"expression":"123x+5y"}`),
			exAnswer: []byte(`{"error":"Expression is not valid"}`),
			exCode: 422,
		},
		{
			name: "calc error: empty expression",
			method: "POST",
			expression: []byte(`{"expression":""}`),
			exAnswer: []byte(`{"error":"Expression is not valid"}`),
			exCode: 422,
		},

		//server error tests
		{
			name: "server error: invalid input data",
			method: "POST",
			expression: []byte(`{"expression?":"1+0"}`),
			exAnswer: []byte(`{"error":"Internal server error"}`),
			exCode: 500,
		},
		{
			name: "server error: invalid type of input data",
			method: "POST",
			expression: []byte(`[{"expression":"2+4"},{"expression":"5-9"}]`),
			exAnswer: []byte(`{"error":"Internal server error"}`),
			exCode: 500,
		},
		{
			name: "server error: invalid method",
			method: "GET",
			expression: []byte(`{"expression":"1.1+90"}`),
			exAnswer: []byte(`{"error":"Access denied"}`),
			exCode: 405,
		},
		*/
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			r := httptest.NewRequest(test.method, "localhost:8080/api/v1/calculate", bytes.NewBufferString(string(test.expression)))
			w := httptest.NewRecorder()
			
			//the stands stood up...
			CalcHandler(w, r)

			answer := w.Result()

			defer answer.Body.Close()

			data, err := io.ReadAll(answer.Body)
			if err != nil {
				t.Errorf("failed to read body: %v", err)
			}

			var stData successOutputData
			var stAnswer successOutputData

			json.Unmarshal(data, &stData)
			json.Unmarshal(test.exAnswer, &stAnswer)

			if answer.StatusCode != test.exCode || stData.Result != stAnswer.Result {
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
