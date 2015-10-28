package graphql

type AstLocation struct {
	Start  int
	End    int
	Source *Source
}

func NewAstLocation(loc *AstLocation) *AstLocation {
	if loc == nil {
		loc = &AstLocation{}
	}
	return &AstLocation{
		Start:  loc.Start,
		End:    loc.End,
		Source: loc.Source,
	}
}
