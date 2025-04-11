package validator

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	Errors map[string]string
}

// We create a new Validator using a factory function
// We have used factory functions before
func NewValidator() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

func (v *Validator) ValidData() bool {
	return len(v.Errors) == 0
}

// Add an error entry to the error map.
// the field and the failed validation message to be sent to the client
func (v *Validator) AddError(field string, message string) {
	_, exists := v.Errors[field]
	if !exists {
		v.Errors[field] = message
	}
}

func (v *Validator) Check(ok bool, field string, message string) {
	if !ok {
		v.AddError(field, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxLength(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func MinLength(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func IsValidEmail(email string) bool {
	return EmailRX.MatchString(email)
}
