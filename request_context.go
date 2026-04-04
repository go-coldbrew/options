package options

import "context"

// RequestContext holds both request-scoped options and log-scoped fields
// in a single context value, reducing context.WithValue allocations from 2 to 1.
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
func getOrCreateRequestContext(ctx context.Context) (context.Context, *RequestContext) {
	if rc := RequestContextFromContext(ctx); rc != nil {
		return ctx, rc
	}
	rc := &RequestContext{}
	return context.WithValue(ctx, requestContextKey, rc), rc
}

// Opts returns the Options, creating if needed.
func (rc *RequestContext) Opts() *Options {
	if rc.opts == nil {
		rc.opts = &Options{}
	}
	return rc.opts
}

// LogFields returns the Options used for log fields, creating if needed.
func (rc *RequestContext) LogFields() *Options {
	if rc.logFields == nil {
		rc.logFields = &Options{}
	}
	return rc.logFields
}

// AddToLogFields adds a key-value pair to the log fields stored in ctx.
// If ctx is nil, context.Background() is used.
func AddToLogFields(ctx context.Context, key string, value any) context.Context {
	if key == "" {
		return ctx
	}
	if ctx == nil {
		ctx = context.Background()
	}
	ctx, rc := getOrCreateRequestContext(ctx)
	rc.LogFields().Add(key, value)
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
