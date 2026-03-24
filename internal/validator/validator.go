package validator

import (
	"regexp"
	"slices"
)

var (
	// EmailRX provides a robust regular expression pattern for validating email addresses.
	// It conforms broadly to W3C standards for email input types.
	EmailRX = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
)

// Validator is a type that encapsulates a map of validation errors.
// It maps field or validation rule names to specific error messages.
type Validator struct {
	Errors map[string]string
}

// New initializes and returns a new Validator instance with an empty map for errors.
func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

// Valid determines if the validation passed successfully.
// Returns true when the map contains no error entries, meaning no rules failed.
func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

// AddError saves an error descriptor in the Validator structure for a given validation key.
// It will only retain the very first error encountered for any particular key.
func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

// Check evaluates a conditional rule and inserts an error message into the error
// map pointing to the specific validation key. The map mutation only happens if ok is false.
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

// PermittedValues checks if a given value is present within a variadic list
// of valid permitted choices using generic types that support comparison operations.
func PermittedValues[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

// Matches verifies whether the given string conforms entirely to the pattern
// specified by the compiled regular expression.
func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

// Unique inspects a typed slice and ensures that every value contained within
// appears exactly once, relying on a localized map. Generics ensure comparable capability.
func Unique[T comparable](values []T) bool {
	uniqueValues := make(map[T]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}
