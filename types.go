package graphql

// type Schema interface{}

type Result struct {
	Data   interface{}      `json:"data"`
	Errors []FormattedError `json:"errors,omitempty"`
}

func (r *Result) HasErrors() bool {
	return (len(r.Errors) > 0)
}
