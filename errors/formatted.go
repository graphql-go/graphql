package graphqlerrors

import (
	"errors"
	"github.com/chris-ramon/graphql/language/location"
)

type FormattedError struct {
	Message   string
	Locations []location.SourceLocation
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
