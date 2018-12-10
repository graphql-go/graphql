package graphql

import (
	"context"

	"github.com/graphql-go/graphql/gqlerrors"
)

type TraceQueryFinishFunc func([]gqlerrors.FormattedError)
type TraceFieldFinishFunc func(gqlerrors.FormattedError)

type Tracer interface {
	TraceQuery(ctx context.Context, queryString string, operationName string) (context.Context, TraceQueryFinishFunc) 
	TraceField(ctx context.Context, fieldName string) TraceFieldFinishFunc
}
