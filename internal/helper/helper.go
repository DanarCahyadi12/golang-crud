package helper

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

func GetFirstValidationErrorAndConvert(validationError error) string {
	var field string
	var tag string
	var param string

	var message string

	if errs, ok := validationError.(validator.ValidationErrors); ok {
		for _, err := range errs {
			field = err.Field()
			tag = err.Tag()
			param = err.Param()
			break
		}
	}

	switch tag {
	case "email":
		message = "Email format is invalid"
		break
	case "min":
		message = fmt.Sprintf("%s %s %s character", field, tag, param)
		break
	default:
		if param == "" {
			message = fmt.Sprintf("%s %s", field, tag)
		} else {
			message = fmt.Sprintf("%s %s %s", field, tag, param)
		}

		break
	}

	return message
}
