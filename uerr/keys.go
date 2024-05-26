package uerr

import "net/http"

const (
	GenericError               = "Error"
	ResourceAlreadyExistsError = "ResourceAlreadyExistsError"
	ResourceNotFoundError      = "ResourceNotFoundError"
	WrongInputParameterError   = "WrongInputParameterError"
	UnauthorizedError          = "UnauthorizedError"
	ForbiddenError             = "ForbiddenError"
)

var httpCodes = map[string]int{
	GenericError:               http.StatusInternalServerError,
	ResourceAlreadyExistsError: http.StatusConflict,
	ResourceNotFoundError:      http.StatusNotFound,
	WrongInputParameterError:   http.StatusUnprocessableEntity,
	UnauthorizedError:          http.StatusUnauthorized,
	ForbiddenError:             http.StatusForbidden,
}

// HTTPCode maps the error keys to an HTTP status code.
func HTTPCode(err error) int {
	key := GetKey(err)
	if key == "" {
		return http.StatusInternalServerError
	}

	code, ok := httpCodes[key]
	if !ok {
		return http.StatusInternalServerError
	}

	return code
}

// IsResourceNotFound returns true if the error is a ResourceNotFoundError.
func IsResourceNotFound(err error) bool {
	return Is(err, ResourceNotFoundError)
}

// IsResourceAlreadyExists returns true if the error is a IsResourceAlreadyExists.
func IsResourceAlreadyExists(err error) bool {
	return Is(err, ResourceAlreadyExistsError)
}

// IsWrongInputParameter returns true if the error is a WrongInputParameterError.
func IsWrongInputParameter(err error) bool {
	return Is(err, WrongInputParameterError)
}

// IsUnauthorized returns true if the error is a UnauthorizedError.
func IsUnauthorized(err error) bool {
	return Is(err, UnauthorizedError)
}

// IsForbidden returns true if the error is a ForbiddenError.
func IsForbidden(err error) bool {
	return Is(err, ForbiddenError)
}

func Is(err error, key string) bool {
	return GetKey(err) == key
}
