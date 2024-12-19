
package golang_calc

import "net/http/httptest"

func TestCalcHandler(t *testing.T) {
	tests := []struct{
		name string
		expression []byte
		exAnswer []byte //excepted answer
		exCode int //excepted return status code from server
	}{
		//success tests
		{
			name: "stage 1",
			expression: []byte(`{"expression":"1+1"}`),
			exAnswer: "2",
			exCode: 200,
		},
		{
			name: "stage 2",
			expression: []byte(`{"expression":"2+2*2"}`),
			exAnswer: "6",
			exCode: 200,
		},
		{
			name: "stage 3",
			expresison: "56/7-2*4",
			exAnswer: "0",
			exCode: 200,
		},
		{
			name: "stage 4",
			expression: "1/4+9/9/2",
			exAnswer: "0.75",
			exCode: 200,
		},
		{
			name: "stage 5",
			expression: "(13-1.4)*17",
			exAnswer: "197.2",
			exCode: 200,
		},
		{
			name: "stage 6",
			expression: "(18-9)/(54.6+7.4)*(0+0.1)",
			exAnswer: "55.8",
			exCode: 200,
		},
		{
			name: "stage 7",
			expression: "12.5*(9.006+(12.4+0.0001)/7.7 - 7)*7.052",
			exAnswer: "318.785888961",
			exCode: 200,
		},
		{
			name: "stage 8",
			expression: "1984.985-(((985.09835+986.04)/87.32+(12-4)+7/7/9.754)*0.007)-1.00001",
			exAnswer: "1983.770166216",
			exCode: 200,
		},
		{
			name: "stage 9",
			expression: "234.0958-213487.2345987",
			exAnswer: "-213253.1387987",
			exCode: 200,
		},
		{
			name: "stage 10",
			expression: "(0.1-0.1)*(((1234-4+964)-90)/1234)/1234",
			exAnswer: "0",
			exCode: 200,
		},

		//calc error tests
		{
			name: "calc error: no numbers",
			expression: "+",
			exAnswer: "Expression is not valid",
			exCode: 422,
		},
		{
			name: "calc error: too many dots",
			expression: "1+1.1.1",
			exAnswer: "Expression is not valid",
			exCode: 422,
		},
		{
			name: "calc error: incorrect expression?",
			expression: "15. +(2)",
			exAnswer: "Expression is not valid",
			exCode: 422,
		},
		{
			name: "calc error: no closed bracket",
			expression: "78.5*(12.09-14/(35/6)",
			exAnswer: "Expression is not valid",
			exCode: 422,
		},
		{
			name: "calc error: division by zero",
			expression: "1295.9030003/(49-7*7)",
			exAnswer: "Expression is not valid",
			exCode: 422,
		},
		{
			name: "calc error: letters",
			expression: []byte(`{"expression":"123x+5y"}`),
			exAnswer: []byte(`{"error":"Expression is not valid"}`),
			exCode: 422,
		},
		{
			name: "calc error: "
		}

		//server error tests
		{
			name: "server error test 1",
			expression: []byte(`{"expression?":"1+0"}`)
			exAnswer: []byte(`{"error":"Internal server error"}`)
			exCode: 500,
		},
		{
			name: "server error: "
		}
	}


}
