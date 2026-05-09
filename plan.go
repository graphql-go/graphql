package graphql

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/graphql-go/graphql/gqlerrors"
	"github.com/graphql-go/graphql/language/ast"
)

// Plan is a precomputed execution shape for a (schema, document,
// operationName) triple. Same triple → same Plan; cache + reuse.
//
// PlanQuery does once-per-query work that Execute otherwise repeats
// every request: identify the operation, collect fragments, resolve
// every selection's *FieldDefinition, walk the selection tree to
// build per-level field lists, pre-coerce literal arguments, and
// pre-compute @include/@skip predicates that don't reference
// variables. ExecutePlan walks the plan and only does work that's
// inherently per-request: variable substitution, abstract type
// runtime resolution, resolver invocation, and result materialization.
//
// Plans are bound to the *Schema pointer they were planned against.
// If the schema is rebuilt, plans become stale (the *Schema pointer
// changes) and callers should re-plan.
type Plan struct {
	schema     *Schema
	operation  *ast.OperationDefinition
	fragments  map[string]ast.Definition
	rootType   *Object
	root       *selectionPlan
	isMutation bool
}

// selectionPlan is a pre-collected, source-ordered list of fields to
// emit for one selection set under a known parent runtime type. The
// plan tree mirrors the document, with one selectionPlan per object-
// returning field's sub-selection.
type selectionPlan struct {
	parentType *Object
	fields     []*fieldPlan
}

// fieldPlan is one entry in a selectionPlan: enough to resolve, run,
// and complete a single field without re-walking the schema or the
// AST. Sub-selections are pre-planned for object returns; abstract
// returns lazily plan per-concrete-type at execute time (cached on
// abstractAlternatives).
type fieldPlan struct {
	responseKey string
	fieldName   string
	fieldDef    *FieldDefinition
	fieldASTs   []*ast.Field // [0] is the canonical AST for arg lookup; full slice flows into ResolveInfo
	args        argPlan
	returnType  Output

	// skipPredicate evaluates the field's combined @skip / @include
	// directives against request variables. nil ⇒ always include
	// (constant-true at plan time, the common case).
	skipPredicate func(map[string]interface{}) bool

	// sub is set when returnType (after unwrapping NonNull and List)
	// resolves to a single concrete *Object; abstractAlternatives is
	// set when it resolves to an Interface or Union; both nil for
	// leaf fields (Scalar/Enum) and for fields whose sub-selection
	// the planner couldn't analyse (e.g. Object returns whose own
	// sub-selection contained inputs we don't yet handle — falls
	// back to runtime collectFields).
	sub                  *selectionPlan
	abstractAlternatives map[*Object]*selectionPlan
}

// argPlan separates static arg values (resolvable once, at plan
// time) from dynamic ones (referenced via variables). The common
// case — all literal args, or no args — yields a fully-built `static`
// map and skips getArgumentValues at execute time entirely.
type argPlan struct {
	// static holds the resolved arg map for the all-literals case.
	// Read-only; the executor copies before handing to the resolver
	// so the resolver can mutate without aliasing the plan.
	static map[string]interface{}

	// hasVariables forces the executor to fall through to the
	// runtime getArgumentValues path. This is correct but slower; the
	// goal is for the common literal-args case to skip it.
	hasVariables bool

	// fieldDefArgs / argASTs are kept for the dynamic fallback path;
	// nil when hasVariables is false.
	fieldDefArgs []*Argument
	argASTs      []*ast.Argument
}

