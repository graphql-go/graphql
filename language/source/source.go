package source

const (
	defaultName = "GraphQL"
)

type Source struct {
	Name string
	Body string
}

func NewSource(name string, body string) *Source {
	s := &Source{
		Name: name,
		Body: body,
	}
	if s.Name == "" {
		s.Name = defaultName
	}
	return s
}
