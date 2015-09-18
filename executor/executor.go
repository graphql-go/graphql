package executor

import (
	"fmt"

	"github.com/chris-ramon/graphql-go/errors"
	"github.com/chris-ramon/graphql-go/language/ast"
	"github.com/chris-ramon/graphql-go/types"
	"reflect"
)

type ExecuteParams struct {
	Schema        types.GraphQLSchema
	Root          map[string]interface{}
	AST           *ast.Document
	OperationName string
	Args          map[string]interface{}
}

func Execute(p ExecuteParams, resultChan chan *types.GraphQLResult) {
	var errors []graphqlerrors.GraphQLFormattedError
	var result types.GraphQLResult
	params := BuildExecutionCtxParams{
		Schema:        p.Schema,
		Root:          p.Root,
		AST:           p.AST,
		OperationName: p.OperationName,
		Args:          p.Args,
		Errors:        errors,
		Result:        &result,
		ResultChan:    resultChan,
	}
	exeContext := buildExecutionContext(params)
	if result.HasErrors() {
		return
	}
	eOperationParams := ExecuteOperationParams{
		ExecutionContext: exeContext,
		Root:             p.Root,
		Operation:        exeContext.Operation,
	}
	executeOperation(eOperationParams, resultChan)
}

type BuildExecutionCtxParams struct {
	Schema        types.GraphQLSchema
	Root          map[string]interface{}
	AST           *ast.Document
	OperationName string
	Args          map[string]interface{}
	Errors        []graphqlerrors.GraphQLFormattedError
	Result        *types.GraphQLResult
	ResultChan    chan *types.GraphQLResult
}
type ExecutionContext struct {
	Schema         types.GraphQLSchema
	Fragments      map[string]ast.Definition
	Root           map[string]interface{}
	Operation      ast.Definition
	VariableValues map[string]interface{}
	Errors         []graphqlerrors.GraphQLFormattedError
}

func buildExecutionContext(p BuildExecutionCtxParams) (eCtx ExecutionContext) {
	operations := map[string]ast.Definition{}
	fragments := map[string]ast.Definition{}
	for _, statement := range p.AST.Definitions {
		switch stm := statement.(type) {
		case *ast.OperationDefinition:
			key := ""
			if stm.GetName() != nil && stm.GetName().Value != "" {
				key = stm.GetName().Value
			}
			operations[key] = stm
		case *ast.FragmentDefinition:
			key := ""
			if stm.GetName() != nil && stm.GetName().Value != "" {
				key = stm.GetName().Value
			}
			fragments[key] = stm
		default:
			err := graphqlerrors.NewGraphQLFormattedError(
				fmt.Sprintf("GraphQL cannot execute a request containing a %v", statement.GetKind()),
			)
			p.Result.Errors = append(p.Result.Errors, err)
			p.ResultChan <- p.Result
			return eCtx
		}
	}
	if (p.OperationName == "") && (len(operations) != 1) {
		err := graphqlerrors.NewGraphQLFormattedError("Must provide operation name if query contains multiple operations.")
		p.Result.Errors = append(p.Result.Errors, err)
		p.ResultChan <- p.Result
		return eCtx
	}
	opName := p.OperationName
	if opName == "" {
		// get first opName
		for k, _ := range operations {
			opName = k
			break
		}
	}
	operation, found := operations[opName]
	if !found {
		err := graphqlerrors.NewGraphQLFormattedError(fmt.Sprintf(`Unknown operation named "%v".`, opName))
		p.Result.Errors = append(p.Result.Errors, err)
		p.ResultChan <- p.Result
		return eCtx
	}
	variableValues, err := getVariableValues(p.Schema, operation.GetVariableDefinitions(), p.Args)
	if err != nil {
		p.Result.Errors = append(p.Result.Errors, graphqlerrors.FormatError(err))
		p.ResultChan <- p.Result
		return eCtx
	}

	eCtx.Schema = p.Schema
	eCtx.Fragments = fragments
	eCtx.Root = p.Root
	eCtx.Operation = operation
	eCtx.VariableValues = variableValues
	eCtx.Errors = p.Errors
	return eCtx
}

