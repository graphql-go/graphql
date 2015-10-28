package graphql

type ValidationResult struct {
	IsValid bool
	Errors  []FormattedError
}

func ValidateDocument(schema Schema, ast *AstDocument) (vr ValidationResult) {
	vr.IsValid = true
	return vr
}
