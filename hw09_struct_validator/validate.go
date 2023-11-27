package hw09structvalidator

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

const (
	ValidateTag             = "validate"
	ValidateNestedStructTag = "nested"
)

func Validate(v interface{}) error {
	value := reflect.ValueOf(v)
	if value.Kind() != reflect.Struct {
		return ErrPassedArgumentIsNotStruct
	}

	validationErrors := validateStruct(value)

	if len(validationErrors) != 0 {
		return validationErrors
	}
	return nil
}

func validateStruct(value reflect.Value) ValidationErrors {
	validationErrors := ValidationErrors{}
	for i := 0; i < value.NumField(); i++ {
		structField := value.Type().Field(i)
		fieldValue := value.Field(i)

		validateTagValue := structField.Tag.Get(ValidateTag)
		if validateTagValue == "" {
			continue
		}

		if validateTagValue == ValidateNestedStructTag {
			if fieldValue.Kind() != reflect.Struct {
				continue
			}

			errs := validateStruct(fieldValue)
			if errs != nil {
				validationErrors = append(validationErrors, errs...)
			}
			continue
		}

		errs := validateField(fieldValue, structField.Type, structField.Name, validateTagValue)
		if errs != nil {
			validationErrors = append(validationErrors, errs...)
		}
	}
	return validationErrors
}

func validateField(v reflect.Value, t reflect.Type, fieldName, validateTagValue string) ValidationErrors {
	validationErrors := ValidationErrors{}

	validators, err := getValidators(validateTagValue)
	if err != nil {
		validationErrors = append(validationErrors, ValidationError{
			Field: fieldName,
			Err:   err,
		})
		return validationErrors
	}

	switch t.Kind() {
	case reflect.Int:
		validationErrors = validateInt(int(v.Int()), validators, fieldName)
	case reflect.String:
		validationErrors = validateString(v.String(), validators, fieldName)
	case reflect.Slice:
		switch t.Elem().Kind() {
		case reflect.Int:
			for i := 0; i < v.Len(); i++ {
				errs := validateInt(int(v.Index(i).Int()), validators, fieldName)
				if errs != nil {
					validationErrors = append(validationErrors, errs...)
				}
			}
		case reflect.String:
			for i := 0; i < v.Len(); i++ {
				errs := validateString(v.Index(i).String(), validators, fieldName)
				if errs != nil {
					validationErrors = append(validationErrors, errs...)
				}
			}
		}
	}
	return validationErrors
}

func validateInt(i int, validators []validator, fieldName string) ValidationErrors {
	validationErrors := ValidationErrors{}

	for _, v := range validators {
		switch v.name {
		case "min":
			minValue, err := strconv.Atoi(v.value)
			if err != nil {
				validationErrors = append(validationErrors, ValidationError{
					Field: fieldName,
					Err:   errors.Wrap(ErrIncorrectTag, err.Error()),
				})
			}
			if i < minValue {
				validationErrors = append(validationErrors, ValidationError{
					Field: fieldName,
					Err:   ErrMinRoolViolated,
				})
			}
		case "max":
			maxValue, err := strconv.Atoi(v.value)
			if err != nil {
				validationErrors = append(validationErrors, ValidationError{
					Field: fieldName,
					Err:   errors.Wrap(ErrIncorrectTag, err.Error()),
				})
			}
			if i > maxValue {
				validationErrors = append(validationErrors, ValidationError{
					Field: fieldName,
					Err:   ErrMaxRoolViolated,
				})
			}
		case "in":
			inValuesString := strings.Split(v.value, ",")
			matched := false
			for _, inValueStr := range inValuesString {
				inValue, err := strconv.Atoi(inValueStr)
				if err != nil {
					validationErrors = append(validationErrors, ValidationError{
						Field: fieldName,
						Err:   errors.Wrap(ErrIncorrectTag, err.Error()),
					})
					return validationErrors
				}

				if i == inValue {
					matched = true
					break
				}
			}
			if !matched {
				validationErrors = append(validationErrors, ValidationError{
					Field: fieldName,
					Err:   ErrInRoolViolated,
				})
			}
		}
	}

	return validationErrors
}

func validateString(s string, validators []validator, fieldName string) ValidationErrors {
	validationErrors := ValidationErrors{}

	for _, v := range validators {
		switch v.name {
		case "len":
			requiredLen, err := strconv.Atoi(v.value)
			if err != nil {
				validationErrors = append(validationErrors, ValidationError{
					Field: fieldName,
					Err:   errors.Wrap(ErrIncorrectTag, err.Error()),
				})
			}
			sLen := len(s)
			if sLen != requiredLen {
				validationErrors = append(validationErrors, ValidationError{
					Field: fieldName,
					Err:   ErrLenRoolViolated,
				})
			}
		case "regexp":
			matched, err := regexp.MatchString(v.value, s)
			if err != nil {
				validationErrors = append(validationErrors, ValidationError{
					Field: fieldName,
					Err:   ErrIncorrectRegexpExpr,
				})
			}

			if !matched {
				validationErrors = append(validationErrors, ValidationError{
					Field: fieldName,
					Err:   ErrRegexpRoolViolated,
				})
			}
		case "in":
			inStrings := strings.Split(v.value, ",")
			matched := false
			for _, in := range inStrings {
				if s == in {
					matched = true
				}
			}
			if !matched {
				validationErrors = append(validationErrors, ValidationError{
					Field: fieldName,
					Err:   ErrInRoolViolated,
				})
			}
		}
	}

	return validationErrors
}