type ExecuteOperationParams struct {
	ExecutionContext ExecutionContext
	Root             map[string]interface{}
	Operation        ast.Definition
}

func executeOperation(p ExecuteOperationParams, r chan *types.GraphQLResult) {
	operationType := getOperationRootType(p.ExecutionContext.Schema, p.Operation, r)

	collectFieldsParams := CollectFieldsParams{
		ExeContext:    p.ExecutionContext,
		OperationType: operationType,
		SelectionSet:  p.Operation.GetSelectionSet(),
	}
	fields := collectFields(collectFieldsParams)
	executeFieldsParams := ExecuteFieldsParams{
		ExecutionContext: p.ExecutionContext,
		ParentType:       operationType,
		Source:           p.Root,
		Fields:           fields,
	}
	if p.Operation.GetOperation() == "mutation" {
		executeFieldsSerially(executeFieldsParams, r)
		return
	}
	executeFields(executeFieldsParams, r)
}

// Extracts the root type of the operation from the schema.
func getOperationRootType(schema types.GraphQLSchema, operation ast.Definition, r chan *types.GraphQLResult) (objType *types.GraphQLObjectType) {
	if operation == nil {
		var result types.GraphQLResult
		err := graphqlerrors.NewGraphQLFormattedError("Can only execute queries and mutations")
		result.Errors = append(result.Errors, err)
		r <- &result
		return objType
	}
	switch operation.GetOperation() {
	case "query":
		return schema.GetQueryType()
	case "mutation":
		mutationType := schema.GetMutationType()
		if mutationType.Name == "" {
			var result types.GraphQLResult
			err := graphqlerrors.NewGraphQLFormattedError("Schema is not configured for mutations")
			result.Errors = append(result.Errors, err)
			r <- &result
			return objType
		}
		return mutationType
	default:
		var result types.GraphQLResult
		err := graphqlerrors.NewGraphQLFormattedError("Can only execute queries and mutations")
		result.Errors = append(result.Errors, err)
		r <- &result
		return objType
	}
}

type ExecuteFieldsParams struct {
	ExecutionContext ExecutionContext
	ParentType       *types.GraphQLObjectType
	Source           interface{}
	Fields           map[string][]*ast.Field
}

// Implements the "Evaluating selection sets" section of the spec for "write" mode.
func executeFieldsSerially(p ExecuteFieldsParams, resultChan chan *types.GraphQLResult) {
	if p.Source == nil {
		p.Source = map[string]interface{}{}
	}
	if p.Fields == nil {
		p.Fields = map[string][]*ast.Field{}
	}
	var result types.GraphQLResult

	finalResults := map[string]interface{}{}
	for responseName, fieldASTs := range p.Fields {
		resolved, hasNoFieldDefs, errs := resolveField(p.ExecutionContext, p.ParentType, p.Source, fieldASTs)
		if hasNoFieldDefs {
			continue
		}
		for _, err := range errs {
			result.Errors = append(result.Errors, err)
		}
		finalResults[responseName] = resolved
	}
	result.Data = finalResults
	resultChan <- &result
}

// Implements the "Evaluating selection sets" section of the spec for "read" mode.
func executeFields(p ExecuteFieldsParams, resultChan chan *types.GraphQLResult) {
	if p.Source == nil {
		p.Source = map[string]interface{}{}
	}
	if p.Fields == nil {
		p.Fields = map[string][]*ast.Field{}
	}
	var result types.GraphQLResult

	finalResults := map[string]interface{}{}
	for responseName, fieldASTs := range p.Fields {

		resolved, hasNoFieldDefs, resolvedErrs := resolveField(p.ExecutionContext, p.ParentType, p.Source, fieldASTs)
		if hasNoFieldDefs {
			continue
		}
		for _, err := range resolvedErrs {
			result.Errors = append(result.Errors, err)
		}
		finalResults[responseName] = resolved
	}

	result.Data = finalResults
	resultChan <- &result
}

type CollectFieldsParams struct {
	ExeContext           ExecutionContext
	OperationType        *types.GraphQLObjectType
	SelectionSet         *ast.SelectionSet
	Fields               map[string][]*ast.Field
	VisitedFragmentNames map[string]bool
}

