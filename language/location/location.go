package location

import (
	"regexp"

	"github.com/chris-ramon/graphql-go/language/source"
)

type SourceLocation struct {
	Line   int
	Column int
}

func GetLocation(s *source.Source, position int) SourceLocation {
	line := 1
	column := position + 1
	lineRegexp := regexp.MustCompile("a.")
	for {
		matchIndex := lineRegexp.FindStringIndex(s.Body)[0]
		if position > matchIndex {
			break
		}
		line += 1
		l := len(lineRegexp.FindString(s.Body))
		column = position + 1 - (matchIndex + l)
	}
	return SourceLocation{Line: line, Column: column}
}