// PlanQuery walks the document, picks the named operation (or the
// only one), pre-resolves fragments + the entire selection tree, and
// returns a Plan ready to be executed against the same schema.
//
// Errors are spec-aligned: missing operation name, unknown operation
// name, ambiguous operation, and document containing a non-operation
// non-fragment definition all surface here. Field-level errors (e.g.
// "no such field on type X") that the existing executor surfaces as
// dispatch-time errors continue to surface there — PlanQuery only
// fails on document-level errors.
func PlanQuery(schema *Schema, doc *ast.Document, operationName string) (*Plan, error) {
	if schema == nil {
		return nil, errors.New("graphql: PlanQuery: schema is nil")
	}
	if doc == nil {
		return nil, errors.New("graphql: PlanQuery: document is nil")
	}

	var operation *ast.OperationDefinition
	fragments := map[string]ast.Definition{}
	for _, definition := range doc.Definitions {
		switch d := definition.(type) {
		case *ast.OperationDefinition:
			if operationName == "" && operation != nil {
				return nil, errors.New("Must provide operation name if query contains multiple operations.")
			}
			if operationName == "" || (d.GetName() != nil && d.GetName().Value == operationName) {
				operation = d
			}
		case *ast.FragmentDefinition:
			key := ""
			if d.GetName() != nil && d.GetName().Value != "" {
				key = d.GetName().Value
			}
			fragments[key] = d
		default:
			return nil, fmt.Errorf("GraphQL cannot execute a request containing a %v", definition.GetKind())
		}
	}
	if operation == nil {
		if operationName != "" {
			return nil, fmt.Errorf(`Unknown operation named "%v".`, operationName)
		}
		return nil, errors.New("Must provide an operation.")
	}

	rootType, err := getOperationRootType(*schema, operation)
	if err != nil {
		return nil, err
	}

	plan := &Plan{
		schema:     schema,
		operation:  operation,
		fragments:  fragments,
		rootType:   rootType,
		isMutation: operation.GetOperation() == ast.OperationTypeMutation,
	}
	plan.root = plan.planSelectionSet(rootType, operation.GetSelectionSet(), nil)
	return plan, nil
}

// planSelectionSet pre-collects the fields under one selection-set
// for a known parent type, recursing into sub-selections. Two-phase:
//
//  1. collectInto walks the selection set + fragments, grouping
//     fields by responseKey and merging their ASTs.
//  2. For each grouped field, planFieldChildren walks the union of
//     all merged ASTs' sub-selections — ensuring fragments that name
//     the same response key contribute their sub-fields. Without
//     this, `... { x { a } } ... { x { b } }` would only see one
//     fragment's `x.{...}`.
//
// visitedFragmentNames is threaded along to avoid infinite recursion
// in mutually-referencing fragments — same shape as the runtime
// collectFields uses.
func (p *Plan) planSelectionSet(parentType *Object, selectionSet *ast.SelectionSet, visitedFragmentNames map[string]bool) *selectionPlan {
	if selectionSet == nil {
		return nil
	}
	if visitedFragmentNames == nil {
		visitedFragmentNames = map[string]bool{}
	}
	sp := &selectionPlan{parentType: parentType}
	keyed := map[string]int{}
	p.collectInto(parentType, selectionSet, visitedFragmentNames, sp, keyed, nil)
	if len(sp.fields) == 0 {
		return nil
	}
	// Phase 2: plan sub-selections for each merged field group.
	for _, fp := range sp.fields {
		if fp.fieldDef == nil {
			continue
		}
		p.planMergedFieldChildren(fp)
	}
	return sp
}

// planMergedFieldChildren walks every AST in fp.fieldASTs to merge
// their sub-selections into a single sub-plan (or per-concrete-type
// abstract alternatives). Mirrors completeObjectValue's loop over
// fieldASTs, but at plan time so the executor can use the result
// directly.
func (p *Plan) planMergedFieldChildren(fp *fieldPlan) {
	t := unwrapNamedType(fp.returnType)
	switch concrete := t.(type) {
	case *Object:
		fp.sub = p.planMergedSelectionsForType(concrete, fp.fieldASTs)
	case *Interface:
		possibleTypes := p.schema.PossibleTypes(concrete)
		fp.abstractAlternatives = make(map[*Object]*selectionPlan, len(possibleTypes))
		for _, pt := range possibleTypes {
			fp.abstractAlternatives[pt] = p.planMergedSelectionsForType(pt, fp.fieldASTs)
		}
	case *Union:
		possibleTypes := p.schema.PossibleTypes(concrete)
		fp.abstractAlternatives = make(map[*Object]*selectionPlan, len(possibleTypes))
		for _, pt := range possibleTypes {
			fp.abstractAlternatives[pt] = p.planMergedSelectionsForType(pt, fp.fieldASTs)
		}
	}
}