// Given a selectionSet, adds all of the fields in that selection to
// the passed in map of fields, and returns it at the end.
func collectFields(p CollectFieldsParams) map[string][]*ast.Field {

	fields := p.Fields
	if fields == nil {
		fields = map[string][]*ast.Field{}
	}
	if p.VisitedFragmentNames == nil {
		p.VisitedFragmentNames = map[string]bool{}
	}
	if p.SelectionSet == nil {
		return fields
	}
	for _, iSelection := range p.SelectionSet.Selections {
		switch selection := iSelection.(type) {
		case *ast.Field:
			if !shouldIncludeNode(p.ExeContext, selection.Directives) {
				continue
			}
			name := getFieldEntryKey(selection)
			if _, ok := fields[name]; !ok {
				fields[name] = []*ast.Field{}
			}
			fields[name] = append(fields[name], selection)
		case *ast.InlineFragment:

			if !shouldIncludeNode(p.ExeContext, selection.Directives) ||
				!doesFragmentConditionMatch(p.ExeContext, selection, p.OperationType) {
				continue
			}
			innerParams := CollectFieldsParams{
				ExeContext:           p.ExeContext,
				OperationType:        p.OperationType,
				SelectionSet:         selection.SelectionSet,
				Fields:               fields,
				VisitedFragmentNames: p.VisitedFragmentNames,
			}
			collectFields(innerParams)
		case *ast.FragmentSpread:
			fragName := ""
			if selection.Name != nil {
				fragName = selection.Name.Value
			}
			if visited, ok := p.VisitedFragmentNames[fragName]; (ok && visited) ||
				!shouldIncludeNode(p.ExeContext, selection.Directives) {
				continue
			}
			p.VisitedFragmentNames[fragName] = true
			fragment, hasFragment := p.ExeContext.Fragments[fragName]
			if !hasFragment {
				continue
			}

			if fragment, ok := fragment.(*ast.FragmentDefinition); ok {
				if !shouldIncludeNode(p.ExeContext, fragment.Directives) ||
					!doesFragmentConditionMatch(p.ExeContext, fragment, p.OperationType) {
					continue
				}
				innerParams := CollectFieldsParams{
					ExeContext:           p.ExeContext,
					OperationType:        p.OperationType,
					SelectionSet:         fragment.GetSelectionSet(),
					Fields:               fields,
					VisitedFragmentNames: p.VisitedFragmentNames,
				}
				collectFields(innerParams)
			}
		}
	}
	return fields
}

// Determines if a field should be included based on the @include and @skip
// directives, where @skip has higher precedence than @include.
func shouldIncludeNode(eCtx ExecutionContext, directives []*ast.Directive) bool {

	defaultReturnValue := true

	var skipAST *ast.Directive
	var includeAST *ast.Directive
	for _, directive := range directives {
		if directive == nil || directive.Name == nil {
			continue
		}
		if directive.Name.Value == types.GraphQLSkipDirective.Name {
			skipAST = directive
			break
		}
	}
	if skipAST != nil {
		argValues, err := getArgumentValues(
			types.GraphQLSkipDirective.Args,
			skipAST.Arguments,
			eCtx.VariableValues,
		)
		if err != nil {
			return defaultReturnValue
		}
		if skipIf, ok := argValues["if"]; ok {
			if boolSkipIf, ok := skipIf.(bool); ok {
				return !boolSkipIf
			}
		}
		return defaultReturnValue
	}
	for _, directive := range directives {
		if directive == nil || directive.Name == nil {
			continue
		}
		if directive.Name.Value == types.GraphQLIncludeDirective.Name {
			includeAST = directive
			break
		}
	}
	if includeAST != nil {
		argValues, err := getArgumentValues(
			types.GraphQLIncludeDirective.Args,
			includeAST.Arguments,
			eCtx.VariableValues,
		)
		if err != nil {
			return defaultReturnValue
		}
		if includeIf, ok := argValues["if"]; ok {
			if boolIncludeIf, ok := includeIf.(bool); ok {
				return boolIncludeIf
			}
		}
		return defaultReturnValue
	}
	return defaultReturnValue
}

