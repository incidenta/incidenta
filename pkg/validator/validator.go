package validator

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func Validate(s interface{}) error {
	errOrig := validate.Struct(s)
	if errOrig == nil {
		return nil
	}
	if _, ok := errOrig.(*validator.InvalidValidationError); ok {
		return errOrig
	}
	for _, err := range errOrig.(validator.ValidationErrors) {
		name := strings.ToLower(err.StructField())
		return ErrValidationFailed{Field: name}
	}
	return nil
}