// planMergedSelectionsForType collects the union of every AST's
// SelectionSet under one concrete parent type, returning a
// selectionPlan that mirrors what completeObjectValue's runtime
// collectFields loop would produce.
func (p *Plan) planMergedSelectionsForType(parentType *Object, fieldASTs []*ast.Field) *selectionPlan {
	sp := &selectionPlan{parentType: parentType}
	keyed := map[string]int{}
	visited := map[string]bool{}
	for _, f := range fieldASTs {
		if f == nil || f.SelectionSet == nil {
			continue
		}
		p.collectInto(parentType, f.SelectionSet, visited, sp, keyed, nil)
	}
	if len(sp.fields) == 0 {
		return nil
	}
	for _, fp := range sp.fields {
		if fp.fieldDef == nil {
			continue
		}
		p.planMergedFieldChildren(fp)
	}
	return sp
}

// collectInto mirrors executor.collectFields: walks selections,
// follows fragment spreads + inline fragments, evaluates @include /
// @skip directives at plan time when constant. Per-field
// skipPredicates carry the dynamic part forward to ExecutePlan.
//
// parentPred carries variable-driven @skip / @include from any
// enclosing inline fragment or fragment spread. It is AND-composed
// with each field's own predicate when a new fieldPlan is created so
// that fragment-level gates are honored at execute time.
//
// keyed maps responseKey → index in sp.fields so repeat selections
// of the same response key merge their fieldASTs (matches
// collectFields's `fields[name] = append(fields[name], selection)`).
func (p *Plan) collectInto(parentType *Object, selectionSet *ast.SelectionSet, visitedFragmentNames map[string]bool, sp *selectionPlan, keyed map[string]int, parentPred func(map[string]interface{}) bool) {
	for _, iSelection := range selectionSet.Selections {
		switch sel := iSelection.(type) {
		case *ast.Field:
			pred, alwaysSkip := planDirectives(sel.Directives)
			if alwaysSkip {
				continue
			}
			responseKey := getFieldEntryKey(sel)
			if idx, ok := keyed[responseKey]; ok {
				// Merge with an earlier same-key field (sub-selection
				// merging happens at execute time via collectFields on
				// the sub-selection — for now we just keep all ASTs and
				// let the runtime path stitch sub-selections; the
				// plan-time precompute conservatively re-plans the
				// first AST's sub-selection, which is correct because
				// validation rules guarantee mergeable selections refer
				// to the same field).
				sp.fields[idx].fieldASTs = append(sp.fields[idx].fieldASTs, sel)
				continue
			}
			fieldName := ""
			if sel.Name != nil {
				fieldName = sel.Name.Value
			}
			fieldDef := getFieldDef(*p.schema, parentType, fieldName)
			if fieldDef == nil {
				// Unknown field: keep it in the plan with a nil
				// fieldDef so ExecutePlan can mirror the
				// hasNoFieldDefs branch (skip the response key).
			}
			fp := &fieldPlan{
				responseKey:   responseKey,
				fieldName:     fieldName,
				fieldDef:      fieldDef,
				fieldASTs:     []*ast.Field{sel},
				skipPredicate: andPredicates(parentPred, pred),
			}
			if fieldDef != nil {
				fp.returnType = fieldDef.Type
				fp.args = planArguments(fieldDef.Args, sel.Arguments)
			}
			keyed[responseKey] = len(sp.fields)
			sp.fields = append(sp.fields, fp)

		case *ast.InlineFragment:
			pred, alwaysSkip := planDirectives(sel.Directives)
			if alwaysSkip {
				continue
			}
			if !planFragmentMatches(*p.schema, sel.TypeCondition, parentType) {
				continue
			}
			if sel.SelectionSet != nil {
				p.collectInto(parentType, sel.SelectionSet, visitedFragmentNames, sp, keyed, andPredicates(parentPred, pred))
			}

		case *ast.FragmentSpread:
			pred, alwaysSkip := planDirectives(sel.Directives)
			if alwaysSkip {
				continue
			}
			fragName := ""
			if sel.Name != nil {
				fragName = sel.Name.Value
			}
			if visitedFragmentNames[fragName] {
				continue
			}
			frag, ok := p.fragments[fragName]
			if !ok {
				continue
			}
			fragDef, ok := frag.(*ast.FragmentDefinition)
			if !ok {
				continue
			}
			visitedFragmentNames[fragName] = true
			if !planFragmentMatches(*p.schema, fragDef.TypeCondition, parentType) {
				continue
			}
			if fragDef.GetSelectionSet() != nil {
				p.collectInto(parentType, fragDef.GetSelectionSet(), visitedFragmentNames, sp, keyed, andPredicates(parentPred, pred))
			}
		}
	}
}