// Determines if a fragment is applicable to the given type.
func doesFragmentConditionMatch(eCtx ExecutionContext, fragment ast.Node, ttype *types.GraphQLObjectType) bool {

	switch fragment := fragment.(type) {
	case *ast.FragmentDefinition:
		conditionalType, err := typeFromAST(eCtx.Schema, fragment.TypeCondition)
		if err != nil {
			return false
		}
		if conditionalType == ttype {
			return true
		}

		if conditionalType, ok := conditionalType.(types.GraphQLAbstractType); ok {
			return conditionalType.IsPossibleType(ttype)
		}
	case *ast.InlineFragment:
		conditionalType, err := typeFromAST(eCtx.Schema, fragment.TypeCondition)
		if err != nil {
			return false
		}
		if conditionalType == ttype {
			return true
		}

		if conditionalType, ok := conditionalType.(types.GraphQLAbstractType); ok {
			return conditionalType.IsPossibleType(ttype)
		}
	}

	return false
}

// Implements the logic to compute the key of a given field’s entry
func getFieldEntryKey(node *ast.Field) string {

	if node.Alias != nil && node.Alias.Value != "" {
		return node.Alias.Value
	}
	if node.Name != nil && node.Name.Value != "" {
		return node.Name.Value
	}
	return ""
}

/**
 * Resolves the field on the given source object. In particular, this
 * figures out the value that the field returns by calling its resolve function,
 * then calls completeValue to complete promises, serialize scalars, or execute
 * the sub-selection-set for objects.
 */
func resolveField(eCtx ExecutionContext, parentType *types.GraphQLObjectType, source interface{}, fieldASTs []*ast.Field) (result interface{}, hasNoFieldDefs bool, errs []graphqlerrors.GraphQLFormattedError) {
	// catch panic
	defer func() (interface{}, bool, []graphqlerrors.GraphQLFormattedError) {
		if r := recover(); r != nil {
			err := graphqlerrors.NewLocatedError(
				fmt.Sprintf("%v", r),
				graphqlerrors.FieldASTsToNodeASTs(fieldASTs),
			)
			errs = append(errs, graphqlerrors.FormatError(err))
			return result, hasNoFieldDefs, errs
		}
		return result, hasNoFieldDefs, errs
	}()

	hasNoFieldDefs = false

	fieldAST := fieldASTs[0]
	fieldName := ""
	if fieldAST.Name != nil {
		fieldName = fieldAST.Name.Value
	}

	fieldDef := getFieldDef(eCtx.Schema, parentType, fieldName)
	if fieldDef == nil {
		hasNoFieldDefs = true
		return nil, hasNoFieldDefs, errs
	}
	returnType := fieldDef.Type
	resolveFn := fieldDef.Resolve
	if resolveFn == nil {
		resolveFn = defaultResolveFn
	}

	// Build a map of arguments from the field.arguments AST, using the
	// variables scope to fulfill any variable references.
	// TODO: find a way to memoize, in case this field is within a List type.
	args, _ := getArgumentValues(fieldDef.Args, fieldAST.Arguments, eCtx.VariableValues)

	// The resolve function's optional third argument is a collection of
	// information about the current execution state.
	info := types.GraphQLResolveInfo{
		FieldName:      fieldName,
		FieldASTs:      fieldASTs,
		ReturnType:     returnType,
		ParentType:     parentType,
		Schema:         eCtx.Schema,
		Fragments:      eCtx.Fragments,
		RootValue:      eCtx.Root,
		Operation:      eCtx.Operation,
		VariableValues: eCtx.VariableValues,
	}

	// TODO: If an error occurs while calling the field `resolve` function, ensure that
	// it is wrapped as a GraphQLError with locations. Log this error and return
	// null if allowed, otherwise throw the error so the parent field can handle
	// it.
	result = resolveFn(types.GQLFRParams{
		Source: source,
		Args:   args,
		Info:   info,
	})
	completed, errs := completeValueCatchingError(eCtx, returnType, fieldASTs, info, result)
	return completed, hasNoFieldDefs, errs
}

