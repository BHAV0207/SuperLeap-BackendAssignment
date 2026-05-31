package validators

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

func ValidateFullName(fl validator.FieldLevel) bool {
	name := strings.TrimSpace(fl.Field().String())
	if name == "" {
		return false
	}

	// Implementation of "middle ground": ensure at least one space
	// which usually indicates a First and Last name.
	parts := strings.Fields(name)
	return len(parts) >= 2
}