// andPredicates returns a predicate that is true only when both inputs
// are true. nil is treated as the constant-true predicate, so the
// common "no enclosing gate" / "no field-level directive" cases avoid
// allocating a closure.
func andPredicates(a, b func(map[string]interface{}) bool) func(map[string]interface{}) bool {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	return func(vars map[string]interface{}) bool {
		if !a(vars) {
			return false
		}
		return b(vars)
	}
}

// unwrapNamedType peels NonNull and List wrappers to expose the
// underlying named type (Object / Interface / Union / Scalar / Enum).
func unwrapNamedType(t Output) Output {
	for {
		switch tt := t.(type) {
		case *NonNull:
			t = tt.OfType.(Output)
		case *List:
			t = tt.OfType.(Output)
		default:
			return tt
		}
	}
}

// planArguments separates static arguments (all-literals → fully
// resolvable at plan time) from dynamic ones (any variable reference
// → defer to per-request getArgumentValues).
//
// We deliberately don't try to do partial pre-coercion (e.g.
// literal-args mixed with one variable arg → coerce the literals
// only) in this first cut: the all-or-nothing split is correct,
// trivial, and captures the common case (the bench query has all
// literal args).
func planArguments(argDefs []*Argument, argASTs []*ast.Argument) argPlan {
	if len(argDefs) == 0 && len(argASTs) == 0 {
		return argPlan{}
	}
	if astHasVariables(argASTs) {
		return argPlan{
			hasVariables: true,
			fieldDefArgs: argDefs,
			argASTs:      argASTs,
		}
	}
	// All-literal: coerce once. Pass nil variableValues — guaranteed
	// no variable refs.
	static := getArgumentValues(argDefs, argASTs, nil)
	if len(static) == 0 {
		return argPlan{}
	}
	return argPlan{static: static}
}

// astHasVariables walks the argument AST tree looking for any
// ast.Variable node. Returns true on the first hit.
func astHasVariables(argASTs []*ast.Argument) bool {
	for _, a := range argASTs {
		if a == nil {
			continue
		}
		if valueHasVariables(a.Value) {
			return true
		}
	}
	return false
}

func valueHasVariables(v ast.Value) bool {
	switch n := v.(type) {
	case nil:
		return false
	case *ast.Variable:
		return true
	case *ast.ListValue:
		for _, item := range n.Values {
			if valueHasVariables(item) {
				return true
			}
		}
	case *ast.ObjectValue:
		for _, f := range n.Fields {
			if f != nil && valueHasVariables(f.Value) {
				return true
			}
		}
	}
	return false
}

