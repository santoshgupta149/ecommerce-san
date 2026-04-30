// pkg/validator/validator.go
package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// FormatErrors converts Go's raw validation errors into human-readable messages
func FormatErrors(err error) []ValidationError {
	var errs []ValidationError

	// Check if it's actually a validator.ValidationErrors type
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		// e.g. malformed JSON — not a field error
		errs = append(errs, ValidationError{
			Field:   "request",
			Message: "invalid JSON body",
		})
		return errs
	}

	for _, e := range validationErrors {
		errs = append(errs, ValidationError{
			Field:   e.Field(),
			Message: buildMessage(e),
		})
	}

	return errs
}

func buildMessage(e validator.FieldError) string {
	field := e.Field()   // "Mobile"
	tag   := e.Tag()     // "min", "required", "email" ...
	param := e.Param()   // "10" for min=10

	switch tag {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", field, param)
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", field, param)
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters", field, param)
	case "oneof":
		return fmt.Sprintf("%s must be one of: %s", field, param)
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, param)
	case "gte":
		return fmt.Sprintf("%s must be %s or greater", field, param)
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, param)
	case "lte":
		return fmt.Sprintf("%s must be %s or less", field, param)
	case "alphanum":
		return fmt.Sprintf("%s must contain only letters and numbers", field)
	case "numeric":
		return fmt.Sprintf("%s must contain only numbers", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "e164":
		return fmt.Sprintf("%s must be a valid phone number (e.g. +919876543210)", field)
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}