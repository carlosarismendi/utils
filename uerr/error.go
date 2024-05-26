package uerr

import (
	"encoding/json"
)

type UError struct {
	key     string
	message string
	cause   error
}

// NewError creates a new UError with given key and message.
func NewError(key, message string) *UError {
	return &UError{
		key:     key,
		message: message,
		cause:   nil,
	}
}

// FromBytes only builds Key and Message build from bytes, ignoring Cause.
func FromBytes(b []byte) (*UError, error) {
	type errDetails struct {
		Key     string `json:"key"`
		Message string `json:"message"`
	}

	type errStruct struct {
		Error *errDetails `json:"error"`
	}

	var e errStruct

	if err := json.Unmarshal(b, &e); err != nil {
		return nil, err
	}
	uerr := &UError{
		key:     e.Error.Key,
		message: e.Error.Message,
	}
	return uerr, nil
}

// nolint:lll // long line needed
// MarshalJSON returns the json representation of the error, adding a parent key "error".
// Example 1: using a fmt.Errorf(...) as error cause.
// err := NewError("myKey", "myMessage").Cause(fmt.Errorf("myCause"))
// err.MarshalJSON() => {"error":{"key":"myKey","message":"myMessage","cause":"myCause"}}
//
// Example 2: using another UError as error cause.
// err := NewError("myKey", "myMessage").Cause(NewError("myCauseKey", "myCauseMessage"))
// err.MarshalJSON() => {"error":{"key":"myKey","message":"myMessage","cause":{"error":{"key":"myKey","message":"myMessage"}}}}
func (c *UError) MarshalJSON() ([]byte, error) {
	type err struct {
		Key     string `json:"key"`
		Message string `json:"message"`
		Cause   any    `json:"cause,omitempty"`
	}

	resErr := &err{
		Key:     c.key,
		Message: c.message,
		Cause:   c.cause,
	}
	if c.cause != nil {
		if _, ok := c.cause.(*UError); !ok {
			resErr.Cause = c.cause.Error()
		}
	}

	m := make(map[string]interface{})

	m["error"] = resErr

	return json.Marshal(m)
}

// UnmarshalJSON builds a UError by calling FromBytes.
func (c *UError) UnmarshalJSON(b []byte) error {
	uerr, err := FromBytes(b)
	if err != nil {
		return err
	}
	*c = *uerr
	return nil
}

// String returns the json representation of the error by calling MarshalJSON.
func (c *UError) String() string {
	b, _ := json.Marshal(c)
	return string(b)
}

// Error returns the string representation of the error by calling String.
func (c *UError) Error() string {
	return c.String()
}

// WithCause adds a cause to the error.
func (c *UError) WithCause(err error) *UError {
	c.cause = err
	return c
}

// GetKey returns the key of the error if it is an UError or empty string if it is another error type.
func GetKey(err error) string {
	if uErr, ok := err.(*UError); ok {
		return uErr.key
	}

	return ""
}

// GetMessage returns the message of the error if it is an UError, err.Error() otherwise.
func GetMessage(err error) string {
	if uErr, ok := err.(*UError); ok {
		return uErr.message
	}

	return err.Error()
}
