package v1

import "fmt"

type ServerError struct {
	Message     string `json:"error"`
	Description string `json:"error_description"`
}

func IsServerError(err error) bool {
	_, ok := err.(ServerError)
	return ok
}

func (s ServerError) Error() string {
	if len(s.Description) == 0 {
		return s.Message
	}
	return fmt.Sprintf("%s: %s", s.Message, s.Description)
}
