package monitoring

import (
	"github.com/graphql-go/graphql"
	"go.opentelemetry.io/otel"
)

// Middleware to trace resolver functions
func TraceResolver(resolverFunc func(p graphql.ResolveParams) (interface{}, error)) graphql.FieldResolveFn {
	return func(p graphql.ResolveParams) (interface{}, error) {
		// Create a new span for the resolver function
		spanName := p.Info.ParentType.Name() + "." + p.Info.FieldName
		ctx, span := otel.GetTracerProvider().Tracer("").Start(p.Context, spanName)
		defer span.End()

		// Perform the resolver logic
		return resolverFunc(graphql.ResolveParams{Source: p.Source, Args: p.Args, Info: p.Info, Context: ctx})
	}
}
