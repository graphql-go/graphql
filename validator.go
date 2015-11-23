package graphql

import (
	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
	"github.com/graphql-go/graphql/language/visitor"
)

type ValidationResult struct {
	IsValid bool
	Errors  []gqlerrors.FormattedError
}

func ValidateDocument(schema *Schema, astDoc *ast.Document, rules []ValidationRuleFn) (vr ValidationResult) {
	if len(rules) == 0 {
		rules = SpecifiedRules
	}

	vr.IsValid = false
	if schema == nil {
		vr.Errors = append(vr.Errors, gqlerrors.NewFormattedError("Must provide schema"))
		return vr
	}
	if astDoc == nil {
		vr.Errors = append(vr.Errors, gqlerrors.NewFormattedError("Must provide document"))
		return vr
	}
	vr.Errors = visitUsingRules(schema, astDoc, rules)
	if len(vr.Errors) == 0 {
		vr.IsValid = true
	}
	return vr
}

func visitUsingRules(schema *Schema, astDoc *ast.Document, rules []ValidationRuleFn) (errors []gqlerrors.FormattedError) {
	typeInfo := NewTypeInfo(schema)
	context := NewValidationContext(schema, astDoc, typeInfo)

	var visitInstance func(astNode ast.Node, instance *ValidationRuleInstance)

	visitInstance = func(astNode ast.Node, instance *ValidationRuleInstance) {
		visitor.Visit(astNode, &visitor.VisitorOptions{
			Enter: func(p visitor.VisitFuncParams) (string, interface{}) {
				var action = visitor.ActionNoChange
				var result interface{}
				switch node := p.Node.(type) {
				case ast.Node:
					// Collect type information about the current position in the AST.
					typeInfo.Enter(node)

					// Do not visit top level fragment definitions if this instance will
					// visit those fragments inline because it
					// provided `visitSpreadFragments`.
					kind := node.GetKind()

					if kind == kinds.FragmentDefinition &&
						p.Key != nil && instance.VisitSpreadFragments == true {
						return visitor.ActionSkip, nil
					}

					// Get the visitor function from the validation instance, and if it
					// exists, call it with the visitor arguments.
					enterFn := visitor.GetVisitFn(instance.VisitorOpts, false, kind)
					if enterFn != nil {
						action, result = enterFn(p)
					}

					// If the visitor returned an error, log it and do not visit any
					// deeper nodes.
					if err, ok := result.(error); ok && err != nil {
						errors = append(errors, gqlerrors.FormatError(err))
						action = visitor.ActionSkip
					}
					if err, ok := result.([]error); ok && err != nil {
						errors = append(errors, gqlerrors.FormatErrors(err...)...)
						action = visitor.ActionSkip
					}

					// If any validation instances provide the flag `visitSpreadFragments`
					// and this node is a fragment spread, visit the fragment definition
					// from this point.
					if action == visitor.ActionNoChange && result == nil &&
						instance.VisitSpreadFragments == true && kind == kinds.FragmentSpread {
						node, _ := node.(*ast.FragmentSpread)
						name := node.Name
						nameVal := ""
						if name != nil {
							nameVal = name.Value
						}
						fragment := context.Fragment(nameVal)
						if fragment != nil {
							visitInstance(fragment, instance)
						}
					}

					// If the result is "false" (ie action === Action.Skip), we're not visiting any descendent nodes,
					// but need to update typeInfo.
					if action == visitor.ActionSkip {
						typeInfo.Leave(node)
					}

				}

				return action, result
			},
			Leave: func(p visitor.VisitFuncParams) (string, interface{}) {
				var action = visitor.ActionNoChange
				var result interface{}
				switch node := p.Node.(type) {
				case ast.Node:
					kind := node.GetKind()

					// Get the visitor function from the validation instance, and if it
					// exists, call it with the visitor arguments.
					leaveFn := visitor.GetVisitFn(instance.VisitorOpts, true, kind)
					if leaveFn != nil {
						action, result = leaveFn(p)
					}

					// If the visitor returned an error, log it and do not visit any
					// deeper nodes.
					if err, ok := result.(error); ok && err != nil {
						errors = append(errors, gqlerrors.FormatError(err))
						action = visitor.ActionSkip
					}
					if err, ok := result.([]error); ok && err != nil {
						errors = append(errors, gqlerrors.FormatErrors(err...)...)
						action = visitor.ActionSkip
					}

					// Update typeInfo.
					typeInfo.Leave(node)
				}
				return action, result
			},
		}, nil)
	}

	instances := []*ValidationRuleInstance{}
	for _, rule := range rules {
		instance := rule(context)
		instances = append(instances, instance)
	}
	for _, instance := range instances {
		visitInstance(astDoc, instance)
	}
	return errors
}

type ValidationContext struct {
	schema    *Schema
	astDoc    *ast.Document
	typeInfo  *TypeInfo
	fragments map[string]*ast.FragmentDefinition
}

func NewValidationContext(schema *Schema, astDoc *ast.Document, typeInfo *TypeInfo) *ValidationContext {
	return &ValidationContext{
		schema:   schema,
		astDoc:   astDoc,
		typeInfo: typeInfo,
	}
}

func (ctx *ValidationContext) Schema() *Schema {
	return ctx.schema
}
func (ctx *ValidationContext) Document() *ast.Document {
	return ctx.astDoc
}

func (ctx *ValidationContext) Fragment(name string) *ast.FragmentDefinition {
	if len(ctx.fragments) == 0 {
		if ctx.Document() == nil {
			return nil
		}
		defs := ctx.Document().Definitions
		fragments := map[string]*ast.FragmentDefinition{}
		for _, def := range defs {
			if def, ok := def.(*ast.FragmentDefinition); ok {
				defName := ""
				if def.Name != nil {
					defName = def.Name.Value
				}
				fragments[defName] = def
			}
		}
		ctx.fragments = fragments
	}
	f, _ := ctx.fragments[name]
	return f
}

func (ctx *ValidationContext) Type() Output {
	return ctx.typeInfo.Type()
}
func (ctx *ValidationContext) ParentType() Composite {
	return ctx.typeInfo.ParentType()
}
func (ctx *ValidationContext) InputType() Input {
	return ctx.typeInfo.InputType()
}
func (ctx *ValidationContext) FieldDef() *FieldDefinition {
	return ctx.typeInfo.FieldDef()
}
func (ctx *ValidationContext) Directive() *Directive {
	return ctx.typeInfo.Directive()
}
func (ctx *ValidationContext) Argument() *Argument {
	return ctx.typeInfo.Argument()
}
