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

	fields := collectFields(collectFieldsParams{
		ExeContext:   exeContext,
		RuntimeType:  operationType,
		SelectionSet: exeContext.Operation.GetSelectionSet(),
	})

	for _, fieldASTs := range fields {
		for _, field := range fieldASTs {
			fieldDef, ok := operationType.Fields()[field.Name.Value]
			if !ok {
				continue
			}
			cost += astFieldCost(field, fieldDef)
		}
	}

	return cost, nil
}

// astFieldCost will recursively determine the cost of a field including its children.
func astFieldCost(field *ast.Field, fieldDef *FieldDefinition) int {
	cost := fieldDef.Cost

	if field.SelectionSet == nil {
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

	for _, s := range field.SelectionSet.Selections {
		if f, ok := s.(*ast.Field); ok {
			fieldDef, ok := parent.Fields()[f.Name.Value]
			if !ok {
				continue
			}
			cost += astFieldCost(f, fieldDef)
		}
	}

	return cost
}