func completeValueCatchingError(eCtx ExecutionContext, returnType types.GraphQLType, fieldASTs []*ast.Field, info types.GraphQLResolveInfo, result interface{}) (completed interface{}, errs []graphqlerrors.GraphQLFormattedError) {
	// catch panic
	defer func() (interface{}, []graphqlerrors.GraphQLFormattedError) {
		if r := recover(); r != nil {
			err := graphqlerrors.NewLocatedError(
				fmt.Sprintf("%v", r),
				graphqlerrors.FieldASTsToNodeASTs(fieldASTs),
			)
			errs = append(errs, graphqlerrors.FormatError(err))
			return completed, errs
		}
		return completed, errs
	}()

	if returnType, ok := returnType.(*types.GraphQLNonNull); ok {
		return completeValue(eCtx, returnType, fieldASTs, info, result)
	}
	completed, completedErrs := completeValue(eCtx, returnType, fieldASTs, info, result)
	if len(completedErrs) > 0 {
		return completed, append(errs, completedErrs...)
	}
	resultVal := reflect.ValueOf(completed)
	if resultVal.IsValid() && resultVal.Type().Kind() == reflect.Func {
		if propertyFn, ok := completed.(func() interface{}); ok {
			return propertyFn(), nil
		}
		return completed, append(errs, graphqlerrors.NewGraphQLFormattedError("Error resolving func. Expected `func() interface{}` signature"))
	}
	return completed, errs

}

