package respond

import "net/http"

type APIError struct {
	Status  int
	Message string
	Err     error
}

func (e *APIError) Error() string {
	return e.Message
}

func (e *APIError) Unwrap() error {
	return e.Err
}

func BadRequest(msg string, err error) *APIError {
	return &APIError{
		Status:  http.StatusBadRequest,
		Message: msg,
		Err:     err,
	}
}

func NotFound(msg string) *APIError {
	return &APIError{
		Status:  http.StatusNotFound,
		Message: msg,
	}
}
