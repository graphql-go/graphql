package types

type GraphQLDirective struct {
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Args        []*GraphQLArgument `json:"args"`
	OnOperation bool               `json:"onOperation"`
	OnFragment  bool               `json:"onFragment"`
	OnField     bool               `json:"onField"`
}

func NewGraphQLDirective(config *GraphQLDirective) *GraphQLDirective {
	if config == nil {
		config = &GraphQLDirective{}
	}
	return &GraphQLDirective{
		Name:        config.Name,
		Description: config.Description,
		Args:        config.Args,
		OnOperation: config.OnOperation,
		OnFragment:  config.OnFragment,
		OnField:     config.OnField,
	}
}

var GraphQLIncludeDirective *GraphQLDirective = NewGraphQLDirective(&GraphQLDirective{
	Name: "include",
	Description: "Directs the executor to include this field or fragment only when " +
		"the `if` argument is true.",
	Args: []*GraphQLArgument{
		&GraphQLArgument{
			Name:        "if",
			Type:        NewGraphQLNonNull(GraphQLBoolean),
			Description: "Included when true.",
		},
	},
	OnOperation: false,
	OnFragment:  true,
	OnField:     true,
})

var GraphQLSkipDirective *GraphQLDirective = NewGraphQLDirective(&GraphQLDirective{
	Name: "skip",
	Description: "Directs the executor to skip this field or fragment when the `if` " +
		"argument is true.",
	Args: []*GraphQLArgument{
		&GraphQLArgument{
			Name:        "if",
			Type:        NewGraphQLNonNull(GraphQLBoolean),
			Description: "Skipped when true.",
		},
	},
	OnOperation: false,
	OnFragment:  true,
	OnField:     true,
})
