package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var fieldMessages = map[string]string{
	"Name":   "name",
	"Email":  "email",
	"Phone":  "phone",
	"Source": "source",
	"Status": "status",
	"ID":     "id",
	"Leads":  "leads",
}

func fieldLabel(field string) string {
	if label, ok := fieldMessages[field]; ok {
		return label
	}
	return strings.ToLower(field)
}

// ParseValidationError converts go-playground/validator errors into
// a human-readable string suitable for returning to API clients.
func ParseValidationError(err error) string {
	var ve validator.ValidationErrors
	var ok bool

	if ve, ok = err.(validator.ValidationErrors); !ok {
		// Not a validation error (e.g. malformed JSON) — return as-is.
		return err.Error()
	}

	messages := make([]string, 0, len(ve))

	for _, fe := range ve {
		field := fieldLabel(fe.Field())

		var msg string
		switch fe.Tag() {
		case "required":
			msg = fmt.Sprintf("%s is required", field)
		case "email":
			msg = fmt.Sprintf("%s must be a valid email address", field)
		case "fullname":
			msg = fmt.Sprintf("%s must include both first and last name (e.g. 'Aman Gupta')", field)
		case "phone":
			msg = fmt.Sprintf("%s must be a valid phone number (e.g. +91-9876543210)", field)
		case "min":
			msg = fmt.Sprintf("%s must be at least %s characters", field, fe.Param())
		case "max":
			msg = fmt.Sprintf("%s must be at most %s characters", field, fe.Param())
		case "oneof":
			msg = fmt.Sprintf("%s must be one of: %s", field, strings.ReplaceAll(fe.Param(), " ", ", "))
		case "uuid":
			msg = fmt.Sprintf("%s must be a valid UUID", field)
		case "dive":
			msg = fmt.Sprintf("one or more items in %s are invalid", field)
		default:
			msg = fmt.Sprintf("%s failed validation ('%s')", field, fe.Tag())
		}

		messages = append(messages, msg)
	}

	return strings.Join(messages, "; ")
}
