package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/BMokarzel/OTEL/server-a/pkg/errors"
)

func (r *Router) ErrorHandler(w http.ResponseWriter, rr *http.Request, err error) {
	var genericError any

	var statusCode int
	switch err := err.(type) {
	case *errors.BadRequestError:
		_, e := errors.AsBadRequestError(err)
		statusCode = http.StatusBadRequest
		genericError = e
	case *errors.NotFound:
		_, e := errors.AsNotFoundError(err)
		statusCode = http.StatusNotFound
		genericError = e
	case *errors.UnprocessableError:
		_, e := errors.AsUnprocessableError(err)
		statusCode = http.StatusUnprocessableEntity
		genericError = e
	case *errors.InternalServerError:
		_, e := errors.AsInternalServerError(err)
		statusCode = http.StatusInternalServerError
		genericError = e
	default:
		statusCode = http.StatusInternalServerError
		msg := fmt.Sprintf("Internal server error. Error: %s", err)
		genericError = errors.NewInternalServerError(msg)
	}

	body, _ := json.Marshal(genericError)

	w.WriteHeader(statusCode)
	w.Write(body)
}
