package errors

import "fmt"

type UnprocessableError struct {
	GenericError
}

func (e *UnprocessableError) Error() string {
	return fmt.Sprintf("unauthorized error: %s", e.GenericError.Message)
}

func NewUnprocessableError(message string) *UnprocessableError {
	return &UnprocessableError{
		GenericError{
			Message: message,
		},
	}
}

func AsUnprocessableError(err error) (bool, UnprocessableError) {
	e, ok := err.(*UnprocessableError)
	if !ok {
		return false, UnprocessableError{}
	}
	return ok, *e
}

func IsUnprocessableError(err error) bool {
	_, ok := err.(*UnprocessableError)
	if !ok {
		return false
	}
	return ok
}
