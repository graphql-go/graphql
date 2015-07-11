package errors

type GraphQLFormattedError struct {
	Message   string
	Locations []struct {
		Line   int
		Column int
	}
}
