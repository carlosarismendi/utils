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
		HTTPCode int    `json:"httpCode"`
		Key      string `json:"key"`
		Message  string `json:"message"`
		Cause    string `json:"cause,omitempty"`
	}

	m := make(map[string]interface{})

	m["error"] = &err{
		HTTPCode: c.getHTTPCode(c.key),
		Key:      c.key,
		Message:  c.message,
		Cause:    c.cause,
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

func (c *UtilError) getHTTPCode(key string) int {
	switch key {
	case ResourceAlreadyExistsError:
		return http.StatusConflict
	case ResourceNotFoundError:
		return http.StatusNotFound
	case WrongInputParameterError:
		return http.StatusUnprocessableEntity
	case GenericError:
		return http.StatusTeapot
	default:
		return http.StatusInternalServerError
	}
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
