package calculator

import (
	"log"
	"strconv"
	"strings"

	"golang_calc/internal/calc_libs/errors"
	"golang_calc/internal/calc_libs/expressions"
	"golang_calc/internal/database"

	"github.com/google/uuid"
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

func operationCalc(expression []string, operations []rune) (aString string, aErr error) {
	var outputString string

	defer func() {
		if ex := recover(); ex != nil {
			aString = ""
			aErr = errors.ErrIncorrectInput
		}
	}()

	for index := 1; index < len(expression)-1; index += 2 {
		if _, found := find([]rune(expression[index])[0], operations); found {
			id := uuid.New().ID()

			arg1, err1 := strconv.ParseFloat(expression[index-1], 64)
			arg2, err2 := strconv.ParseFloat(expression[index+1], 64)

			if err1 != nil || err2 != nil {
				log.Printf("[ERROR] failed to convert arg1/arg2 to float64")
				return "", err2
			}

			database.DataBase.Insert(
				expressions.NewExpression(
					id,
					arg1,
					arg2,
					expression[index],
				),

				'{' == rune(expression[index-1][0]),
				'{' == rune(expression[index+1][0]),
			)

			strID := strconv.Itoa(int(id))

			expression[index-1] = ""
			expression[index] = ""
			expression[index+1] = "{" + strID + "}"
		}
	}

	for _, subStr := range expression {
		outputString += subStr
	}

	return outputString, nil
}

func tempCalculate(expression string) (string, error) {
	var (
		flag bool  = false
		errE error = nil

		//result string

		//err error
	)

	operators := []rune{'*', '/', '+', '-'}
	validSymbols := []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.', '{', '}', 'a', 'b', 'c', 'd', 'e', 'f'}

	for {
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
			} else if _, found := find(r, validSymbols); found {
				tempString += string(r)
			} else {
				return "", errors.ErrIncorrectInput
			}
		}

		tempStrings = append(tempStrings, tempString)

		if flag {
			expression, errE = operationCalc(tempStrings, []rune{'+', '-'})

			if errE != nil {
				return "", errE
			}

			break
		} else {
			expression, errE = operationCalc(tempStrings, []rune{'*', '/'})

			if errE != nil {
				return "", errE
			}

			flag = true
		}
	}

	/*
		if err != nil {
			return "", errors.ErrIncorrectInput
		}
	*/

	return expression, nil
}

func Parser(expression string) (string, error) {
	var ( // brackets variables
		rbCounter int  = 0
		lbCounter int  = 0
		ok        bool = false

		result string

		err error
	)

	tempExpression := expression
	tempExpression = strings.Replace(tempExpression, " ", "", -1)

	// checking...
	if tempExpression == "" {
		return "", errors.ErrEmptyExpression
	}

	for _, num := range []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'} {
		if _, found := find(num, []rune(tempExpression)); found {
			ok = true
			break
		}
	}

	if !ok {
		return "", errors.ErrNoNumbers
	}

	for _, br := range tempExpression {
		switch br {
		case '(':
			lbCounter++
		case ')':
			rbCounter++
		}
	}

	if lbCounter != rbCounter {
		return "", errors.ErrIncorrectInput
	}

start:
	if indexL, left := find('(', []rune(tempExpression)); left {
		counter := -1

		for indexR := indexL; indexR < len([]rune(tempExpression)); indexR++ {
			if []rune(tempExpression)[indexR] == '(' {
				counter++
			} else if []rune(tempExpression)[indexR] == ')' {
				if counter == 0 {
					result, err = Parser(tempExpression[indexL+1 : indexR])
					if err != nil {
						return "", err
					}

					tempExpression = tempExpression[:indexL] + result + tempExpression[indexR+1:]
					goto start
				} else {
					counter--
				}
			}
		}
	} else {
		result, err = tempCalculate(tempExpression)
		if err != nil {
			return "", err
		}
	}

	return result, nil
}
