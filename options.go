package options

import (
	"context"
	"sync"
)

type contextKey string

var (
	optionsKey contextKey = "ColdBrewOptions"
)

// Options are request options passed from ColdBrew to server
type Options struct {
	sync.Map
}

// FromContext fetchs options from provided context
// if no options found, return nil
func FromContext(ctx context.Context) *Options {
	if h := ctx.Value(optionsKey); h != nil {
		if options, ok := h.(*Options); ok {
			return options
		}
	}
	return nil
}

// AddToOptions adds options to context
// if no options found, create a new one and adds the provided options to it and returns the new context
func AddToOptions(ctx context.Context, key string, value interface{}) context.Context {
	h := FromContext(ctx)
	if h == nil {
		ctx = context.WithValue(ctx, optionsKey, new(Options))
		h = FromContext(ctx)
	}
	if h != nil && key != "" {
		h.Add(key, value)
	}
	return ctx
}

// Add to Options
// can be used to add options to context
func (o *Options) Add(key string, value interface{}) {
	o.Store(key, value)
}

// Del an options
// can be used to delete options from context
func (o *Options) Del(key string) {
	o.Delete(key)
}

// Get an options
// can be used to get options from context
func (o *Options) Get(key string) (interface{}, bool) {
	if o == nil {
		return nil, false
	}
	value, found := o.Load(key)
	return value, found
}
