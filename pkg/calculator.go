package main

import (
    "strconv"
    "strings"
    "errors"
    "fmt"
)

func Find(item rune, arr []rune) (int, bool) {
    var index int = 0
    for _, item2 := range arr {
	if item == item2 {
            return index, true
        }
	index++
    }
    return -1, false
}

func OperationCalc(expression []string, operations []rune) (string, error) {
	var (
		outputString string = ""
		result float64 = 0
		number1 float64 = 0
		number2 float64 = 0
		err1 error = nil
		err2 error = nil
		tempIndex int = -1
	)
	
	for index := 1; index < len(expression)-1; index += 2 {
		if _, found := Find([]rune(expression[index])[0], operations); found {
			number2, err2 = strconv.ParseFloat(expression[index+1], 64)

			if expression[index-1] == "" {
				number1 = result
				err1 = nil
				expression[tempIndex] = ""
			} else {
				number1, err1 = strconv.ParseFloat(expression[index-1], 64)
			}

			if err1 != nil || err2 != nil {
				return "0", errors.New("Invalid input")
			}

			switch expression[index] {
			case "*":
				result = number1 * number2
			case "/":
				if number2 == 0 {
					return "0", errors.New("Error: division by zero")
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
		}
	}
	
	for _, subStr := range expression {
		outputString += subStr
	}

	return outputString, nil
}

func TempCalculate(expression string) (float64, error) {
    var (
	    flag bool = false
    	    errE error = nil
    )

    operators := []rune{'*', '/', '+', '-'}
    numbers := []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.'}

    start:
    symbols := []rune(expression)
    tempStrings := []string{}
    
    var tempString string
    
    for _, r := range symbols {
	if r == ' ' {
	    continue
    	} else if _, found := Find(r, operators); found {
            tempStrings = append(tempStrings, tempString)
            tempStrings = append(tempStrings, string(r))
            tempString = ""
    	} else if _, found := Find(r, numbers); found {
            tempString += string(r)
        } else {
            return 0, errors.New("Invalid input")
        }
    }

    tempStrings = append(tempStrings, tempString)
    
    if flag {
    	expression, errE = OperationCalc(tempStrings, []rune{'+', '-'})

	if errE != nil {
		return 0, errE
	}
    } else {
	expression, errE = OperationCalc(tempStrings, []rune{'*', '/'})

	if errE != nil {
		return 0, errE
	}
	flag = true
	goto start
    }
    
    result, err := strconv.ParseFloat(expression, 64)

    if err != nil {
	return 0, errors.New("Invalid input")
    }
    
    return result, nil
}

func Calc(expression string) (float64, error) {
	var (
		result float64 = 0
	 	err error
	)
	
	tempExpression := expression
	tempExpression = strings.Replace(tempExpression, " ", "", -1)
	
	st:
	if indexL, left := Find('(', []rune(tempExpression)); left {
		counter := -1

		for indexR := indexL; indexR < len([]rune(tempExpression)); indexR++ {
			if []rune(tempExpression)[indexR] == '(' {
				counter++
			} else if []rune(tempExpression)[indexR] == ')' {
				if counter == 0 {
					result, err = Calc(tempExpression[indexL+1:indexR])
					if err != nil {
						return 0, err
					}
					tempExpression = tempExpression[:indexL] + fmt.Sprintf("%f", result) + tempExpression[indexR+1:]
					goto st
				} else {
					counter--
				}
			}
		}
	} else {
		result, err = TempCalculate(tempExpression)
		if err != nil {
			return 0, err
		}
	}

	return result, nil
}