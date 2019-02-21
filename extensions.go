package graphql

import (
	"context"

	"github.com/graphql-go/graphql/gqlerrors"
)

// Extension is an interface for extensions in graphql
type Extension interface {
	// Init is used to help you initialize the extension
	Init(context.Context, *Params)
	// Name returns the name of the extension (make sure it's custom)
	Name() string
	// ParseDidStart ...
	ParseDidStart(context.Context)
	// ParseDidStart ...
	ParseEnded(context.Context, error)
	// ValidationDidStart ...
	ValidationDidStart(context.Context)
	// ValidationEnded ...
	ValidationEnded(context.Context, []gqlerrors.FormattedError)
	// ExecutionDidStart notifies about the start of the execution
	ExecutionDidStart(context.Context)
	// ExecutionEnded notifies about the end of the execution
	ExecutionEnded(context.Context)
	// ResolveFieldDidStart notifies about the start of the resolving of a field
	ResolveFieldDidStart(context.Context, *ResolveInfo)
	// ResolveFieldEnded notifies about the end of the resolving of a field
	ResolveFieldEnded(context.Context, *ResolveInfo)
	// HasResult returns if the extension wants to add data to the result
	HasResult() bool
	// GetResult returns the data that the extension wants to add to the result
	GetResult(context.Context) interface{}
}

// handleExtensionsInits handles all the init functions for all the extensions in the schema
func handleExtensionsInits(p *Params) {
	for _, ext := range p.Schema.extensions {
		ext.Init(p.Context, p)
	}
}

// handleExtensionsParseDidStart ...
func handleExtensionsParseDidStart(p *Params) {
	for _, ext := range p.Schema.extensions {
		ext.ParseDidStart(p.Context)
	}
}

// handleExtensionsParseEnded ...
func handleExtensionsParseEnded(p *Params, err error) {
	for _, ext := range p.Schema.extensions {
		ext.ParseEnded(p.Context, err)
	}
}

// handleExtensionsValidationDidStart ...
func handleExtensionsValidationDidStart(p *Params) {
	for _, ext := range p.Schema.extensions {
		ext.ValidationDidStart(p.Context)
	}
}

// handleExtensionsValidationEnded ...
func handleExtensionsValidationEnded(p *Params, errs []gqlerrors.FormattedError) {
	for _, ext := range p.Schema.extensions {
		ext.ValidationEnded(p.Context, errs)
	}
}

// handleExecutionDidStart handles the ExecutionDidStart functions
func handleExtensionsExecutionDidStart(p *ExecuteParams) {
	for _, ext := range p.Schema.extensions {
		ext.ExecutionDidStart(p.Context)
	}
}

// handleExecutionEnded handles the notification of the extensions about the end of the execution
func handleExtensionsExecutionEnded(p *ExecuteParams) {
	for _, ext := range p.Schema.extensions {
		ext.ExecutionEnded(p.Context)
	}
}

// handleResolveFieldDidStart handles the notification of the extensions about the start of a resolve function
func handleExtensionsResolveFieldDidStart(exts []Extension, p *executionContext, i *ResolveInfo) {
	for _, ext := range exts {
		ext.ResolveFieldDidStart(p.Context, i)
	}
}

// handleResolveFieldEnded handles the notification of the extensions about the end of a resolve function
func handleExtensionsResolveFieldEnded(exts []Extension, p *executionContext, i *ResolveInfo) {
	for _, ext := range exts {
		ext.ResolveFieldEnded(p.Context, i)
	}
}
