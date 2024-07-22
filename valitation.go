package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidatorFunc is the interface for custom validation functions.
type ValidatorFunc interface {
	Validate(fl FieldLevel) bool
}

// FieldLevel provides a field-level validation context.
type FieldLevel interface {
	Field() interface{}
}

// validatorAdapter adapts the go-playground validator function to the ValidatorFunc interface.
type validatorAdapter struct {
	fn func(fl validator.FieldLevel) bool
}

func (a validatorAdapter) Validate(fl FieldLevel) bool {
	return a.fn(fl.(validator.FieldLevel))
}

// Validator encapsulates the go-playground validator instance.
type Validator struct {
	validate *validator.Validate
}

// New creates a new Validator instance.
func New() *Validator {
	return &Validator{validate: validator.New()}
}

// Register registers a custom validation function.
func (v *Validator) Register(tag string, fn ValidatorFunc) error {
	return v.validate.RegisterValidation(tag, func(fl validator.FieldLevel) bool {
		return fn.Validate(fl)
	})
}

// Rule represents a single validation rule.
type Rule struct {
	Field    string
	Tag      string
	Message  string
}

// ValidateStruct validates a struct based on custom rules and returns custom validation errors.
func (v *Validator) ValidateStruct(s interface{}, rules []Rule) Errors {
	var errs Errors
	val := reflect.ValueOf(s)
	for _, rule := range rules {
		fieldVal := getField(val, rule.Field)
		if !fieldVal.IsValid() {
			continue
		}
		if err := v.validate.Var(fieldVal.Interface(), rule.Tag); err != nil {
			fieldName := rule.Field
			if jsonTag := getJSONTag(val.Type(), rule.Field); jsonTag != "" {
				fieldName = jsonTag
			}
			errs = append(errs, Error{
				Field:   fieldName,
				Message: rule.Message,
			})
		}
	}
	return errs
}

// getField returns the value of a field by name.
func getField(val reflect.Value, name string) reflect.Value {
	return val.FieldByName(name)
}

// getJSONTag returns the json tag value for a given field.
func getJSONTag(t reflect.Type, fieldName string) string {
	field, found := t.Elem().FieldByName(fieldName)
	if !found {
		return ""
	}
	tag := field.Tag.Get("json")
	if tag == "" || tag == "-" {
		return ""
	}
	return strings.Split(tag, ",")[0]
}

// Error represents a single validation error.
type Error struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}

// Errors is a collection of validation errors.
type Errors []Error

func (e Errors) Error() string {
	return fmt.Sprintf("validation errors: %v", []Error(e))
}

// SetMessage allows setting custom messages for all validation errors.
func (e *Errors) SetMessage(field, message string) {
	for _, err := range *e {
		if err.Field == field {
			err.Message = message
		}
	}
}

