package source

const (
	name = "GraphQL"
)

type Source struct {
	body  string
	name  string
	runes []rune
}

func New(name, body string) *Source {
	return &Source{
		name:  name,
		body:  body,
		runes: []rune(body),
	}
}

func (s *Source) Name() string {
	return s.name
}

func (s *Source) Body() string {
	return s.body
}

func (s *Source) RuneAt(i int) rune {
	if i >= len(s.runes) {
		return 0
	}
	return s.runes[i]
}
