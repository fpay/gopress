package gopress

import (
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

// NewTracingMiddleware wraps an http.Handler and traces incoming requests.
// Additionally, it adds the span to the request's context.
func NewTracingMiddleware(tr opentracing.Tracer, operationName string) MiddlewareFunc {
	return func(next HandlerFunc) HandlerFunc {
		return func(c Context) error {
			l := ContextLogger(c)
			// Try to join to a trace propagated in `req`.
			wireContext, err := tr.Extract(
				opentracing.TextMap,
				opentracing.HTTPHeadersCarrier(c.Request().Header),
			)
			if err != nil {
				l.Errorf("error encountered while trying to extract span: %+v\n", err)
			}
			// create span
			if span := tr.StartSpan(operationName, ext.RPCServerOption(wireContext)); span != nil {
				// store span in context
				r := c.Request()
				r = r.WithContext(opentracing.ContextWithSpan(r.Context(), span))
				c.SetRequest(r)
				defer span.Finish()
			}

			// perform handler
			if err := next(c); err != nil {
				c.Error(err)
			}
			return nil
		}
	}
}
