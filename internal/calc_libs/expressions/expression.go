package expressions

import "golang_calc/internal/application"

type Expression struct {
	Id            int    `json:"id"`
	Arg1          string `json:"arg1"`
	Arg2          string `json:"arg2"`
	Operation     string `json:"operation"`
	OperationTime string `json:"operation_time"`
}

func NewExpression(id int, arg1, arg2, operation string) *Expression {
	var operationTime string

	switch operation {
	case "+":
		operationTime = application.App.Configuration.PlusTime
	case "-":
		operationTime = application.App.Configuration.MinusTime
	case "/":
		operationTime = application.App.Configuration.DivisionTime
	case "*":
		operationTime = application.App.Configuration.MultipTime
	}

	return &Expression{
		Id:            id,
		Arg1:          arg1,
		Arg2:          arg2,
		Operation:     operation,
		OperationTime: operationTime,
	}
}

type Expressions struct {
	Exprs []Expression `json:"expressions"`
}

func NewExpressions() *Expressions {
	return &Expressions{Exprs: []Expression{}}
}

type ExpressionInfo struct {
	Id     int     `json:"id"`
	Status string  `json:"status"`
	Result float64 `json:"result"`
}
