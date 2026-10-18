package validators

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func newValidator(t *testing.T) *validator.Validate {
	t.Helper()
	v := validator.New()
	if err := v.RegisterValidation("phone", ValidatePhone); err != nil {
		t.Fatal(err)
	}
	if err := v.RegisterValidation("fullname", ValidateFullName); err != nil {
		t.Fatal(err)
	}
	return v
}

func TestValidatePhone(t *testing.T) {
	v := newValidator(t)

	valid := []string{"", "9876543210", "+919876543210", "+91 9876543210", "+1-123456789012"}
	invalid := []string{"12345", "abcdefghij", "+91 98765", "98765432101234"}

	for _, p := range valid {
		if err := v.Var(p, "phone"); err != nil {
			t.Errorf("expected %q to be valid, got %v", p, err)
		}
	}
	for _, p := range invalid {
		if err := v.Var(p, "phone"); err == nil {
			t.Errorf("expected %q to be invalid", p)
		}
	}
}

func TestValidateFullName(t *testing.T) {
	v := newValidator(t)

	valid := []string{"Bhavya Jain", "  John   Smith  ", "Mary Ann Lee"}
	invalid := []string{"", "   ", "Bhavya"}

	for _, n := range valid {
		if err := v.Var(n, "fullname"); err != nil {
			t.Errorf("expected %q to be valid, got %v", n, err)
		}
	}
	for _, n := range invalid {
		if err := v.Var(n, "fullname"); err == nil {
			t.Errorf("expected %q to be invalid", n)
		}
	}
}

func TestIsValidStatus(t *testing.T) {
	for _, s := range []string{"NEW", "CONTACTED", "QUALIFIED", "CONVERTED", "LOST"} {
		if !IsValidStatus(s) {
			t.Errorf("expected %s to be a valid status", s)
		}
	}
	for _, s := range []string{"", "new", "DONE"} {
		if IsValidStatus(s) {
			t.Errorf("expected %q to be invalid", s)
		}
	}
}
