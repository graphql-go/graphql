package languageerrors

import (
	"errors"
	"fmt"

	"github.com/chris-ramon/graphql-go/language/location"
	"github.com/chris-ramon/graphql-go/language/source"
)

func Error(s *source.Source, position int, description string) error {
	location := location.GetLocation(s, position)
	err := fmt.Sprintf(`
	Syntax error %s (%s:%s) %s
	`, s.Name, location.Line, location.Column, description)
	//err := fmt.Sprintf(`
	//Syntax error %s (%s:%s) %s \n\n %s
	//`, s.Name, location.Line, location.Column, description, highlightSourceAtLocation(s, location))
	return errors.New(err)
}

//func highlightSourceAtLocation(s source.Source, location.SourceLocation) {
//line := location.Line
//prevLineNum := fmt.Sprintf("%d", (line - 1))
//lineNum := fmt.Sprintf("%d", line)
//nextLineNum := fmt.Sprintf("%d", (line + 1))
//padLen := len(nextLineNum)
//}
