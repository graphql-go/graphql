package gqlerrors

import (
	"github.com/graphql-go/graphql/language/location"
	"github.com/pkg/errors"
)

type FormattedError struct {
	Message   string                    `json:"message"`
	Locations []location.SourceLocation `json:"locations"`
	cause     error
}

func (g FormattedError) Error() string {
	return g.Message
}

func (g FormattedError) Cause() error {
	return g.cause
}

func NewFormattedError(message string) FormattedError {
	err := errors.New(message)
	return FormatError(err)
}

func FormatError(err error) FormattedError {
	switch err := err.(type) {
	case FormattedError:
		return err
	case *Error:
		return FormattedError{
			Message:   err.Error(),
			Locations: err.Locations,
		}
	case Error:
		return FormattedError{
			Message:   err.Error(),
			Locations: err.Locations,
		}
	default:
		return FormattedError{
			Message:   err.Error(),
			Locations: []location.SourceLocation{},
			cause:     err,
		}
	}
}

func FormatErrors(errs ...error) []FormattedError {
	formattedErrors := []FormattedError{}
	for _, err := range errs {
		formattedErrors = append(formattedErrors, FormatError(err))
	}
	return formattedErrors
}
