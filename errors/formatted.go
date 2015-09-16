package graphqlerrors

import (
	"errors"
	"github.com/chris-ramon/graphql-go/language/location"
)

type GraphQLFormattedError struct {
	Message   string
	Locations []location.SourceLocation
}

func (g GraphQLFormattedError) Error() string {
	return g.Message
}

func NewGraphQLFormattedError(message string) GraphQLFormattedError {
	err := errors.New(message)
	return FormatError(err)
}

func FormatError(err error) GraphQLFormattedError {
	switch err := err.(type) {
	case GraphQLFormattedError:
		return err
	case GraphQLError:
		return GraphQLFormattedError{
			Message:   err.Error(),
			Locations: err.Locations,
		}
	default:
		return GraphQLFormattedError{
			Message:   err.Error(),
			Locations: []location.SourceLocation{},
		}
	}
}

func FormatErrors(errs ...error) []GraphQLFormattedError {
	formattedErrors := []GraphQLFormattedError{}
	for _, err := range errs {
		formattedErrors = append(formattedErrors, FormatError(err))
	}
	return formattedErrors
}
