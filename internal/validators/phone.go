package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func ValidatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	if phone == "" {
		return true
	}
	
	// Basic regex for phone numbers
	// Supports optional +, optional country code, and 10-12 digits
	re := regexp.MustCompile(`^(\+\d{1,3}[- ]?)?\d{10,12}$`)
	return re.MatchString(phone)
}
