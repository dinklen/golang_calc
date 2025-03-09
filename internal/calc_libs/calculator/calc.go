package calculator

import (
	"fmt"
	"log"
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

func operationCalc(expression []string, operations []rune) (aString string, lid uint32, aErr error) {
	var lastID uint32

	var outputString string

	defer func() {
		if ex := recover(); ex != nil {
			aString = ""
			aErr = errors.ErrIncorrectInput
		}
	}()

	for index := 1; index < len(expression)-1; index += 2 {
		if _, found := find([]rune(expression[index])[0], operations); found {
			lastID = uuid.New().ID()

			database.DataBase.Insert(
				expressions.NewExpression(
					lastID,
					expression[index-1],
					expression[index+1],
					expression[index],
				),

				rune(expression[index-1][0]) != '{',
				rune(expression[index+1][0]) != '{',
			)

			strID := fmt.Sprintf("%v", lastID)

			expression[index-1] = ""
			expression[index] = ""
			expression[index+1] = "{" + strID + "}"
		}
	}

	for _, subStr := range expression {
		outputString += subStr
	}

	return outputString, lastID, nil
}

func tempCalculate(expression string) (string, uint32, error) {
	var (
		flag bool  = false
		errE error = nil

		lid uint32
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
				return "", 0, errors.ErrIncorrectInput
			}
		}

		tempStrings = append(tempStrings, tempString)

		if flag {
			expression, lid, errE = operationCalc(tempStrings, []rune{'+', '-'})

			if errE != nil {
				return "", 0, errE
			}

			break
		} else {
			expression, lid, errE = operationCalc(tempStrings, []rune{'*', '/'})

			if errE != nil {
				return "", 0, errE
			}

			flag = true
		}
	}

	return expression, lid, nil
}

func checker(expression string) error {
	for index := 0; index < len(expression); index++ {
		if rune(expression[index]) == '{' || rune(expression[index]) == '}' {
			return errors.ErrIncorrectInput
		}
	}

	return nil
}

func Parser(expression string) (string, uint32, error) {
	var ( // brackets variables
		rbCounter int  = 0
		lbCounter int  = 0
		ok        bool = false

		result string

		lid uint32

		err error
	)

	if err = checker(expression); err != nil {
		log.Printf("[ERROR] invalid input: %v", err)
		return "", 0, err
	}

	tempExpression := expression
	tempExpression = strings.Replace(tempExpression, " ", "", -1)

	if tempExpression == "" {
		return "", 0, errors.ErrEmptyExpression
	}

	for _, num := range []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'} {
		if _, found := find(num, []rune(tempExpression)); found {
			ok = true
			break
		}
	}

	if !ok {
		return "", 0, errors.ErrNoNumbers
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
		return "", 0, errors.ErrIncorrectInput
	}

start:
	if indexL, left := find('(', []rune(tempExpression)); left {
		counter := -1

		for indexR := indexL; indexR < len([]rune(tempExpression)); indexR++ {
			if []rune(tempExpression)[indexR] == '(' {
				counter++
			} else if []rune(tempExpression)[indexR] == ')' {
				if counter == 0 {
					result, lid, err = Parser(tempExpression[indexL+1 : indexR])
					if err != nil {
						return "", 0, err
					}

					tempExpression = tempExpression[:indexL] + result + tempExpression[indexR+1:]
					goto start
				} else {
					counter--
				}
			}
		}
	} else {
		result, lid, err = tempCalculate(tempExpression)
		if err != nil {
			return "", 0, err
		}
	}

	return result, lid, nil
}