// planDirectives evaluates @include and @skip directives at plan
// time when their `if` argument is a literal; returns a
// skipPredicate (nil if always-include) and an alwaysSkip flag (true
// if literal evaluation produced a definitive skip).
func planDirectives(directives []*ast.Directive) (pred func(map[string]interface{}) bool, alwaysSkip bool) {
	var skipDir, includeDir *ast.Directive
	for _, d := range directives {
		if d == nil || d.Name == nil {
			continue
		}
		switch d.Name.Value {
		case SkipDirective.Name:
			skipDir = d
		case IncludeDirective.Name:
			includeDir = d
		}
	}
	if skipDir == nil && includeDir == nil {
		return nil, false
	}
	// Evaluate constants where possible; surface a runtime predicate
	// for the variable-driven cases.
	var skipDyn, includeDyn *ast.Directive
	if skipDir != nil {
		if astHasVariables(skipDir.Arguments) {
			skipDyn = skipDir
		} else {
			vals := getArgumentValues(SkipDirective.Args, skipDir.Arguments, nil)
			if v, ok := vals["if"].(bool); ok && v {
				return nil, true
			}
		}
	}
	if includeDir != nil {
		if astHasVariables(includeDir.Arguments) {
			includeDyn = includeDir
		} else {
			vals := getArgumentValues(IncludeDirective.Args, includeDir.Arguments, nil)
			if v, ok := vals["if"].(bool); ok && !v {
				return nil, true
			}
		}
	}
	if skipDyn == nil && includeDyn == nil {
		return nil, false
	}
	return func(vars map[string]interface{}) bool {
		if skipDyn != nil {
			vals := getArgumentValues(SkipDirective.Args, skipDyn.Arguments, vars)
			if v, ok := vals["if"].(bool); ok && v {
				return false // excluded
			}
		}
		if includeDyn != nil {
			vals := getArgumentValues(IncludeDirective.Args, includeDyn.Arguments, vars)
			if v, ok := vals["if"].(bool); ok && !v {
				return false // excluded
			}
		}
		return true
	}, false
}

// planFragmentMatches mirrors doesFragmentConditionMatch: a missing
// type condition matches anything; otherwise the condition resolves
// against the schema and must equal — or, for abstract types,
// include — the runtime parent type.
func planFragmentMatches(schema Schema, typeConditionAST *ast.Named, runtime *Object) bool {
	if typeConditionAST == nil {
		return true
	}
	conditionalType, err := typeFromAST(schema, typeConditionAST)
	if err != nil {
		return false
	}
	if conditionalType == runtime {
		return true
	}
	if conditionalType.Name() == runtime.Name() {
		return true
	}
	switch ct := conditionalType.(type) {
	case *Interface:
		return schema.IsPossibleType(ct, runtime)
	case *Union:
		return schema.IsPossibleType(ct, runtime)
	}
	return false
}

// ExecutePlan runs a planned operation. Args / Root / Context still
// flow in per-request; everything else (selection shape, field
// resolution, literal args, directive predicates, sub-plans) is
// taken from the plan.
//
// The walker mirrors executor.executeOperation → executeFields →
// resolveField → completeValue, but skips collectFields and
// getFieldDef on the hot path. Per-field arguments come from the
// argPlan: static (no variables) bypasses getArgumentValues entirely.
func ExecutePlan(plan *Plan, p ExecuteParams) (result *Result) {
	if plan == nil {
		return &Result{Errors: gqlerrors.FormatErrors(errors.New("graphql: ExecutePlan: plan is nil"))}
	}
	ctx := p.Context
	if ctx == nil {
		ctx = context.Background()
	}

	extErrs, executionFinishFn := handleExtensionsExecutionDidStart(&p)
	if len(extErrs) != 0 {
		return &Result{Errors: extErrs}
	}
	defer func() {
		extErrs := executionFinishFn(result)
		if len(extErrs) != 0 {
			result.Errors = append(result.Errors, extErrs...)
		}
		addExtensionResults(&p, result)
	}()

	resultChannel := make(chan *Result, 2)
	go func() {
		out := &Result{}
		defer func() {
			if err := recover(); err != nil {
				if e, ok := err.(error); ok {
					out.Errors = append(out.Errors, gqlerrors.FormatError(e))
				} else {
					out.Errors = append(out.Errors, gqlerrors.FormatError(fmt.Errorf("%v", err)))
				}
			}
			resultChannel <- out
		}()

		// Plan is bound to plan.schema (sub-plans, abstractAlternatives,
		// and field defs were resolved against it). Use that same schema
		// here so variable coercion and abstract-type resolution stay
		// consistent with what the plan was built against — p.Schema is
		// ignored to avoid silent drift if the caller passes a rebuilt
		// schema with the same shape but different *Object pointers.
		execSchema := *plan.schema
		variableValues, err := getVariableValues(execSchema, plan.operation.GetVariableDefinitions(), p.Args)
		if err != nil {
			out.Errors = append(out.Errors, gqlerrors.FormatError(err))
			return
		}

		eCtx := &executionContext{
			Schema:         execSchema,
			Fragments:      plan.fragments,
			Root:           p.Root,
			Operation:      plan.operation,
			VariableValues: variableValues,
			Context:        ctx,
		}

		data := executePlannedSelection(eCtx, plan.root, p.Root, plan.rootType, nil)
		// Mutations run serially with each field's result
		// dethunked depth-first; queries run all then dethunk
		// breadth-first. The traversal here just runs the appropriate
		// dethunker on the assembled map.
		if plan.isMutation {
			dethunkMapDepthFirst(data)
		} else {
			dethunkMapWithBreadthFirstTraversal(data)
		}
		out.Data = data
		out.Errors = append(out.Errors, eCtx.Errors...)
	}()

	select {
	case <-ctx.Done():
		r := &Result{}
		r.Errors = append(r.Errors, gqlerrors.FormatError(ctx.Err()))
		return r
	case r := <-resultChannel:
		return r
	}
}

