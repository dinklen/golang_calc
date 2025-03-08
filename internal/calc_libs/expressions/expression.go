package expressions

import "golang_calc/internal/config"

type Expression struct {
	Id            uint32  `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int     `json:"operation_time"`
}

func NewExpression(id uint32, arg1, arg2 float64, operation string) *Expression {
	var operationTime int

	if operation == "+" {
		operationTime = config.Conf.PlusTime
	} else if operation == "-" {
		operationTime = config.Conf.MinusTime
	} else if operation == "*" {
		operationTime = config.Conf.MultipTime
	} else if operation == "/" {
		operationTime = config.Conf.DivisionTime
	}

	/*
		switch operation {
		case "+":
			operationTime = config.Conf.PlusTime
		case "-":
			operationTime = config.Conf.MinusTime
		case "/":
			operationTime = config.Conf.DivisionTime
		case "*":
			operationTime = config.Conf.MultipTime
		}
	*/

	return &Expression{
		Id:            id,
		Arg1:          arg1,
		Arg2:          arg2,
		Operation:     operation,
		OperationTime: operationTime,
	}
}

type Expressions struct {
	Exprs []*Expression `json:"expressions"`
}

func NewExpressions() *Expressions {
	return &Expressions{Exprs: []*Expression{}}
}

type ExpressionInfo struct {
	Id     int     `json:"id"`
	Status string  `json:"status"`
	Result float64 `json:"result"`
}
