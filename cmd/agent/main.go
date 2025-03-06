package main

import (
	"log"

	"golang_calc/internal/application"
)

/*
var (
    outputString string = ""
    result float64 = 0
    number1 float64 = 0
    number2 float64 = 0
    err1 error = nil
    err2 error = nil
    tempIndex int = -1
)

=========================================================

number2, err2 = strconv.ParseFloat(expression[index+1], 64)

if index > 1 && expression[index-1] == "" {
    number1 = result
    err1 = nil
    expression[tempIndex] = ""
} else {
    number1, err1 = strconv.ParseFloat(expression[index-1], 64)
}

if err1 != nil || err2 != nil {
    return "0", ErrIncorrectInput
}

switch expression[index] {
case "*":
    result = number1 * number2
case "/":
    if number2 == 0 {
        return "0", ErrDivisionByZero
    }
    result = number1 / number2
case "+":
    result = number1 + number2
case "-":
    result = number1 - number2
}

expression[index-1] = fmt.Sprintf("%f", result)
expression[index] = ""
expression[index+1] = ""

tempIndex = index-1
*/

func main() {
	if err := application.App.RunAgent(); err != nil {
		log.Fatal("failed to start agent: ", err)
	}
}
