package validation

import (
	"reflect"
	"testing"

	"github.com/go-playground/validator/v10"
)

type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   string `json:"age"`
}

func TestValidateStruct(t *testing.T) {
	v := New()

	validUser := User{
		Name:  "JohnDoe",
		Email: "john.doe@example.com",
		Age:   "30",
	}

	// Register a custom validation function
	v.Register("is-even", validatorAdapter{fn: func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		if len(value) == 0 {
			return false
		}
		// Check if the string length is even
		return len(value)%2 == 0
	}})

	rules := []Rule{
		{
			Field:    "Name",
			Tag:      "required,alpha",
			Message:  "Name must contain only letters",
		},
		{
			Field:    "Email",
			Tag:      "required,email",
			Message:  "Email must be a valid email address",
		},
		{
			Field:    "Age",
			Tag:      "required,numeric,is-even",
			Message:  "Age must be an even integer",
		},
	}

	if errs := v.ValidateStruct(validUser, rules); errs != nil {
		t.Errorf("expected no error, got %v", errs)
	}

	invalidUser := User{
		Name:  "John123",
		Email: "invalid-email",
		Age:   "31",
	}

	expectedErrors := Errors{
		{Field: "name", Message: "Name must contain only letters"},
		{Field: "email", Message: "Email must be a valid email address"},
		{Field: "age", Message: "Age must be an even integer"},
	}

	errs := v.ValidateStruct(invalidUser, rules)
	if !reflect.DeepEqual(errs, expectedErrors) {
		t.Errorf("expected %v, got %v", expectedErrors, errs)
	}
}

func TestValidate(t *testing.T) {
	v := New()

	validUser := User{
		Name:  "JohnDoe",
		Email: "john.doe@example.com",
		Age:   "30",
	}

	// Register a custom validation function
	v.Register("is-even", validatorAdapter{fn: func(fl validator.FieldLevel) bool {
		value := fl.Field().String()
		if len(value) == 0 {
			return false
		}
		// Check if the string length is even
		return len(value)%2 == 0
	}})

	rules := []Rule{
		{
			Field:    "Name",
			Tag:      "required,alpha",
			Message:  "Name must contain only letters",
		},
		{
			Field:    "Email",
			Tag:      "required,email",
			Message:  "Email must be a valid email address",
		},
		{
			Field:    "Age",
			Tag:      "required,numeric,is-even",
			Message:  "Age must be an even integer",
		},
	}

	if errs := v.ValidateStruct(validUser, rules); errs != nil {
		t.Errorf("expected no error, got %v", errs)
	}

	invalidUser := User{
		Name:  "John123",
		Email: "invalid-email",
		Age:   "31",
	}

	expectedErrors := Errors{
		{Field: "name", Message: "Name must contain only letters"},
		{Field: "email", Message: "Email must be a valid email address"},
		{Field: "age", Message: "Age must be an even integer"},
	}

	errs := v.ValidateStruct(invalidUser, rules)
	if !reflect.DeepEqual(errs, expectedErrors) {
		t.Errorf("expected %v, got %v", expectedErrors, errs)
	}
}

