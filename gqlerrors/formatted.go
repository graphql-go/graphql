package gqlerrors

import (
	"errors"
	"fmt"
	"runtime"

	"github.com/sprucehealth/graphql/language/location"
)

const (
	InternalError = "INTERNAL"
)

type FormattedError struct {
	Message     string                    `json:"message"`
	Type        string                    `json:"type,omitempty"`
	UserMessage string                    `json:"userMessage,omitempty"`
	Locations   []location.SourceLocation `json:"locations"`
	StackTrace  string                    `json:"-"`
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
	case runtime.Error:
		return FormattedError{
			Message:    err.Error(),
			Type:       InternalError,
			StackTrace: stackTrace(),
		}
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

func FormatPanic(r interface{}) FormattedError {
	if e, ok := r.(error); ok {
		return FormatError(e)
	}
	return FormattedError{
		Message:    fmt.Sprintf("panic %v", r),
		Type:       InternalError,
		StackTrace: stackTrace(),
	}
}

func FormatErrors(errs ...error) []FormattedError {
	formattedErrors := []FormattedError{}
	for _, err := range errs {
		formattedErrors = append(formattedErrors, FormatError(err))
	}
	return formattedErrors
}

func stackTrace() string {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}
