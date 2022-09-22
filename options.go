package options

import (
	"context"
	"strings"
	"sync"
)

type contextKey string

var (
	optionsKey contextKey = "ColdBrewOptions"
)

// Options are request options passed from Orion to server
type Options struct {
	sync.Map
}

// FromContext fetchs options from provided context
func FromContext(ctx context.Context) *Options {
	if h := ctx.Value(optionsKey); h != nil {
		if options, ok := h.(*Options); ok {
			return options
		}
	}
	return nil
}

// AddToOptions adds options to context
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
func (o Options) Add(key string, value interface{}) {
	o.Add(key, value)
}

// Del an options
func (o Options) Del(key string) {
	o.Delete(key)
}

// Get an options
func (o Options) Get(key string) (interface{}, bool) {
	value, found := o.Load(strings.ToLower(key))
	return value, found
}
