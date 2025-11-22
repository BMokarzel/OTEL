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
		break
	case *errors.NotFound:
		_, e := errors.AsNotFoundError(err)
		statusCode = http.StatusNotFound
		genericError = e
		break
	case *errors.UnprocessableError:
		_, e := errors.AsUnprocessableError(err)
		statusCode = http.StatusUnprocessableEntity
		genericError = e
		break
	case *errors.InternalServerError:
		_, e := errors.AsInternalServerError(err)
		statusCode = http.StatusInternalServerError
		genericError = e
		break
	default:
		statusCode = http.StatusInternalServerError
		msg := fmt.Sprintf("Internal server error. Error: %s", err)
		genericError = errors.NewInternalServerError(msg)
		break
	}

	body, _ := json.Marshal(genericError)

	w.WriteHeader(statusCode)
	w.Write(body)
}
