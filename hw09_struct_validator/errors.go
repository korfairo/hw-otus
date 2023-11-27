package hw09structvalidator

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrPassedArgumentIsNotStruct = errors.New("passed value is not a struct")
	ErrIncorrectTag              = errors.New("incorrect tag")

	ErrMinRoolViolated = errors.New("value less than minimum")
	ErrMaxRoolViolated = errors.New("value greater than maximum")
	ErrInRoolViolated  = errors.New("value is not in the group of possible values")

	ErrLenRoolViolated     = errors.New("value length does not match the requirement")
	ErrIncorrectRegexpExpr = errors.New("regexp expression in validator tag is incorrect")
	ErrRegexpRoolViolated  = errors.New("string doesn't matched the regexp pattern")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var sb strings.Builder
	for _, ve := range v {
		sb.WriteString(fmt.Sprintf("%s: %s\n", ve.Field, ve.Err.Error()))
	}
	return sb.String()
}
