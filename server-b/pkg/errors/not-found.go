package errors

import (
	"fmt"
)

type NotFound struct {
	GenericError
}

func (e *NotFound) Error() string {
	return fmt.Sprintf("internal server error: %s", e.GenericError.Message)
}

func NewNotFoundError(message string) *NotFound {
	return &NotFound{
		GenericError{
			Message: message,
		},
	}
}

func AsNotFoundError(err error) (bool, NotFound) {
	e, ok := err.(*NotFound)
	if !ok {
		return false, NotFound{}
	}
	return ok, *e
}

func IsNotFoundError(err error) bool {
	_, ok := err.(*NotFound)
	if !ok {
		return false
	}
	return ok
}