// executePlannedSelection runs one selection plan against a parent
// source value, returning the assembled response map. Walks fields
// in source order (sp.fields is built that way at plan time).
//
// Mutation vs. query traversal is handled at the top level in
// ExecutePlan via dethunkMapDepthFirst / dethunkMapWithBreadthFirstTraversal,
// so this walker is the same for both.
func executePlannedSelection(eCtx *executionContext, sp *selectionPlan, source interface{}, parentType *Object, path *ResponsePath) map[string]interface{} {
	if sp == nil {
		return map[string]interface{}{}
	}
	if source == nil {
		source = map[string]interface{}{}
	}
	finalResults := make(map[string]interface{}, len(sp.fields))
	for _, fp := range sp.fields {
		if fp.skipPredicate != nil && !fp.skipPredicate(eCtx.VariableValues) {
			continue
		}
		if fp.fieldDef == nil {
			// Mirrors executeSubFields' hasNoFieldDefs branch: silently
			// skip unknown fields. Validation should have rejected
			// these but we match runtime behavior for safety.
			continue
		}
		fieldPath := path.WithKey(fp.responseKey)
		resolved, ok := resolvePlannedField(eCtx, parentType, source, fp, fieldPath)
		if !ok {
			continue
		}
		finalResults[fp.responseKey] = resolved
	}
	return finalResults
}

// resolvePlannedField mirrors resolveField but uses the plan's
// pre-resolved fieldDef + pre-coerced static args + pre-decided
// returnType. Pure-literal arguments skip getArgumentValues entirely;
// variable-bearing args fall through to the existing per-request
// coercion.
func resolvePlannedField(eCtx *executionContext, parentType *Object, source interface{}, fp *fieldPlan, path *ResponsePath) (result interface{}, ok bool) {
	var returnType Output
	defer func() {
		if r := recover(); r != nil {
			handleFieldError(r, FieldASTsToNodeASTs(fp.fieldASTs), path, returnType, eCtx)
			ok = true
		}
	}()
	fieldDef := fp.fieldDef
	returnType = fp.returnType
	resolveFn := fieldDef.Resolve
	if resolveFn == nil {
		resolveFn = DefaultResolveFn
	}

	// Resolvers expect a non-nil Args map (the existing resolveField
	// path always passes the result of getArgumentValues, which is
	// never nil even when empty). Match that contract.
	var args map[string]interface{}
	switch {
	case fp.args.hasVariables:
		args = getArgumentValues(fp.args.fieldDefArgs, fp.args.argASTs, eCtx.VariableValues)
	case fp.args.static != nil:
		// Resolvers may mutate the args map; copy to keep the plan's
		// static map immutable across requests.
		args = make(map[string]interface{}, len(fp.args.static))
		for k, v := range fp.args.static {
			args[k] = v
		}
	default:
		args = map[string]interface{}{}
	}

	info := ResolveInfo{
		FieldName:      fp.fieldName,
		FieldASTs:      fp.fieldASTs,
		Path:           path,
		ReturnType:     returnType,
		ParentType:     parentType,
		Schema:         eCtx.Schema,
		Fragments:      eCtx.Fragments,
		RootValue:      eCtx.Root,
		Operation:      eCtx.Operation,
		VariableValues: eCtx.VariableValues,
	}

	extErrs, resolveFieldFinishFn := handleExtensionsResolveFieldDidStart(eCtx.Schema.extensions, eCtx, &info)
	if len(extErrs) != 0 {
		eCtx.Errors = append(eCtx.Errors, extErrs...)
	}

	var resolveFnError error
	result, resolveFnError = resolveFn(ResolveParams{
		Source:  source,
		Args:    args,
		Info:    info,
		Context: eCtx.Context,
	})

	extErrs = resolveFieldFinishFn(result, resolveFnError)
	if len(extErrs) != 0 {
		eCtx.Errors = append(eCtx.Errors, extErrs...)
	}
	if resolveFnError != nil {
		panic(resolveFnError)
	}

	completed := completePlannedValueCatchingError(eCtx, returnType, fp, info, path, result)
	return completed, true
}

