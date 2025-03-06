package errors

import "errors"

var (
	// calc errors
	ErrIncorrectInput  = errors.New("incorrect input")
	ErrDivisionByZero  = errors.New("division by zero")
	ErrEmptyExpression = errors.New("empty expression")
	ErrNoNumbers       = errors.New("no numbers")

	// server errors
	ErrIncorrectMethod = errors.New("incorrect method")
	ErrIncorrectQuery  = errors.New("incorrect query")
)
