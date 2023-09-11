package uerr

import (
	"encoding/json"
	"net/http"
)

type UError struct {
	key     string
	message string
	cause   string
}

func NewError(key, message string) *UError {
	return &UError{
		key:     key,
		message: message,
		cause:   "",
	}
}

func (c *UError) MarshalJSON() ([]byte, error) {
	type err struct {
		Key     string `json:"key"`
		Message string `json:"message"`
		Cause   string `json:"cause,omitempty"`
	}

	m := make(map[string]interface{})

	m["error"] = &err{
		Key:     c.key,
		Message: c.message,
		Cause:   c.cause,
	}

	return json.Marshal(m)
}

func (c *UError) String() string {
	b, _ := json.Marshal(c)
	return string(b)
}

func (c *UError) Error() string {
	return c.String()
}

func (c *UError) WithCause(err error) *UError {
	c.cause = err.Error()
	return c
}

func GetKey(err error) string {
	if uErr, ok := err.(*UError); ok {
		return uErr.key
	}

	return ""
}

func GetMessage(err error) string {
	if uErr, ok := err.(*UError); ok {
		return uErr.message
	}

	return err.Error()
}

func IsResourceNotFound(err error) bool {
	return GetKey(err) == ResourceNotFoundError
}

func HTTPCode(err error) int {
	key := GetKey(err)
	switch key {
	case WrongInputParameterError:
		return http.StatusUnprocessableEntity

	case ResourceNotFoundError:
		return http.StatusNotFound

	case UnauthorizedError:
		return http.StatusUnauthorized

	case ForbiddenError:
		return http.StatusForbidden

	case ResourceAlreadyExistsError:
		return http.StatusConflict

	case GenericError:
		return http.StatusTeapot

	default:
		return http.StatusInternalServerError
	}
}
