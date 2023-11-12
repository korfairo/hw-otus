package hw09structvalidator

import "strings"

const ValidatorsSeparator = "|"

type validator struct {
	name  string
	value string
}

func getValidators(tagValue string) (validators []validator, err error) {
	tagValue = strings.TrimSpace(tagValue)

	validatorStrings := strings.Split(tagValue, ValidatorsSeparator)
	for _, validatorString := range validatorStrings {
		keyValue := strings.Split(validatorString, ":")
		if len(keyValue) != 2 {
			return nil, ErrIncorrectTag
		}

		validators = append(validators, validator{
			name:  keyValue[0],
			value: keyValue[1],
		})
	}
	return validators, nil
}
