package gqlerrors

import (
	"errors"

	"github.com/graphql-go/graphql/language/location"
)

// FormattedError contains user and machine readable, formatted error messages.
type FormattedError struct {
	Message   string                    `json:"message"`
	Locations []location.SourceLocation `json:"locations"`
	Type      string                    `json:"type,omitempty"`
}

// ErrorType describes the type of the error. For example "NOT_FOUND" might be
// an error type the server wants to communicate to the client.
type ErrorType string

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

// FormatError from a plain error type.
func FormatError(err error) FormattedError {
	var errorType string
	formattedErrorType, isFormattedErrorType := err.(FormattedErrorType)
	if isFormattedErrorType {
		errorType = formattedErrorType.ErrorType()
	}

	switch err := err.(type) {
	case FormattedError:
		return err
	case *Error:
		return FormattedError{
			Message:   err.Error(),
			Locations: err.Locations,
			Type:      errorType,
		}
	case Error:
		return FormattedError{
			Message:   err.Error(),
			Locations: err.Locations,
			Type:      errorType,
		}
	default:
		return FormattedError{
			Message:   err.Error(),
			Locations: []location.SourceLocation{},
			Type:      errorType,
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
