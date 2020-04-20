package validator

import "fmt"

type ErrValidationFailed struct {
	Field string
}

func IsErrValidationFailed(err error) bool {
	_, ok := err.(ErrValidationFailed)
	return ok
}

func (e ErrValidationFailed) Error() string {
	return fmt.Sprintf("Incorrect value [field=%s]", e.Field)
}
