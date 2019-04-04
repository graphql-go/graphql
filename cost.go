package graphql

import "github.com/GannettDigital/graphql/language/ast"

type fieldDefiner interface {
	Fields() FieldDefinitionMap
}

// QueryComplexity returns the complexity cost of the given query.
//
// The cost is calculated by adding up the costs of the various fields
func QueryComplexity(p ExecuteParams) (int, error) {
	var cost int

	exeContext, err := buildExecutionContext(buildExecutionCtxParams{
		Schema:        p.Schema,
		Root:          p.Root,
		AST:           p.AST,
		OperationName: p.OperationName,
		Args:          p.Args,
		Errors:        nil,
		Result:        &Result{},
		Context:       p.Context,
	})
	if err != nil {
		return 0, err
	}

	operationType, err := getOperationRootType(p.Schema, exeContext.Operation)
	if err != nil {
		return 0, err
	}

	cost += selectionSetCost(exeContext.Operation.GetSelectionSet(), operationType, exeContext)

	return cost, nil
}

// astFieldCost will recursively determine the cost of a field including its children.
func astFieldCost(field *ast.Field, fieldDef *FieldDefinition, exeContext *executionContext) int {
	cost := fieldDef.Cost

	set := field.GetSelectionSet()
	if set == nil {
		return cost
	}
	fType := fieldDef.Type
	if nonNullType, ok := fieldDef.Type.(*NonNull); ok {
		fType = nonNullType.OfType
	}
	if listType, ok := fType.(*List); ok {
		fType = listType.OfType
	}
	if nonNullType, ok := fType.(*NonNull); ok {
		fType = nonNullType.OfType
	}
	parent, ok := fType.(fieldDefiner)
	if !ok {
		return cost
	}
	cost += selectionSetCost(set, parent, exeContext)

	return cost
}

// selectionSetCost will return the cost for a given selection set.
func selectionSetCost(set *ast.SelectionSet, parent fieldDefiner, exeContext *executionContext) int {
	if set == nil {
		return 0
	}
	var cost int

	for _, iSelection := range set.Selections {
		switch selection := iSelection.(type) {
		case *ast.Field:
			fieldDef, ok := parent.Fields()[selection.Name.Value]
			if !ok {
				continue
			}
			cost += astFieldCost(selection, fieldDef, exeContext)
		case *ast.InlineFragment:
			selectionType := selection.TypeCondition
			parentInterface, ok := parent.(*Interface)
			if !ok || selectionType == nil || parentInterface == nil {
				cost += selectionSetCost(selection.SelectionSet, parent, exeContext)
				continue
			}
			for _, object := range exeContext.Schema.implementations[parentInterface.Name()] {
				if object.Name() == selectionType.Name.Value {
					cost += selectionSetCost(selection.SelectionSet, object, exeContext)
				}
			}
		case *ast.FragmentSpread:
			fragment, ok := exeContext.Fragments[selection.Name.Value]
			if !ok {
				continue
			}
			fragmentDef, ok := fragment.(*ast.FragmentDefinition)
			if !ok {
				continue
			}
			fragmentType, err := typeFromAST(exeContext.Schema, fragmentDef.TypeCondition)
			if err != nil {
				continue
			}
			fragmentObject, ok := fragmentType.(fieldDefiner)
			if !ok {
				continue
			}
			cost += selectionSetCost(fragment.GetSelectionSet(), fragmentObject, exeContext)
		}
	}

	return cost
}
