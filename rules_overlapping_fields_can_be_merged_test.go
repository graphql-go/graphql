package graphql_test

import (
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/testutil"
)

func TestValidate_OverlappingFieldsCanBeMerged_UniqueFields(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.OverlappingFieldsCanBeMergedRule, `
      fragment uniqueFields on Dog {
        name
        nickname
      }
    `)
}
func TestValidate_OverlappingFieldsCanBeMerged_IdenticalFields(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.OverlappingFieldsCanBeMergedRule, `
      fragment mergeIdenticalFields on Dog {
        name
        name
      }
    `)
}
func TestValidate_OverlappingFieldsCanBeMerged_IdenticalFieldsWithIdenticalArgs(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.OverlappingFieldsCanBeMergedRule, `
      fragment mergeIdenticalFieldsWithIdenticalArgs on Dog {
        doesKnowCommand(dogCommand: SIT)
        doesKnowCommand(dogCommand: SIT)
      }
    `)
}
func TestValidate_OverlappingFieldsCanBeMerged_IdenticalFieldsWithIdenticalDirectives(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.OverlappingFieldsCanBeMergedRule, `
      fragment mergeSameFieldsWithSameDirectives on Dog {
        name @include(if: true)
        name @include(if: true)
      }
    `)
}
func TestValidate_OverlappingFieldsCanBeMerged_DifferentArgsWithDifferentAliases(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.OverlappingFieldsCanBeMergedRule, `
      fragment differentArgsWithDifferentAliases on Dog {
        knowsSit: doesKnowCommand(dogCommand: SIT)
        knowsDown: doesKnowCommand(dogCommand: DOWN)
      }
    `)
}
func TestValidate_OverlappingFieldsCanBeMerged_DifferentDirectivesWithDifferentAliases(t *testing.T) {
	testutil.ExpectPassesRule(t, graphql.OverlappingFieldsCanBeMergedRule, `
      fragment differentDirectivesWithDifferentAliases on Dog {
        nameIfTrue: name @include(if: true)
        nameIfFalse: name @include(if: false)
      }
    `)
}
func TestValidate_OverlappingFieldsCanBeMerged_SameAliasesWithDifferentFieldTargets(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.OverlappingFieldsCanBeMergedRule, `
      fragment sameAliasesWithDifferentFieldTargets on Dog {
        fido: name
        fido: nickname
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`Fields "fido" conflict because name and nickname are different fields.`, 3, 9, 4, 9),
	})
}
func TestValidate_OverlappingFieldsCanBeMerged_AliasMakingDirectFieldAccess(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.OverlappingFieldsCanBeMergedRule, `
      fragment aliasMaskingDirectFieldAccess on Dog {
        name: nickname
        name
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`Fields "name" conflict because nickname and name are different fields.`, 3, 9, 4, 9),
	})
}
func TestValidate_OverlappingFieldsCanBeMerged_ConflictingArgs(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.OverlappingFieldsCanBeMergedRule, `
      fragment conflictingArgs on Dog {
        doesKnowCommand(dogCommand: SIT)
        doesKnowCommand(dogCommand: HEEL)
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`Fields "doesKnowCommand" conflict because they have differing arguments.`, 3, 9, 4, 9),
	})
}
func TestValidate_OverlappingFieldsCanBeMerged_ConflictingDirectives(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.OverlappingFieldsCanBeMergedRule, `
      fragment conflictingDirectiveArgs on Dog {
        name @include(if: true)
        name @skip(if: false)
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`Fields "name" conflict because they have differing directives.`, 3, 9, 4, 9),
	})
}
func TestValidate_OverlappingFieldsCanBeMerged_ConflictingDirectiveArgs(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.OverlappingFieldsCanBeMergedRule, `
      fragment conflictingDirectiveArgs on Dog {
        name @include(if: true)
        name @include(if: false)
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`Fields "name" conflict because they have differing directives.`, 3, 9, 4, 9),
	})
}
func TestValidate_OverlappingFieldsCanBeMerged_ConflictingArgsWithMatchingDirectives(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.OverlappingFieldsCanBeMergedRule, `
      fragment conflictingArgsWithMatchingDirectiveArgs on Dog {
        doesKnowCommand(dogCommand: SIT) @include(if: true)
        doesKnowCommand(dogCommand: HEEL) @include(if: true)
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`Fields "doesKnowCommand" conflict because they have differing arguments.`, 3, 9, 4, 9),
	})
}
func TestValidate_OverlappingFieldsCanBeMerged_ConflictingDirectivesWithMatchingArgs(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.OverlappingFieldsCanBeMergedRule, `
      fragment conflictingDirectiveArgsWithMatchingArgs on Dog {
        doesKnowCommand(dogCommand: SIT) @include(if: true)
        doesKnowCommand(dogCommand: SIT) @skip(if: false)
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`Fields "doesKnowCommand" conflict because they have differing directives.`, 3, 9, 4, 9),
	})
}
func TestValidate_OverlappingFieldsCanBeMerged_EncountersConflictInFragments(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.OverlappingFieldsCanBeMergedRule, `
      {
        ...A
        ...B
      }
      fragment A on Type {
        x: a
      }
      fragment B on Type {
        x: b
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`Fields "x" conflict because a and b are different fields.`, 7, 9, 10, 9),
	})
}
func TestValidate_OverlappingFieldsCanBeMerged_ReportsEachConflictOnce(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.OverlappingFieldsCanBeMergedRule, `
      {
        f1 {
          ...A
          ...B
        }
        f2 {
          ...B
          ...A
        }
        f3 {
          ...A
          ...B
          x: c
        }
      }
      fragment A on Type {
        x: a
      }
      fragment B on Type {
        x: b
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`Fields "x" conflict because a and b are different fields.`, 18, 9, 21, 9),
		testutil.RuleError(`Fields "x" conflict because a and c are different fields.`, 18, 9, 14, 11),
		testutil.RuleError(`Fields "x" conflict because b and c are different fields.`, 21, 9, 14, 11),
	})
}
func TestValidate_OverlappingFieldsCanBeMerged_DeepConflict(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.OverlappingFieldsCanBeMergedRule, `
      {
        field {
          x: a
        },
        field {
          x: b
        }
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(`Fields "field" conflict because subfields "x" conflict because a and b are different fields.`,
			3, 9, 6, 9, 4, 11, 7, 11),
	})
}
func TestValidate_OverlappingFieldsCanBeMerged_DeepConflictWithMultipleIssues(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.OverlappingFieldsCanBeMergedRule, `
      {
        field {
          x: a
          y: c
        },
        field {
          x: b
          y: d
        }
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(
			`Fields "field" conflict because subfields "x" conflict because a and b are different fields and `+
				`subfields "y" conflict because c and d are different fields.`,
			3, 9, 7, 9, 4, 11, 8, 11, 5, 11, 9, 11),
	})
}
func TestValidate_OverlappingFieldsCanBeMerged_VeryDeepConflict(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.OverlappingFieldsCanBeMergedRule, `
      {
        field {
          deepField {
            x: a
          }
        },
        field {
          deepField {
            x: b
          }
        }
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(
			`Fields "field" conflict because subfields "deepField" conflict because subfields "x" conflict because `+
				`a and b are different fields.`,
			3, 9, 8, 9, 4, 11, 9, 11, 5, 13, 10, 13),
	})
}
func TestValidate_OverlappingFieldsCanBeMerged_ReportsDeepConflictToNearestCommonAncestor(t *testing.T) {
	testutil.ExpectFailsRule(t, graphql.OverlappingFieldsCanBeMergedRule, `
      {
        field {
          deepField {
            x: a
          }
          deepField {
            x: b
          }
        },
        field {
          deepField {
            y
          }
        }
      }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(
			`Fields "deepField" conflict because subfields "x" conflict because `+
				`a and b are different fields.`,
			4, 11, 7, 11, 5, 13, 8, 13),
	})
}

var stringBoxObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "StringBox",
	Fields: graphql.Fields{
		"scalar": &graphql.Field{
			Type: graphql.String,
		},
	},
})
var intBoxObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "IntBox",
	Fields: graphql.Fields{
		"scalar": &graphql.Field{
			Type: graphql.Int,
		},
	},
})
var nonNullStringBox1Object = graphql.NewObject(graphql.ObjectConfig{
	Name: "NonNullStringBox1",
	Fields: graphql.Fields{
		"scalar": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
})
var nonNullStringBox2Object = graphql.NewObject(graphql.ObjectConfig{
	Name: "NonNullStringBox2",
	Fields: graphql.Fields{
		"scalar": &graphql.Field{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
})
var boxUnionObject = graphql.NewUnion(graphql.UnionConfig{
	Name: "BoxUnion",
	ResolveType: func(value interface{}, info graphql.ResolveInfo) *graphql.Object {
		return stringBoxObject
	},
	Types: []*graphql.Object{
		stringBoxObject,
		intBoxObject,
		nonNullStringBox1Object,
		nonNullStringBox2Object,
	},
})

var connectionObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "Connection",
	Fields: graphql.Fields{
		"edges": &graphql.Field{
			Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
				Name: "Edge",
				Fields: graphql.Fields{
					"node": &graphql.Field{
						Type: graphql.NewObject(graphql.ObjectConfig{
							Name: "Node",
							Fields: graphql.Fields{
								"id": &graphql.Field{
									Type: graphql.ID,
								},
								"name": &graphql.Field{
									Type: graphql.String,
								},
							},
						}),
					},
				},
			})),
		},
	},
})
var schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: graphql.NewObject(graphql.ObjectConfig{
		Name: "QueryRoot",
		Fields: graphql.Fields{
			"boxUnion": &graphql.Field{
				Type: boxUnionObject,
			},
			"connection": &graphql.Field{
				Type: connectionObject,
			},
		},
	}),
})

func TestValidate_OverlappingFieldsCanBeMerged_ReturnTypesMustBeUnambiguous_ConflictingScalarReturnTypes(t *testing.T) {
	testutil.ExpectFailsRuleWithSchema(t, &schema, graphql.OverlappingFieldsCanBeMergedRule, `
        {
          boxUnion {
            ...on IntBox {
              scalar
            }
            ...on StringBox {
              scalar
            }
          }
        }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(
			`Fields "scalar" conflict because they return differing types Int and String.`,
			5, 15, 8, 15),
	})
}
func TestValidate_OverlappingFieldsCanBeMerged_ReturnTypesMustBeUnambiguous_SameWrappedScalarReturnTypes(t *testing.T) {
	testutil.ExpectPassesRuleWithSchema(t, &schema, graphql.OverlappingFieldsCanBeMergedRule, `
        {
          boxUnion {
            ...on NonNullStringBox1 {
              scalar
            }
            ...on NonNullStringBox2 {
              scalar
            }
          }
        }
    `)
}
func TestValidate_OverlappingFieldsCanBeMerged_ReturnTypesMustBeUnambiguous_ComparesDeepTypesIncludingList(t *testing.T) {
	testutil.ExpectFailsRuleWithSchema(t, &schema, graphql.OverlappingFieldsCanBeMergedRule, `
        {
          connection {
            ...edgeID
            edges {
              node {
                id: name
              }
            }
          }
        }

        fragment edgeID on Connection {
          edges {
            node {
              id
            }
          }
        }
    `, []gqlerrors.FormattedError{
		testutil.RuleError(
			`Fields "edges" conflict because subfields "node" conflict because subfields "id" conflict because `+
				`id and name are different fields.`,
			14, 11, 5, 13, 15, 13, 6, 15, 16, 15, 7, 17),
	})
}
func TestValidate_OverlappingFieldsCanBeMerged_ReturnTypesMustBeUnambiguous_IgnoresUnknownTypes(t *testing.T) {
	testutil.ExpectPassesRuleWithSchema(t, &schema, graphql.OverlappingFieldsCanBeMergedRule, `
        {
          boxUnion {
            ...on UnknownType {
              scalar
            }
            ...on NonNullStringBox2 {
              scalar
            }
          }
        }
    `)
}
