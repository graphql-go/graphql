package gqlerrors

import (
	"encoding/json"
	"errors"

	"github.com/graphql-go/graphql/language/location"
)

// FormattedError contains user and machine readable, formatted error messages.
type FormattedError struct {
	Message    string                    `json:"message"`
	Locations  []location.SourceLocation `json:"locations"`
	Extensions ErrorExtensions           `json:"-"`
}

// MarshalJSON implements custom JSON marshaling for the `FormattedError` type
// in order to place the `ErrorExtensions` at the top level.
func (g FormattedError) MarshalJSON() ([]byte, error) {
	m := map[string]interface{}{}
	for k, v := range g.Extensions {
		m[k] = v
	}
	m["message"] = g.Message
	m["locations"] = g.Locations
	return json.Marshal(m)
}

// FormattedErrorType is an interface that can be implemented by the underlying
// error that is passed to `FormatError` (or `FormatErrors`). If the error
// implements this interface the `Type` property of `FormattedError` is set to
// this value.
type FormattedErrorType interface {
	ErrorType() string
}

// Error implements the `error` interface.
func (g FormattedError) Error() string {
	return g.Message
}

// NewFormattedError creates a new formatted error from a string.
func NewFormattedError(message string) FormattedError {
	err := errors.New(message)
	return FormatError(err)
}

// NewFormattedErrorWithExtensions creates a new formatted error from a string
// with the given extensions.
func NewFormattedErrorWithExtensions(message string,
	extensions ErrorExtensions,
) FormattedError {
	err := FormatError(errors.New(message))
	err.Extensions = extensions
	return err
}

// FormatError from a plain error type.
func FormatError(err error) FormattedError {
	switch err := err.(type) {
	case FormattedError:
		return err
	case *Error:
		return FormattedError{
			Message:    err.Error(),
			Locations:  err.Locations,
			Extensions: err.Extensions,
		}
	case Error:
		return FormattedError{
			Message:    err.Error(),
			Locations:  err.Locations,
			Extensions: err.Extensions,
		}
	default:
		return FormattedError{
			Message:   err.Error(),
			Locations: []location.SourceLocation{},
		}
	}
}

// FormatErrors creates an array of `FormattedError`s from plain errors.
func FormatErrors(errs ...error) []FormattedError {
	formattedErrors := []FormattedError{}
	for _, err := range errs {
		formattedErrors = append(formattedErrors, FormatError(err))
	}
	return formattedErrors
}
