package graphql

type Directive struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Args        []*Argument `json:"args"`
	OnOperation bool        `json:"onOperation"`
	OnFragment  bool        `json:"onFragment"`
	OnField     bool        `json:"onField"`
}

/**
 * Directives are used by the GraphQL runtime as a way of modifying execution
 * behavior. Type system creators will usually not create these directly.
 */
func NewDirective(config *Directive) *Directive {
	if config == nil {
		config = &Directive{}
	}
	return &Directive{
		Name:        config.Name,
		Description: config.Description,
		Args:        config.Args,
		OnOperation: config.OnOperation,
		OnFragment:  config.OnFragment,
		OnField:     config.OnField,
	}
}

/**
 * Used to conditionally include fields or fragments
 */
var IncludeDirective *Directive = NewDirective(&Directive{
	Name: "include",
	Description: "Directs the executor to include this field or fragment only when " +
		"the `if` argument is true.",
	Args: []*Argument{
		&Argument{
			PrivateName:        "if",
			Type:               NewNonNull(Boolean),
			PrivateDescription: "Included when true.",
		},
	},
	OnOperation: false,
	OnFragment:  true,
	OnField:     true,
})

/**
 * Used to conditionally skip (exclude) fields or fragments
 */
var SkipDirective *Directive = NewDirective(&Directive{
	Name: "skip",
	Description: "Directs the executor to skip this field or fragment when the `if` " +
		"argument is true.",
	Args: []*Argument{
		&Argument{
			PrivateName:        "if",
			Type:               NewNonNull(Boolean),
			PrivateDescription: "Skipped when true.",
		},
	},
	OnOperation: false,
	OnFragment:  true,
	OnField:     true,
})
