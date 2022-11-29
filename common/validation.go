package common

import (
	"regexp"

	"github.com/go-playground/validator"
)

// Go Validator Custom Func
func ValidateStringWhitespace(fl validator.FieldLevel) bool {
	r := regexp.MustCompile("^[a-zA-Z0-9_]+( [a-zA-Z0-9_]+)*$")
	return r.MatchString(fl.Field().String())
}

func ValidateUsername(fl validator.FieldLevel) bool {
	r := regexp.MustCompile("^[A-Za-z0-9_.]+$")
	return r.MatchString(fl.Field().String())
}