func completePlannedValueCatchingError(eCtx *executionContext, returnType Type, fp *fieldPlan, info ResolveInfo, path *ResponsePath, result interface{}) (completed interface{}) {
	defer func() {
		if r := recover(); r != nil {
			handleFieldError(r, FieldASTsToNodeASTs(fp.fieldASTs), path, returnType, eCtx)
		}
	}()
	if rt, ok := returnType.(*NonNull); ok {
		return completePlannedValue(eCtx, rt, fp, info, path, result)
	}
	return completePlannedValue(eCtx, returnType, fp, info, path, result)
}

func completePlannedValue(eCtx *executionContext, returnType Type, fp *fieldPlan, info ResolveInfo, path *ResponsePath, result interface{}) interface{} {
	resultVal := reflect.ValueOf(result)
	if resultVal.IsValid() && resultVal.Kind() == reflect.Func {
		return func() interface{} {
			return completePlannedThunkValueCatchingError(eCtx, returnType, fp, info, path, result)
		}
	}
	if rt, ok := returnType.(*NonNull); ok {
		completed := completePlannedValue(eCtx, rt.OfType, fp, info, path, result)
		if completed == nil {
			err := NewLocatedErrorWithPath(
				fmt.Sprintf("Cannot return null for non-nullable field %v.%v.", info.ParentType, info.FieldName),
				FieldASTsToNodeASTs(fp.fieldASTs),
				path.AsArray(),
			)
			panic(gqlerrors.FormatError(err))
		}
		return completed
	}
	if isNullish(result) {
		return nil
	}
	if rt, ok := returnType.(*List); ok {
		return completePlannedListValue(eCtx, rt, fp, info, path, result)
	}
	if rt, ok := returnType.(*Scalar); ok {
		return completeLeafValue(rt, result)
	}
	if rt, ok := returnType.(*Enum); ok {
		return completeLeafValue(rt, result)
	}
	if rt, ok := returnType.(*Union); ok {
		return completePlannedAbstractValue(eCtx, rt, fp, info, path, result)
	}
	if rt, ok := returnType.(*Interface); ok {
		return completePlannedAbstractValue(eCtx, rt, fp, info, path, result)
	}
	if rt, ok := returnType.(*Object); ok {
		return completePlannedObjectValue(eCtx, rt, fp, info, path, result)
	}
	err := invariantf(false, `Cannot complete value of unexpected type "%v."`, returnType)
	if err != nil {
		panic(gqlerrors.FormatError(err))
	}
	return nil
}

