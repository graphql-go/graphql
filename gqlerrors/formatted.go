package gqlerrors

import (
	"errors"

	"github.com/graphql-go/graphql/language/location"
)

type FormattedError struct {
	Message       string                    `json:"message"`
	Locations     []location.SourceLocation `json:"locations"`
	OriginalError error                     `json:"-"`
}

func (g FormattedError) Error() string {
	return g.Message
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
			Message:       err.Error(),
			Locations:     err.Locations,
			OriginalError: err,
		}
	case Error:
		return FormattedError{
			Message:       err.Error(),
			Locations:     err.Locations,
			OriginalError: err,
		}
	default:
		return FormattedError{
			Message:       err.Error(),
			Locations:     []location.SourceLocation{},
			OriginalError: err,
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
