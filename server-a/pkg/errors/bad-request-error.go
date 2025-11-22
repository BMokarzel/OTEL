package errors

import (
	"fmt"
)

type BadRequestError struct {
	GenericError
}

func (e *BadRequestError) Error() string {
	return fmt.Sprintf("service unavailable error: %s", e.GenericError.Message)
}

func NewBadRequestError(message string) *BadRequestError {
	return &BadRequestError{
		GenericError{
			Message: message,
		},
	}
}

func AsBadRequestError(err error) (bool, BadRequestError) {
	e, ok := err.(*BadRequestError)
	if !ok {
		return false, BadRequestError{}
	}
	return ok, *e
}

func IsBadRequestError(err error) bool {
	_, ok := err.(*BadRequestError)
	if !ok {
		return false
	}
	return ok
}
