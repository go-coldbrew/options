package options

import "context"

// RequestContext holds both request-scoped options and log-scoped fields
// in a single context value, reducing context.WithValue allocations from 2 to 1.
// Both fields are eagerly initialized to avoid data races from concurrent lazy init.
type RequestContext struct {
	opts      *Options
	logFields *Options
}

var requestContextKey contextKey = "ColdBrewRequestContext"

// RequestContextFromContext retrieves the RequestContext from ctx.
// Returns nil if not present.
func RequestContextFromContext(ctx context.Context) *RequestContext {
	if rc, ok := ctx.Value(requestContextKey).(*RequestContext); ok {
		return rc
	}
	return nil
}

// getOrCreateRequestContext returns the existing RequestContext or creates one.
// Both opts and logFields are eagerly allocated so that concurrent access
// from multiple goroutines sharing the same context is safe without
// additional synchronization on the RequestContext itself.
func getOrCreateRequestContext(ctx context.Context) (context.Context, *RequestContext) {
	if rc := RequestContextFromContext(ctx); rc != nil {
		return ctx, rc
	}
	rc := &RequestContext{
		opts:      &Options{},
		logFields: &Options{},
	}
	return context.WithValue(ctx, requestContextKey, rc), rc
}

// Opts returns the Options for request-scoped key-value pairs.
func (rc *RequestContext) Opts() *Options {
	return rc.opts
}

// LogFields returns the Options used for log fields.
func (rc *RequestContext) LogFields() *Options {
	return rc.logFields
}

// AddToLogFields adds a key-value pair to the log fields stored in ctx.
// If ctx is nil, context.Background() is used.
func AddToLogFields(ctx context.Context, key string, value any) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if key == "" {
		return ctx
	}
	ctx, rc := getOrCreateRequestContext(ctx)
	rc.logFields.Add(key, value)
	return ctx
}

// LogFieldsFromContext retrieves the log fields Options from context.
// Returns nil if not present.
func LogFieldsFromContext(ctx context.Context) *Options {
	if ctx == nil {
		return nil
	}
	if rc := RequestContextFromContext(ctx); rc != nil {
		return rc.logFields
	}
	return nil
}
