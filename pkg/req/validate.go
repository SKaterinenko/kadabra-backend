package req

import (
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

func IsValid[T any](payload T) error {
	validate := validator.New()
	err := validate.Struct(payload)

	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			firstError := validationErrors[0]
			fieldName := strings.ToLower(firstError.Field())

			var message string
			switch firstError.Tag() {
			case "required":
				message = fmt.Sprintf("Field %s is required", fieldName)
			case "email":
				message = fmt.Sprintf("Field %s must be a valid email", fieldName)
			case "min":
				message = fmt.Sprintf("Field %s must be at least %s", fieldName, firstError.Param())
			case "max":
				message = fmt.Sprintf("Field %s must be at most %s", fieldName, firstError.Param())
			case "len":
				message = fmt.Sprintf("Field %s must be %s characters long", fieldName, firstError.Param())
			case "gte":
				message = fmt.Sprintf("Field %s must be greater than or equal to %s", fieldName, firstError.Param())
			case "lte":
				message = fmt.Sprintf("Field %s must be less than or equal to %s", fieldName, firstError.Param())
			case "gt":
				message = fmt.Sprintf("Field %s must be greater than %s", fieldName, firstError.Param())
			case "lt":
				message = fmt.Sprintf("Field %s must be less than %s", fieldName, firstError.Param())
			case "oneof":
				message = fmt.Sprintf("Field %s must be one of [%s]", fieldName, firstError.Param())
			default:
				message = fmt.Sprintf("Field %s validation failed on %s", fieldName, firstError.Tag())
			}

			return fmt.Errorf(message)
		}
	}

	return err
}
