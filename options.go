package options

import (
	"context"
	"maps"
	"sync"
)

// SupportPackageIsVersion1 is a compile-time assertion constant.
// Downstream packages reference this to enforce version compatibility.
const SupportPackageIsVersion1 = true

type contextKey string

// Options are request options passed from ColdBrew to server.
// Uses RWMutex + map instead of sync.Map since Options is per-request
// and never shared across goroutines.
type Options struct {
	mu sync.RWMutex
	m  map[string]any
}

// FromContext fetches options from provided context.
// If no options are found, it returns nil.
func FromContext(ctx context.Context) *Options {
	if rc := RequestContextFromContext(ctx); rc != nil {
		return rc.opts
	}
	return nil
}

// AddToOptions adds a key-value pair to the Options stored in ctx.
// If no Options exists in the context, a new one is created.
// Empty keys are silently ignored and do not allocate.
func AddToOptions(ctx context.Context, key string, value any) context.Context {
	if key == "" {
		return ctx
	}
	ctx, rc := getOrCreateRequestContext(ctx)
	rc.Opts().Add(key, value)
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
// Only string keys are supported; non-string keys are silently ignored.
func (o *Options) Store(key, value any) {
	if k, ok := key.(string); ok {
		o.Add(k, value)
	}
}

// Load is a sync.Map-compatible alias for Get.
// Only string keys are supported; non-string keys return (nil, false).
func (o *Options) Load(key any) (any, bool) {
	if k, ok := key.(string); ok {
		return o.Get(k)
	}
	return nil, false
}

// Delete is a sync.Map-compatible alias for Del.
// Only string keys are supported; non-string keys are silently ignored.
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
	if len(o.m) == 0 {
		o.mu.RUnlock()
		return
	}
	snapshot := make(map[string]any, len(o.m))
	maps.Copy(snapshot, o.m)
	o.mu.RUnlock()
	for k, v := range snapshot {
		if !f(k, v) {
			break
		}
	}
}

// RangeSlice calls f sequentially for each key and value, using a slice
// snapshot. This is more efficient than Range for small maps and matches
// the iteration pattern used by LogFields.
func (o *Options) RangeSlice(f func(key, value any) bool) {
	o.mu.RLock()
	if len(o.m) == 0 {
		o.mu.RUnlock()
		return
	}
	type kv struct {
		k string
		v any
	}
	entries := make([]kv, 0, len(o.m))
	for k, v := range o.m {
		entries = append(entries, kv{k, v})
	}
	o.mu.RUnlock()
	for _, e := range entries {
		if !f(e.k, e.v) {
			break
		}
	}
}