func completeValue(eCtx ExecutionContext, returnType types.GraphQLType, fieldASTs []*ast.Field, info types.GraphQLResolveInfo, result interface{}) (interface{}, []graphqlerrors.GraphQLFormattedError) {

	var errs []graphqlerrors.GraphQLFormattedError
	// TODO: explore resolving go-routines in completeValue

	resultVal := reflect.ValueOf(result)
	if resultVal.IsValid() && resultVal.Type().Kind() == reflect.Func {
		if propertyFn, ok := result.(func() interface{}); ok {
			return propertyFn(), nil
		}
		return result, append(errs, graphqlerrors.NewGraphQLFormattedError("Error resolving func. Expected `func() interface{}` signature"))
	}

	if returnType, ok := returnType.(*types.GraphQLNonNull); ok {
		completed, completedErrs := completeValue(eCtx, returnType.OfType, fieldASTs, info, result)
		if len(completedErrs) > 0 {
			errs = append(errs, completedErrs...)
		}
		if completed == nil {
			err := graphqlerrors.NewLocatedError(
				fmt.Sprintf("Cannot return null for non-nullable field %v.%v.", info.ParentType, info.FieldName),
				graphqlerrors.FieldASTsToNodeASTs(fieldASTs),
			)
			errs = append(errs, graphqlerrors.FormatError(err))
		}
		return completed, errs
	}

	if isNullish(result) {
		return result, nil
	}

	// If field type is List, complete each item in the list with the inner type
	if returnType, ok := returnType.(*types.GraphQLList); ok {

		resultVal := reflect.ValueOf(result)
		err := invariant(
			resultVal.IsValid() && resultVal.Type().Kind() == reflect.Slice,
			"User Error: expected iterable, but did not find one.",
		)
		if err != nil {
			return result, append(errs, graphqlerrors.FormatError(err))
		}

		itemType := returnType.OfType
		completedResults := []interface{}{}
		for i := 0; i < resultVal.Len(); i++ {
			val := resultVal.Index(i).Interface()
			completedItem, completedErrs := completeValueCatchingError(eCtx, itemType, fieldASTs, info, val)
			if len(completedErrs) > 0 {
				errs = append(errs, completedErrs...)
				continue
			}
			completedResults = append(completedResults, completedItem)
		}
		return completedResults, errs
	}

	// If field type is Scalar or Enum, serialize to a valid value, returning
	// null if serialization is not possible.
	if returnType, ok := returnType.(*types.GraphQLScalarType); ok {
		err := invariant(returnType.Serialize != nil, "Missing serialize method on type")
		if err != nil {
			return result, append(errs, graphqlerrors.FormatError(err))
		}
		serializedResult := returnType.Serialize(result)
		if isNullish(serializedResult) {
			return nil, nil
		} else {
			return serializedResult, nil
		}
	}
	if returnType, ok := returnType.(*types.GraphQLEnumType); ok {
		err := invariant(returnType.Serialize != nil, "Missing serialize method on type")
		if err != nil {
			return result, append(errs, graphqlerrors.FormatError(err))
		}
		serializedResult := returnType.Serialize(result)
		if isNullish(serializedResult) {
			return nil, nil
		} else {
			return serializedResult, nil
		}
	}

	// Field type must be Object, Interface or Union and expect sub-selections.
	var objectType *types.GraphQLObjectType
	switch returnType := returnType.(type) {
	case *types.GraphQLObjectType:
		objectType = returnType
	case types.GraphQLAbstractType:
		objectType = returnType.GetObjectType(result, info)
		if objectType != nil && !returnType.IsPossibleType(objectType) {
			return result, append(errs, graphqlerrors.NewGraphQLFormattedError(
				fmt.Sprintf(`Runtime Object type "%v" is not a possible type `+
					`for "%v".`, objectType, returnType),
			))
		}
	}
	if objectType == nil {
		return nil, nil
	}

	// If there is an isTypeOf predicate function, call it with the
	// current result. If isTypeOf returns false, then raise an error rather
	// than continuing execution.
	if objectType.IsTypeOf != nil && !objectType.IsTypeOf(result, info) {
		return result, append(errs, graphqlerrors.NewGraphQLFormattedError(
			fmt.Sprintf(`Expected value of type "%v" but got: %T.`, objectType, result),
		))
	}

	// Collect sub-fields to execute to complete this value.
	subFieldASTs := map[string][]*ast.Field{}
	visitedFragmentNames := map[string]bool{}
	for _, fieldAST := range fieldASTs {
		if fieldAST == nil {
			continue
		}
		selectionSet := fieldAST.SelectionSet
		if selectionSet != nil {
			innerParams := CollectFieldsParams{
				ExeContext:           eCtx,
				OperationType:        objectType,
				SelectionSet:         selectionSet,
				Fields:               subFieldASTs,
				VisitedFragmentNames: visitedFragmentNames,
			}
			subFieldASTs = collectFields(innerParams)
		}
	}

	resultChannel := make(chan *types.GraphQLResult)
	executeFieldsParams := ExecuteFieldsParams{
		ExecutionContext: eCtx,
		ParentType:       objectType,
		Source:           result,
		Fields:           subFieldASTs,
	}
	go executeFields(executeFieldsParams, resultChannel)
	results := <-resultChannel
	for _, err := range results.Errors {
		errs = append(errs, err)
	}
	return results.Data, errs

}

func defaultResolveFn(p types.GQLFRParams) interface{} {
	// expects p.Source to be a map
	if sourceMap, ok := p.Source.(map[string]interface{}); ok {
		property := sourceMap[p.Info.FieldName]
		val := reflect.ValueOf(property)
		if val.IsValid() && val.Type().Kind() == reflect.Func {
			// try type casting the func to the most basic func signature
			// for more complex signatures, user have to define ResolveFn
			if propertyFn, ok := property.(func() interface{}); ok {
				return propertyFn()
			}
		}
		return property
	}
	return p.Source
}

func getFieldDef(schema types.GraphQLSchema, parentType *types.GraphQLObjectType, fieldName string) *types.GraphQLFieldDefinition {

	if parentType == nil {
		return nil
	}

	if fieldName == types.SchemaMetaFieldDef.Name &&
		schema.GetQueryType().Name == parentType.Name {
		return types.SchemaMetaFieldDef
	}
	if fieldName == types.TypeMetaFieldDef.Name &&
		schema.GetQueryType().Name == parentType.Name {
		return types.TypeMetaFieldDef
	}
	if fieldName == types.TypeNameMetaFieldDef.Name &&
		schema.GetQueryType().Name == parentType.Name {
		return types.TypeNameMetaFieldDef
	}
	return parentType.GetFields()[fieldName]
}
