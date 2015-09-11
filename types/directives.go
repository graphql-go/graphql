package types

type GraphQLDirective struct {
	Name        string
	Description string
	Args        []*GraphQLArgument
	OnOperation bool
	OnFragment  bool
	OnField     bool
}

func NewGraphQLDirective(config *GraphQLDirective) *GraphQLDirective {
	if config == nil {
		config = &GraphQLDirective{}
	}
	return &GraphQLDirective{
		Name: config.Name,
		Description: config.Description,
		Args: config.Args,
		OnOperation: config.OnOperation,
		OnFragment: config.OnFragment,
		OnField: config.OnField,
	}
}

var GraphQLIncludeDirective *GraphQLDirective = NewGraphQLDirective(&GraphQLDirective{
	Name: "include",
	Description: "Directs the executor to include this field or fragment only when " +
	"the `if` argument is true.",
	Args: []*GraphQLArgument{
		&GraphQLArgument{
			Name: "if",
			Description: "Included when true.",
		},
	},
})