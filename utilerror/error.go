package utilerror

import (
	"encoding/json"
	"net/http"
)

type UtilError struct {
	key     string
	message string
	cause   string
}

func NewError(key, message string) *UtilError {
	return &UtilError{
		key:     key,
		message: message,
		cause:   "",
	}
}

func (c *UtilError) MarshalJSON() ([]byte, error) {
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

func (c *UtilError) String() string {
	b, _ := json.Marshal(c)
	return string(b)
}

func (c *UtilError) Error() string {
	return c.String()
}

func (c *UtilError) WithCause(err error) *UtilError {
	c.cause = err.Error()
	return c
}

func GetKey(err error) string {
	if uErr, ok := err.(*UtilError); ok {
		return uErr.key
	}

	return ""
}

func GetMessage(err error) string {
	if uErr, ok := err.(*UtilError); ok {
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

	case ResourceAlreadyExistsError:
		return http.StatusConflict

	case GenericError:
		return http.StatusTeapot

	default:
		return http.StatusInternalServerError
	}
}
