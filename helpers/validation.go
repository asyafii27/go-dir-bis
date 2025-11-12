package helpers

import (
	"github.com/go-playground/validator/v10"
)

func FormatValidationError(err error) map[string]interface{} {
	errors := make(map[string][]string)
	var message string

	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, e := range ve {
			field := e.Field()
			switch e.Tag() {
			case "required":
				msg := "Kolom " + field + " harus diisi."
				errors[field] = append(errors[field], msg)
				if message == "" {
					message = msg
				}
			case "max":
				msg := "Kolom " + field + " tidak boleh lebih dari " + e.Param() + " karakter."
				errors[field] = append(errors[field], msg)
				if message == "" {
					message = msg
				}
			default:
				msg := "Kolom " + field + " invalid."
				errors[field] = append(errors[field], msg)
				if message == "" {
					message = msg
				}
			}
		}
	} else {
		message = "Invalid request."
	}

	return map[string]interface{}{
		"message": message,
		"errors":  errors,
	}
}
