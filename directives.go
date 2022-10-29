package graphql

const (
	// Operations
	DirectiveLocationQuery              = "QUERY"
	DirectiveLocationMutation           = "MUTATION"
	DirectiveLocationSubscription       = "SUBSCRIPTION"
	DirectiveLocationField              = "FIELD"
	DirectiveLocationFragmentDefinition = "FRAGMENT_DEFINITION"
	DirectiveLocationFragmentSpread     = "FRAGMENT_SPREAD"
	DirectiveLocationInlineFragment     = "INLINE_FRAGMENT"

	// Schema Definitions
	DirectiveLocationSchema               = "SCHEMA"
	DirectiveLocationScalar               = "SCALAR"
	DirectiveLocationObject               = "OBJECT"
	DirectiveLocationFieldDefinition      = "FIELD_DEFINITION"
	DirectiveLocationArgumentDefinition   = "ARGUMENT_DEFINITION"
	DirectiveLocationInterface            = "INTERFACE"
	DirectiveLocationUnion                = "UNION"
	DirectiveLocationEnum                 = "ENUM"
	DirectiveLocationEnumValue            = "ENUM_VALUE"
	DirectiveLocationInputObject          = "INPUT_OBJECT"
	DirectiveLocationInputFieldDefinition = "INPUT_FIELD_DEFINITION"
)

// DefaultDeprecationReason Constant string used for default reason for a deprecation.
const DefaultDeprecationReason = "No longer supported"

// SpecifiedRules The full list of specified directives.
var SpecifiedDirectives = []*Directive{
	IncludeDirective,
	SkipDirective,
	DeprecatedDirective,
}

// Directive structs are used by the GraphQL runtime as a way of modifying execution
// behavior. Type system creators will usually not create these directly.
type Directive struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Locations   []string    `json:"locations"`
	Args        []*Argument `json:"args"`

	err error
}

// DirectiveConfig options for creating a new GraphQLDirective
type DirectiveConfig struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Locations   []string            `json:"locations"`
	Args        FieldConfigArgument `json:"args"`
}

func NewDirective(config DirectiveConfig) *Directive {
	dir := &Directive{}

	// Ensure directive is named
	if dir.err = invariant(config.Name != "", "Directive must be named."); dir.err != nil {
		return dir
	}

	// Ensure directive name is valid
	if dir.err = assertValidName(config.Name); dir.err != nil {
		return dir
	}

	// Ensure locations are provided for directive
	if dir.err = invariant(len(config.Locations) > 0, "Must provide locations for directive."); dir.err != nil {
		return dir
	}

	args := []*Argument{}

	for _, arg := range config.Args {
		if dir.err = assertValidName(arg.Name); dir.err != nil {
			return dir
		}
		args = append(args, &Argument{
			PrivateName:        arg.Name,
			PrivateDescription: arg.Description,
			Type:               arg.Type,
			DefaultValue:       arg.DefaultValue,
		})
	}

	dir.Name = config.Name
	dir.Description = config.Description
	dir.Locations = config.Locations
	dir.Args = args
	return dir
}

// IncludeDirective is used to conditionally include fields or fragments.
var IncludeDirective = NewDirective(DirectiveConfig{
	Name: "include",
	Description: "Directs the executor to include this field or fragment only when " +
		"the `if` argument is true.",
	Locations: []string{
		DirectiveLocationField,
		DirectiveLocationFragmentSpread,
		DirectiveLocationInlineFragment,
	},
	Args: FieldConfigArgument{
		&ArgumentConfig{
			Name:        "if",
			Type:        NewNonNull(Boolean),
			Description: "Included when true.",
		},
	},
})

// SkipDirective Used to conditionally skip (exclude) fields or fragments.
var SkipDirective = NewDirective(DirectiveConfig{
	Name: "skip",
	Description: "Directs the executor to skip this field or fragment when the `if` " +
		"argument is true.",
	Args: FieldConfigArgument{
		&ArgumentConfig{
			Name:        "if",
			Type:        NewNonNull(Boolean),
			Description: "Skipped when true.",
		},
	},
	Locations: []string{
		DirectiveLocationField,
		DirectiveLocationFragmentSpread,
		DirectiveLocationInlineFragment,
	},
})

// DeprecatedDirective  Used to declare element of a GraphQL schema as deprecated.
var DeprecatedDirective = NewDirective(DirectiveConfig{
	Name:        "deprecated",
	Description: "Marks an element of a GraphQL schema as no longer supported.",
	Args: FieldConfigArgument{
		&ArgumentConfig{
			Name: "reason",
			Type: String,
			Description: "Explains why this element was deprecated, usually also including a " +
				"suggestion for how to access supported similar data. Formatted" +
				"in [Markdown](https://daringfireball.net/projects/markdown/).",
			DefaultValue: DefaultDeprecationReason,
		},
	},
	Locations: []string{
		DirectiveLocationFieldDefinition,
		DirectiveLocationEnumValue,
	},
})
