package rules_test

import (
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/gqlerrors"
)

func TestValidate_PossibleFragmentSpreads_OfTheSameObject(t *testing.T) {
	expectPassesRule(t, graphql.PossibleFragmentSpreadsRule, `
      fragment objectWithinObject on Dog { ...dogFragment }
      fragment dogFragment on Dog { barkVolume }
    `)
}
func TestValidate_PossibleFragmentSpreads_OfTheSameObjectWithInlineFragment(t *testing.T) {
	expectPassesRule(t, graphql.PossibleFragmentSpreadsRule, `
      fragment objectWithinObjectAnon on Dog { ... on Dog { barkVolume } }
    `)
}
func TestValidate_PossibleFragmentSpreads_ObjectIntoAnImplementedInterface(t *testing.T) {
	expectPassesRule(t, graphql.PossibleFragmentSpreadsRule, `
      fragment objectWithinInterface on Pet { ...dogFragment }
      fragment dogFragment on Dog { barkVolume }
    `)
}
func TestValidate_PossibleFragmentSpreads_ObjectIntoContainingUnion(t *testing.T) {
	expectPassesRule(t, graphql.PossibleFragmentSpreadsRule, `
      fragment objectWithinUnion on CatOrDog { ...dogFragment }
      fragment dogFragment on Dog { barkVolume }
    `)
}
func TestValidate_PossibleFragmentSpreads_UnionIntoContainedObject(t *testing.T) {
	expectPassesRule(t, graphql.PossibleFragmentSpreadsRule, `
      fragment unionWithinObject on Dog { ...catOrDogFragment }
      fragment catOrDogFragment on CatOrDog { __typename }
    `)
}
func TestValidate_PossibleFragmentSpreads_UnionIntoOverlappingInterface(t *testing.T) {
	expectPassesRule(t, graphql.PossibleFragmentSpreadsRule, `
      fragment unionWithinInterface on Pet { ...catOrDogFragment }
      fragment catOrDogFragment on CatOrDog { __typename }
    `)
}
func TestValidate_PossibleFragmentSpreads_UnionIntoOverlappingUnion(t *testing.T) {
	expectPassesRule(t, graphql.PossibleFragmentSpreadsRule, `
      fragment unionWithinUnion on DogOrHuman { ...catOrDogFragment }
      fragment catOrDogFragment on CatOrDog { __typename }
    `)
}
func TestValidate_PossibleFragmentSpreads_InterfaceIntoImplementedObject(t *testing.T) {
	expectPassesRule(t, graphql.PossibleFragmentSpreadsRule, `
      fragment interfaceWithinObject on Dog { ...petFragment }
      fragment petFragment on Pet { name }
    `)
}
func TestValidate_PossibleFragmentSpreads_InterfaceIntoOverlappingInterface(t *testing.T) {
	expectPassesRule(t, graphql.PossibleFragmentSpreadsRule, `
      fragment interfaceWithinInterface on Pet { ...beingFragment }
      fragment beingFragment on Being { name }
    `)
}
func TestValidate_PossibleFragmentSpreads_InterfaceIntoOverlappingInterfaceInInlineFragment(t *testing.T) {
	expectPassesRule(t, graphql.PossibleFragmentSpreadsRule, `
      fragment interfaceWithinInterface on Pet { ... on Being { name } }
    `)
}
func TestValidate_PossibleFragmentSpreads_InterfaceIntoOverlappingUnion(t *testing.T) {
	expectPassesRule(t, graphql.PossibleFragmentSpreadsRule, `
      fragment interfaceWithinUnion on CatOrDog { ...petFragment }
      fragment petFragment on Pet { name }
    `)
}
func TestValidate_PossibleFragmentSpreads_DifferentObjectIntoObject(t *testing.T) {
	expectFailsRule(t, graphql.PossibleFragmentSpreadsRule, `
      fragment invalidObjectWithinObject on Cat { ...dogFragment }
      fragment dogFragment on Dog { barkVolume }
    `, []gqlerrors.FormattedError{
		ruleError(`Fragment "dogFragment" cannot be spread here as objects of `+
			`type "Cat" can never be of type "Dog".`, 2, 51),
	})
}

