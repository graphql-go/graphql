package languageerrors

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/chris-ramon/graphql-go/errors"
	"github.com/chris-ramon/graphql-go/language/location"
	"github.com/chris-ramon/graphql-go/language/source"
)

func Error(s *source.Source, position int, description string) graphqlerrors.GraphQLError {
	l := location.GetLocation(s, position)
	err := errors.New(fmt.Sprintf("Syntax Error %s (%d:%d) %s\n\n%s", s.Name, l.Line, l.Column, description, highlightSourceAtLocation(s, l)))
	return graphqlerrors.GraphQLError{
		Error:     err,
		Source:    s,
		Positions: []int{position},
	}
}

func highlightSourceAtLocation(s *source.Source, l location.SourceLocation) string {
	line := l.Line
	prevLineNum := fmt.Sprintf("%d", (line - 1))
	lineNum := fmt.Sprintf("%d", line)
	nextLineNum := fmt.Sprintf("%d", (line + 1))
	padLen := len(nextLineNum)
	lines := regexp.MustCompile("\r\n|[\n\r\u2028\u2029]").Split(s.Body, -1)
	var highlight string
	if line >= 2 {
		highlight += fmt.Sprintf("%s: %s\n", lpad(padLen, prevLineNum), lines[line-2])
	}
	highlight += fmt.Sprintf("%s: %s\n", lpad(padLen, lineNum), lines[line-1])
	for i := 1; i < (2 + padLen + l.Column); i++ {
		highlight += " "
	}
	highlight += "^\n"
	if line < len(lines) {
		highlight += fmt.Sprintf("%s: %s\n", lpad(padLen, nextLineNum), lines[line])
	}
	return highlight
}

func lpad(l int, s string) string {
	var r string
	for i := 1; i < (l - len(s) + 1); i++ {
		r += " "
	}
	return r + s
}
