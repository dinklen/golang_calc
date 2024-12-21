package golang_calc

import (
    "strconv"
    "strings"
    "errors"
    "fmt"
)

//errors initialization
var (
	ErrIncorrectInput = errors.New("incorrect input")
	ErrDivisionByZero = errors.New("division by zero")
	ErrEmptyExpression = errors.New("empty expression")
	ErrNoNumbers = errors.New("no numbers")
)

func find(item rune, arr []rune) (int, bool) {
    var index int = 0
    for _, item2 := range arr {
	if item == item2 {
            return index, true
        }
	index++
    }
    return -1, false
}

func operationCalc(expression []string, operations []rune) (string, error) {
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
		if _, found := find([]rune(expression[index])[0], operations); found {
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
		}
	}
	
	for _, subStr := range expression {
		outputString += subStr
	}

	return outputString, nil
}

func tempCalculate(expression string) (float64, error) {
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
	} else if _, found := find(r, operators); found {
            	tempStrings = append(tempStrings, tempString)
            	tempStrings = append(tempStrings, string(r))
            	tempString = ""
    	} else if _, found := find(r, numbers); found {
            	tempString += string(r)
        } else {
            	return 0, ErrIncorrectInput
        }
    }

    tempStrings = append(tempStrings, tempString)
    
    if flag {
    	expression, errE = operationCalc(tempStrings, []rune{'+', '-'})

	if errE != nil {
		return 0, errE
	}
    } else {
	expression, errE = operationCalc(tempStrings, []rune{'*', '/'})

	if errE != nil {
		return 0, errE
	}
	flag = true
	goto start
    }
    
    result, err := strconv.ParseFloat(expression, 64)

    if err != nil {
	return 0, ErrIncorrectInput
    }
    
    return result, nil
}

func Calc(expression string) (float64, error) {
	var (
		result float64 = 0
		rbCounter int = 0 //rigth brackets counter
		lbCounter int = 0 //left brackets counter
		ok bool = false
	 	err error
	)
	
	tempExpression := expression
	tempExpression = strings.Replace(tempExpression, " ", "", -1)

	//checking...

	if tempExpression == "" {return 0, ErrEmptyExpression}

	for _, num := range []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'} {
		if _, found := find(num, []rune(tempExpression)); found {
			ok = true
			break
		}
	}

	if !ok {return 0, ErrNoNumbers}
	
	for _, br := range tempExpression {
		switch br {
		case '(': lbCounter++
		case ')': rbCounter++
		}
	}
	
	if lbCounter != rbCounter {return 0, ErrIncorrectInput}
	
	st:
	if indexL, left := find('(', []rune(tempExpression)); left {
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
		result, err = tempCalculate(tempExpression)
		if err != nil {
			return 0, err
		}
	}

	return result, nil
}
