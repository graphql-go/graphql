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

	typeInfo := NewTypeInfo(&TypeInfoConfig{
		Schema: schema,
	})
	vr.Errors = VisitUsingRules(schema, typeInfo, astDoc, rules)
	if len(vr.Errors) == 0 {
		vr.IsValid = true
	}
	return vr
}

/**
 * VisitUsingRules This uses a specialized visitor which runs multiple visitors in parallel,
 * while maintaining the visitor skip and break API.
 * @internal
 */
func VisitUsingRules(schema *Schema, typeInfo *TypeInfo, astDoc *ast.Document, rules []ValidationRuleFn) []gqlerrors.FormattedError {

	context := NewValidationContext(schema, astDoc, typeInfo)
	visitors := []*visitor.VisitorOptions{}

	for _, rule := range rules {
		instance := rule(context)
		visitors = append(visitors, instance.VisitorOpts)
	}

	// Visit the whole document with each instance of all provided rules.
	visitor.Visit(astDoc, visitor.VisitWithTypeInfo(typeInfo, visitor.VisitInParallel(visitors...)), nil)
	return context.Errors()
}

func visitUsingRulesOld(schema *Schema, astDoc *ast.Document, rules []ValidationRuleFn) []gqlerrors.FormattedError {
	typeInfo := NewTypeInfo(&TypeInfoConfig{
		Schema: schema,
	})
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
					enterFn := visitor.GetVisitFn(instance.VisitorOpts, kind, false)
					if enterFn != nil {
						action, result = enterFn(p)
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
					leaveFn := visitor.GetVisitFn(instance.VisitorOpts, kind, true)
					if leaveFn != nil {
						action, result = leaveFn(p)
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
	return context.Errors()
}

type HasSelectionSet interface {
	GetKind() string
	GetLoc() *ast.Location
	GetSelectionSet() *ast.SelectionSet
}

var _ HasSelectionSet = (*ast.OperationDefinition)(nil)
var _ HasSelectionSet = (*ast.FragmentDefinition)(nil)

type VariableUsage struct {
	Node *ast.Variable
	Type Input
}

type ValidationContext struct {
	schema                         *Schema
	astDoc                         *ast.Document
	typeInfo                       *TypeInfo
	errors                         []gqlerrors.FormattedError
	fragments                      map[string]*ast.FragmentDefinition
	variableUsages                 map[HasSelectionSet][]*VariableUsage
	recursiveVariableUsages        map[*ast.OperationDefinition][]*VariableUsage
	recursivelyReferencedFragments map[*ast.OperationDefinition][]*ast.FragmentDefinition
	fragmentSpreads                map[HasSelectionSet][]*ast.FragmentSpread
}

func NewValidationContext(schema *Schema, astDoc *ast.Document, typeInfo *TypeInfo) *ValidationContext {
	return &ValidationContext{
		schema:                         schema,
		astDoc:                         astDoc,
		typeInfo:                       typeInfo,
		fragments:                      map[string]*ast.FragmentDefinition{},
		variableUsages:                 map[HasSelectionSet][]*VariableUsage{},
		recursiveVariableUsages:        map[*ast.OperationDefinition][]*VariableUsage{},
		recursivelyReferencedFragments: map[*ast.OperationDefinition][]*ast.FragmentDefinition{},
		fragmentSpreads:                map[HasSelectionSet][]*ast.FragmentSpread{},
	}
}

func (ctx *ValidationContext) ReportError(err error) {
	formattedErr := gqlerrors.FormatError(err)
	ctx.errors = append(ctx.errors, formattedErr)
}
func (ctx *ValidationContext) Errors() []gqlerrors.FormattedError {
	return ctx.errors
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
func (ctx *ValidationContext) FragmentSpreads(node HasSelectionSet) []*ast.FragmentSpread {
	if spreads, ok := ctx.fragmentSpreads[node]; ok && spreads != nil {
		return spreads
	}

	spreads := []*ast.FragmentSpread{}
	setsToVisit := []*ast.SelectionSet{node.GetSelectionSet()}

	for {
		if len(setsToVisit) == 0 {
			break
		}
		var set *ast.SelectionSet
		// pop
		set, setsToVisit = setsToVisit[len(setsToVisit)-1], setsToVisit[:len(setsToVisit)-1]
		if set.Selections != nil {
			for _, selection := range set.Selections {
				switch selection := selection.(type) {
				case *ast.FragmentSpread:
					spreads = append(spreads, selection)
				case *ast.Field:
					if selection.SelectionSet != nil {
						setsToVisit = append(setsToVisit, selection.SelectionSet)
					}
				case *ast.InlineFragment:
					if selection.SelectionSet != nil {
						setsToVisit = append(setsToVisit, selection.SelectionSet)
					}
				}
			}
		}
		ctx.fragmentSpreads[node] = spreads
	}
	return spreads
}

func (ctx *ValidationContext) RecursivelyReferencedFragments(operation *ast.OperationDefinition) []*ast.FragmentDefinition {
	if fragments, ok := ctx.recursivelyReferencedFragments[operation]; ok && fragments != nil {
		return fragments
	}

	fragments := []*ast.FragmentDefinition{}
	collectedNames := map[string]bool{}
	nodesToVisit := []HasSelectionSet{operation}

	for {
		if len(nodesToVisit) == 0 {
			break
		}

		var node HasSelectionSet

		node, nodesToVisit = nodesToVisit[len(nodesToVisit)-1], nodesToVisit[:len(nodesToVisit)-1]
		spreads := ctx.FragmentSpreads(node)
		for _, spread := range spreads {
			fragName := ""
			if spread.Name != nil {
				fragName = spread.Name.Value
			}
			if res, ok := collectedNames[fragName]; !ok || !res {
				collectedNames[fragName] = true
				fragment := ctx.Fragment(fragName)
				if fragment != nil {
					fragments = append(fragments, fragment)
					nodesToVisit = append(nodesToVisit, fragment)
				}
			}

		}
	}

	ctx.recursivelyReferencedFragments[operation] = fragments
	return fragments
}
func (ctx *ValidationContext) VariableUsages(node HasSelectionSet) []*VariableUsage {
	if usages, ok := ctx.variableUsages[node]; ok && usages != nil {
		return usages
	}
	usages := []*VariableUsage{}
	typeInfo := NewTypeInfo(&TypeInfoConfig{
		Schema: ctx.schema,
	})

	visitor.Visit(node, visitor.VisitWithTypeInfo(typeInfo, &visitor.VisitorOptions{
		KindFuncMap: map[string]visitor.NamedVisitFuncs{
			kinds.VariableDefinition: visitor.NamedVisitFuncs{
				Kind: func(p visitor.VisitFuncParams) (string, interface{}) {
					return visitor.ActionSkip, nil
				},
			},
			kinds.Variable: visitor.NamedVisitFuncs{
				Kind: func(p visitor.VisitFuncParams) (string, interface{}) {
					if node, ok := p.Node.(*ast.Variable); ok && node != nil {
						usages = append(usages, &VariableUsage{
							Node: node,
							Type: typeInfo.InputType(),
						})
					}
					return visitor.ActionNoChange, nil
				},
			},
		},
	}), nil)

	ctx.variableUsages[node] = usages
	return usages
}
func (ctx *ValidationContext) RecursiveVariableUsages(operation *ast.OperationDefinition) []*VariableUsage {
	if usages, ok := ctx.recursiveVariableUsages[operation]; ok && usages != nil {
		return usages
	}
	usages := ctx.VariableUsages(operation)

	fragments := ctx.RecursivelyReferencedFragments(operation)
	for _, fragment := range fragments {
		fragmentUsages := ctx.VariableUsages(fragment)
		usages = append(usages, fragmentUsages...)
	}

	ctx.recursiveVariableUsages[operation] = usages
	return usages
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
