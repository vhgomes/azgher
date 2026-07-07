package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func Validate(s any) []ValidationError {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return []ValidationError{{Field: "unknown", Message: err.Error()}}
	}

	out := make([]ValidationError, 0, len(validationErrors))
	for _, fe := range validationErrors {
		out = append(out, ValidationError{
			Field:   strings.ToLower(fe.Field()),
			Message: formatMessage(fe),
		})
	}
	return out
}

func formatMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s é obrigatório", fe.Field())
	case "email":
		return "formato de e-mail inválido"
	case "min":
		return fmt.Sprintf("%s deve ter no mínimo %s caracteres", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s deve ter no máximo %s caracteres", fe.Field(), fe.Param())
	default:
		return fmt.Sprintf("%s é inválido", fe.Field())
	}
}
