package graphql

import (
	"fmt"
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
	"github.com/graphql-go/graphql/language/printer"
	"github.com/graphql-go/graphql/language/visitor"
)

/**
 * SpecifiedRules set includes all validation rules defined by the GraphQL spec.
 */
var SpecifiedRules = []ValidationRuleFn{
	ArgumentsOfCorrectTypeRule,
	KnownTypeNamesRule,
	DefaultValuesOfCorrectTypeRule,
}

type ValidationRuleInstance struct {
	VisitorOpts          *visitor.VisitorOptions
	VisitSpreadFragments bool
}
type ValidationRuleFn func(context *ValidationContext) *ValidationRuleInstance

func ArgumentsOfCorrectTypeRule(context *ValidationContext) *ValidationRuleInstance {
	visitorOpts := &visitor.VisitorOptions{
		KindFuncMap: map[string]visitor.NamedVisitFuncs{
			kinds.Argument: visitor.NamedVisitFuncs{
				Kind: func(p visitor.VisitFuncParams) (string, interface{}) {
					var action = visitor.ActionNoChange
					var result interface{}
					if argAST, ok := p.Node.(*ast.Argument); ok {
						value := argAST.Value
						argDef := context.GetArgument()
						if argDef != nil && !isValidLiteralValue(argDef.Type, value) {
							argNameValue := ""
							if argAST.Name != nil {
								argNameValue = argAST.Name.Value
							}
							// TODO: helper to construct gqlerror with message + []ast.Node
							return visitor.ActionNoChange, gqlerrors.NewError(
								fmt.Sprintf(`Argument "%v" expected type "%v" but got: %v.`,
									argNameValue, argDef.Type, printer.Print(value)),
								[]ast.Node{value},
								"",
								nil,
								[]int{},
							)
						}
					}
					return action, result
				},
			},
		},
	}
	return &ValidationRuleInstance{
		VisitorOpts: visitorOpts,
	}
}
func DefaultValuesOfCorrectTypeRule(context *ValidationContext) *ValidationRuleInstance {
	visitorOpts := &visitor.VisitorOptions{
		KindFuncMap: map[string]visitor.NamedVisitFuncs{
			kinds.VariableDefinition: visitor.NamedVisitFuncs{
				Kind: func(p visitor.VisitFuncParams) (string, interface{}) {
					var action = visitor.ActionNoChange
					var result interface{}
					if varDefAST, ok := p.Node.(*ast.VariableDefinition); ok {
						name := ""
						if varDefAST.Variable != nil && varDefAST.Variable.Name != nil {
							name = varDefAST.Variable.Name.Value
						}
						defaultValue := varDefAST.DefaultValue
						ttype := context.GetInputType()

						if ttype, ok := ttype.(*NonNull); ok && defaultValue != nil {
							return visitor.ActionNoChange, gqlerrors.NewError(
								fmt.Sprintf(`Variable "$%v" of type "%v" is required and will not use the default value. Perhaps you meant to use type "%v".`,
									name, ttype, ttype.OfType),
								[]ast.Node{defaultValue},
								"",
								nil,
								[]int{},
							)
						}
						if ttype != nil && defaultValue != nil && !isValidLiteralValue(ttype, defaultValue) {
							return visitor.ActionNoChange, gqlerrors.NewError(
								fmt.Sprintf(`Variable "$%v" of type "%v" has invalid default value: %v.`,
									name, ttype, printer.Print(defaultValue)),
								[]ast.Node{defaultValue},
								"",
								nil,
								[]int{},
							)
						}
					}
					return action, result
				},
			},
		},
	}
	return &ValidationRuleInstance{
		VisitorOpts: visitorOpts,
	}
}
func KnownTypeNamesRule(context *ValidationContext) *ValidationRuleInstance {
	visitorOpts := &visitor.VisitorOptions{
		KindFuncMap: map[string]visitor.NamedVisitFuncs{
			kinds.Named: visitor.NamedVisitFuncs{
				Kind: func(p visitor.VisitFuncParams) (string, interface{}) {
					if node, ok := p.Node.(*ast.Named); ok {
						typeNameValue := ""
						typeName := node.Name
						if typeName != nil {
							typeNameValue = typeName.Value
						}
						ttype := context.GetSchema().GetType(typeNameValue)
						if ttype == nil {
							return visitor.ActionNoChange, gqlerrors.NewError(
								fmt.Sprintf(`Unknown type "%v".`, typeNameValue),
								[]ast.Node{node},
								"",
								nil,
								[]int{},
							)
						}
					}
					return visitor.ActionNoChange, nil
				},
			},
		},
	}
	return &ValidationRuleInstance{
		VisitorOpts: visitorOpts,
	}
}

/**
 * Utility for validators which determines if a value literal AST is valid given
 * an input type.
 *
 * Note that this only validates literal values, variables are assumed to
 * provide values of the correct type.
 */
func isValidLiteralValue(ttype Input, valueAST ast.Value) bool {
	// A value must be provided if the type is non-null.
	if ttype, ok := ttype.(*NonNull); ok {
		if valueAST == nil {
			return false
		}
		ofType, _ := ttype.OfType.(Input)
		return isValidLiteralValue(ofType, valueAST)
	}

	if valueAST == nil {
		return true
	}

	// This function only tests literals, and assumes variables will provide
	// values of the correct type.
	if valueAST.GetKind() == kinds.Variable {
		return true
	}

	// Lists accept a non-list value as a list of one.
	if ttype, ok := ttype.(*List); ok {
		itemType, _ := ttype.OfType.(Input)
		if valueAST, ok := valueAST.(*ast.ListValue); ok {
			for _, value := range valueAST.Values {
				if isValidLiteralValue(itemType, value) == false {
					return false
				}
			}
			return true
		}
		return isValidLiteralValue(itemType, valueAST)

	}

	// Input objects check each defined field and look for undefined fields.
	if ttype, ok := ttype.(*InputObject); ok {
		valueAST, ok := valueAST.(*ast.ObjectValue)
		if !ok {
			return false
		}
		fields := ttype.GetFields()

		// Ensure every provided field is defined.
		// Ensure every defined field is valid.
		fieldASTs := valueAST.Fields
		fieldASTMap := map[string]*ast.ObjectField{}
		for _, fieldAST := range fieldASTs {
			fieldASTName := ""
			if fieldAST.Name != nil {
				fieldASTName = fieldAST.Name.Value
			}

			fieldASTMap[fieldASTName] = fieldAST

			// check if field is defined
			field, ok := fields[fieldASTName]
			if !ok || field == nil {
				return false
			}
		}
		for fieldName, field := range fields {
			fieldAST, _ := fieldASTMap[fieldName]
			var fieldASTValue ast.Value
			if fieldAST != nil {
				fieldASTValue = fieldAST.Value
			}
			if !isValidLiteralValue(field.Type, fieldASTValue) {
				return false
			}
		}
		return true
	}

	if ttype, ok := ttype.(*Scalar); ok {
		return !isNullish(ttype.ParseLiteral(valueAST))
	}
	if ttype, ok := ttype.(*Enum); ok {
		return !isNullish(ttype.ParseLiteral(valueAST))
	}

	// Must be input type (not scalar or enum)
	// Silently fail, instead of panic()
	return false
}
