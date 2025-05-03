package errorsStatus

import (
	"errors"
	"net/http"
)

type ErrorsStatus struct {
	StatusCode int
	Err        error
	Message    string
}

func New(statusCode int, err string, message string) *ErrorsStatus {
	return &ErrorsStatus{
		StatusCode: statusCode,
		Err:        errors.New(err),
		Message:    message,
	}
}

func (r *ErrorsStatus) Error() string {
	return r.Err.Error()
}

func StatusCode(err error) int {
	statusCode := http.StatusBadRequest
	if e, ok := err.(*ErrorsStatus); ok {
		statusCode = e.StatusCode
	}
	return statusCode
}

func Message(err error) string {
	if e, ok := err.(*ErrorsStatus); ok {
		return e.Message
	}
	return ""
}
