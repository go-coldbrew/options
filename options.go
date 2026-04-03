package options

import (
	"context"
	"sync"
)

// SupportPackageIsVersion1 is a compile-time assertion constant.
// Downstream packages reference this to enforce version compatibility.
const SupportPackageIsVersion1 = true

type contextKey string

var (
	optionsKey contextKey = "ColdBrewOptions"
)

// Options are request options passed from ColdBrew to server.
// Uses RWMutex + map instead of sync.Map since Options is per-request
// and never shared across goroutines.
type Options struct {
	mu sync.RWMutex
	m  map[string]any
}

// newOptions creates an Options with an initialized map.
func newOptions() *Options {
	return &Options{m: make(map[string]any, 2)}
}

// FromContext fetches options from provided context.
// If no options are found, it returns nil.
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
func AddToOptions(ctx context.Context, key string, value any) context.Context {
	h := FromContext(ctx)
	if h == nil {
		h = newOptions()
		ctx = context.WithValue(ctx, optionsKey, h)
	}
	if key != "" {
		h.Add(key, value)
	}
	return ctx
}

// Add adds a key-value pair to Options.
// Empty keys are silently ignored.
func (o *Options) Add(key string, value any) {
	if key == "" {
		return
	}
	o.mu.Lock()
	if o.m == nil {
		o.m = make(map[string]any, 2)
	}
	o.m[key] = value
	o.mu.Unlock()
}

// Del deletes an option by key.
func (o *Options) Del(key string) {
	o.mu.Lock()
	if o.m != nil {
		delete(o.m, key)
	}
	o.mu.Unlock()
}

// Get retrieves an option value by key.
func (o *Options) Get(key string) (any, bool) {
	if o == nil {
		return nil, false
	}
	o.mu.RLock()
	if o.m == nil {
		o.mu.RUnlock()
		return nil, false
	}
	v, found := o.m[key]
	o.mu.RUnlock()
	return v, found
}

// Store is a sync.Map-compatible alias for Add.
func (o *Options) Store(key, value any) {
	if k, ok := key.(string); ok {
		o.Add(k, value)
	}
}

// Load is a sync.Map-compatible alias for Get.
func (o *Options) Load(key any) (any, bool) {
	if k, ok := key.(string); ok {
		return o.Get(k)
	}
	return nil, false
}

// Delete is a sync.Map-compatible alias for Del.
func (o *Options) Delete(key any) {
	if k, ok := key.(string); ok {
		o.Del(k)
	}
}

// Range calls f sequentially for each key and value.
// If f returns false, Range stops the iteration.
// The callback may safely call Add/Del on the same Options instance.
func (o *Options) Range(f func(key, value any) bool) {
	o.mu.RLock()
	snapshot := make(map[string]any, len(o.m))
	for k, v := range o.m {
		snapshot[k] = v
	}
	o.mu.RUnlock()
	for k, v := range snapshot {
		if !f(k, v) {
			break
		}
	}
}
