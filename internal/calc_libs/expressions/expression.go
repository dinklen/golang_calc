package expressions

import "golang_calc/internal/config"

type Expression struct {
	Id            uint32 `json:"id"`
	Arg1          string `json:"arg1"`
	Arg2          string `json:"arg2"`
	Operation     string `json:"operation"`
	OperationTime int    `json:"operation_time"`
}

func NewExpression(id uint32, arg1, arg2, operation string) *Expression {
	var operationTime int

	// знаю, но switch/case почему-то не прошёл + спешил
	if operation == "+" {
		operationTime = config.Conf.PlusTime
	} else if operation == "-" {
		operationTime = config.Conf.MinusTime
	} else if operation == "*" {
		operationTime = config.Conf.MultipTime
	} else if operation == "/" {
		operationTime = config.Conf.DivisionTime
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
	Exprs []*Expression `json:"expressions"`
}

func NewExpressions() *Expressions {
	return &Expressions{Exprs: []*Expression{}}
}

type ExpressionInfo struct {
	Id     uint32  `json:"id"`
	Status string  `json:"status"`
	Result float64 `json:"result"`
}
