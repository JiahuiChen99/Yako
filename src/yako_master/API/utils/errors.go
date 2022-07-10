package utils

import "net/http"

// RestError models YakoAPI error response
type RestError struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Error   string `json:"error"`
}

// InternalServerError something went wrong server-side
func InternalServerError(errorMsg string) *RestError {
	return &RestError{
		Message: errorMsg,
		Status:  http.StatusInternalServerError,
		Error:   "internal server error",
	}
}

// BadRequestError the server cannot process the request
func BadRequestError(errorMsg string) *RestError {
	return &RestError{
		Message: errorMsg,
		Status:  http.StatusUnprocessableEntity,
		Error:   "bad request",
	}
}