func TestValidate_PossibleFragmentSpreads_DifferentObjectIntoObjectInInlineFragment(t *testing.T) {
	expectFailsRule(t, graphql.PossibleFragmentSpreadsRule, `
      fragment invalidObjectWithinObjectAnon on Cat {
        ... on Dog { barkVolume }
      }
    `, []gqlerrors.FormattedError{
		ruleError(`Fragment cannot be spread here as objects of `+
			`type "Cat" can never be of type "Dog".`, 3, 9),
	})
}
func TestValidate_PossibleFragmentSpreads_ObjectIntoNotImplementingInterface(t *testing.T) {
	expectFailsRule(t, graphql.PossibleFragmentSpreadsRule, `
      fragment invalidObjectWithinInterface on Pet { ...humanFragment }
      fragment humanFragment on Human { pets { name } }
    `, []gqlerrors.FormattedError{
		ruleError(`Fragment "humanFragment" cannot be spread here as objects of `+
			`type "Pet" can never be of type "Human".`, 2, 54),
	})
}
func TestValidate_PossibleFragmentSpreads_ObjectIntoNotContainingUnion(t *testing.T) {
	expectFailsRule(t, graphql.PossibleFragmentSpreadsRule, `
      fragment invalidObjectWithinUnion on CatOrDog { ...humanFragment }
      fragment humanFragment on Human { pets { name } }
    `, []gqlerrors.FormattedError{
		ruleError(`Fragment "humanFragment" cannot be spread here as objects of `+
			`type "CatOrDog" can never be of type "Human".`, 2, 55),
	})
}
func TestValidate_PossibleFragmentSpreads_UnionIntoNotContainedObject(t *testing.T) {
	expectFailsRule(t, graphql.PossibleFragmentSpreadsRule, `
      fragment invalidUnionWithinObject on Human { ...catOrDogFragment }
      fragment catOrDogFragment on CatOrDog { __typename }
    `, []gqlerrors.FormattedError{
		ruleError(`Fragment "catOrDogFragment" cannot be spread here as objects of `+
			`type "Human" can never be of type "CatOrDog".`, 2, 52),
	})
}
func TestValidate_PossibleFragmentSpreads_UnionIntoNonOverlappingInterface(t *testing.T) {
	expectFailsRule(t, graphql.PossibleFragmentSpreadsRule, `
      fragment invalidUnionWithinInterface on Pet { ...humanOrAlienFragment }
      fragment humanOrAlienFragment on HumanOrAlien { __typename }
    `, []gqlerrors.FormattedError{
		ruleError(`Fragment "humanOrAlienFragment" cannot be spread here as objects of `+
			`type "Pet" can never be of type "HumanOrAlien".`, 2, 53),
	})
}

func TestValidate_PossibleFragmentSpreads_UnionIntoNonOverlappingUnion(t *testing.T) {
	expectFailsRule(t, graphql.PossibleFragmentSpreadsRule, `
      fragment invalidUnionWithinUnion on CatOrDog { ...humanOrAlienFragment }
      fragment humanOrAlienFragment on HumanOrAlien { __typename }
    `, []gqlerrors.FormattedError{
		ruleError(`Fragment "humanOrAlienFragment" cannot be spread here as objects of `+
			`type "CatOrDog" can never be of type "HumanOrAlien".`, 2, 54),
	})
}

func TestValidate_PossibleFragmentSpreads_InterfaceIntoNonImplementingObject(t *testing.T) {
	expectFailsRule(t, graphql.PossibleFragmentSpreadsRule, `
      fragment invalidInterfaceWithinObject on Cat { ...intelligentFragment }
      fragment intelligentFragment on Intelligent { iq }
    `, []gqlerrors.FormattedError{
		ruleError(`Fragment "intelligentFragment" cannot be spread here as objects of `+
			`type "Cat" can never be of type "Intelligent".`, 2, 54),
	})
}
func TestValidate_PossibleFragmentSpreads_InterfaceIntoNonOverlappingInterface(t *testing.T) {
	expectFailsRule(t, graphql.PossibleFragmentSpreadsRule, `
      fragment invalidInterfaceWithinInterface on Pet {
        ...intelligentFragment
      }
      fragment intelligentFragment on Intelligent { iq }
    `, []gqlerrors.FormattedError{
		ruleError(`Fragment "intelligentFragment" cannot be spread here as objects of `+
			`type "Pet" can never be of type "Intelligent".`, 3, 9),
	})
}
func TestValidate_PossibleFragmentSpreads_InterfaceIntoNonOverlappingInterfaceInInlineFragment(t *testing.T) {
	expectFailsRule(t, graphql.PossibleFragmentSpreadsRule, `
      fragment invalidInterfaceWithinInterfaceAnon on Pet {
        ...on Intelligent { iq }
      }
    `, []gqlerrors.FormattedError{
		ruleError(`Fragment cannot be spread here as objects of `+
			`type "Pet" can never be of type "Intelligent".`, 3, 9),
	})
}
func TestValidate_PossibleFragmentSpreads_InterfaceIntoNonOverlappingUnion(t *testing.T) {
	expectFailsRule(t, graphql.PossibleFragmentSpreadsRule, `
      fragment invalidInterfaceWithinUnion on HumanOrAlien { ...petFragment }
      fragment petFragment on Pet { name }
    `, []gqlerrors.FormattedError{
		ruleError(`Fragment "petFragment" cannot be spread here as objects of `+
			`type "HumanOrAlien" can never be of type "Pet".`, 2, 62),
	})
}
