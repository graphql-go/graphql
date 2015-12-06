package location

import (
	"regexp"

	"github.com/graphql-go/graphql/language/source"
)

type SourceLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

func GetLocation(s *source.Source, position int) SourceLocation {
	body := ""
	if s != nil {
		body = s.Body
	}
	line := 1
	column := position + 1
	lineRegexp := regexp.MustCompile("\r\n|[\n\r\u2028\u2029]")
	matches := lineRegexp.FindAllStringIndex(body, -1)
	for _, match := range matches {
		matchIndex := match[0]
		if matchIndex < position {
			line += 1
			l := len(s.Body[match[0]:match[1]])
			column = position + 1 - (matchIndex + l)
			continue
		} else {
			break
		}
	}
	return SourceLocation{Line: line, Column: column}
}
