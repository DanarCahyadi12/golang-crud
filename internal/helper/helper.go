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
	case "min", "max":
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

func FormatNextURLPagination(path string, page int, limit int, pageSize int64) string {
	if int64(page) < pageSize {
		return fmt.Sprintf("http://localhost:8080/%s?page=%d&limit=%d", path, page+1, limit)
	}

	return ""

}

func FormatPrevURLPagination(path string, page int, limit int) string {
	if page <= 1 {
		return ""
	}
	return fmt.Sprintf("http://localhost:8080/%s?page=%d&limit=%d", path, page-1, limit)
}
