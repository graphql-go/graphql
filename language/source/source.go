package source

type Source struct {
	Body string
	Name string
}

func NewSource(body, name string) *Source {
	if name == "" {
		name = "GraphQL"
	}
	return &Source{
		Body: body,
		Name: name,
	}
}