func completePlannedThunkValueCatchingError(eCtx *executionContext, returnType Type, fp *fieldPlan, info ResolveInfo, path *ResponsePath, result interface{}) (completed interface{}) {
	defer func() {
		if r := recover(); r != nil {
			handleFieldError(r, FieldASTsToNodeASTs(fp.fieldASTs), path, returnType, eCtx)
		}
	}()
	propertyFn, ok := result.(func() (interface{}, error))
	if !ok {
		err := gqlerrors.NewFormattedError("Error resolving func. Expected `func() (interface{}, error)` signature")
		panic(gqlerrors.FormatError(err))
	}
	fnResult, err := propertyFn()
	if err != nil {
		panic(gqlerrors.FormatError(err))
	}
	result = fnResult
	if rt, ok := returnType.(*NonNull); ok {
		return completePlannedValue(eCtx, rt, fp, info, path, result)
	}
	return completePlannedValue(eCtx, returnType, fp, info, path, result)
}

func completePlannedListValue(eCtx *executionContext, returnType *List, fp *fieldPlan, info ResolveInfo, path *ResponsePath, result interface{}) interface{} {
	resultVal := reflect.ValueOf(result)
	if resultVal.Kind() == reflect.Ptr {
		resultVal = resultVal.Elem()
	}
	parentTypeName := ""
	if info.ParentType != nil {
		parentTypeName = info.ParentType.Name()
	}
	err := invariantf(
		resultVal.IsValid() && isIterable(result),
		"User Error: expected iterable, but did not find one "+
			"for field %v.%v.", parentTypeName, info.FieldName)
	if err != nil {
		panic(gqlerrors.FormatError(err))
	}
	itemType := returnType.OfType
	completedResults := make([]interface{}, 0, resultVal.Len())
	for i := 0; i < resultVal.Len(); i++ {
		val := resultVal.Index(i).Interface()
		fieldPath := path.WithKey(i)
		completedItem := completePlannedValueCatchingError(eCtx, itemType, fp, info, fieldPath, val)
		completedResults = append(completedResults, completedItem)
	}
	return completedResults
}

func completePlannedObjectValue(eCtx *executionContext, returnType *Object, fp *fieldPlan, info ResolveInfo, path *ResponsePath, result interface{}) interface{} {
	if returnType.IsTypeOf != nil {
		p := IsTypeOfParams{Value: result, Info: info, Context: eCtx.Context}
		if !returnType.IsTypeOf(p) {
			panic(gqlerrors.NewFormattedError(
				fmt.Sprintf(`Expected value of type "%v" but got: %T.`, returnType, result),
			))
		}
	}
	if fp.sub != nil {
		return executePlannedSelection(eCtx, fp.sub, result, returnType, path)
	}
	// Fallback: planner didn't precompute (e.g. selection set was
	// empty per validation, which shouldn't reach here for object
	// types). Surface no-data with a defensive empty map.
	return map[string]interface{}{}
}

func completePlannedAbstractValue(eCtx *executionContext, returnType Abstract, fp *fieldPlan, info ResolveInfo, path *ResponsePath, result interface{}) interface{} {
	var runtimeType *Object
	rtParams := ResolveTypeParams{Value: result, Info: info, Context: eCtx.Context}
	if u, ok := returnType.(*Union); ok && u.ResolveType != nil {
		runtimeType = u.ResolveType(rtParams)
	} else if i, ok := returnType.(*Interface); ok && i.ResolveType != nil {
		runtimeType = i.ResolveType(rtParams)
	} else {
		runtimeType = defaultResolveTypeFn(rtParams, returnType)
	}
	if err := invariantf(runtimeType != nil,
		`Abstract type %v must resolve to an Object type at runtime `+
			`for field %v.%v with value "%v", received "%v".`,
		returnType, info.ParentType, info.FieldName, result, runtimeType,
	); err != nil {
		panic(err)
	}
	if !eCtx.Schema.IsPossibleType(returnType, runtimeType) {
		panic(gqlerrors.NewFormattedError(
			fmt.Sprintf(`Runtime Object type "%v" is not a possible type for "%v".`, runtimeType, returnType),
		))
	}
	if sub, ok := fp.abstractAlternatives[runtimeType]; ok && sub != nil {
		return executePlannedSelection(eCtx, sub, result, runtimeType, path)
	}
	// Defensive fallback: unplanned concrete type (e.g. interface
	// gained a new implementer between plan time and execute time —
	// shouldn't happen in practice since schema rebuilds invalidate
	// the plan, but stay correct).
	return map[string]interface{}{}
}
