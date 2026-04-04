package options

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestContext_SingleContextValue(t *testing.T) {
	ctx := context.Background()

	// Adding options creates a RequestContext
	ctx = AddToOptions(ctx, "opt-key", "opt-val")
	rc := RequestContextFromContext(ctx)
	require.NotNil(t, rc, "RequestContext should exist after AddToOptions")

	// Adding log fields reuses the same RequestContext
	ctx = AddToLogFields(ctx, "log-key", "log-val")
	rc2 := RequestContextFromContext(ctx)
	assert.Same(t, rc, rc2, "should reuse same RequestContext")

	// Both are accessible
	opts := FromContext(ctx)
	require.NotNil(t, opts)
	v, ok := opts.Get("opt-key")
	assert.True(t, ok)
	assert.Equal(t, "opt-val", v)

	lf := LogFieldsFromContext(ctx)
	require.NotNil(t, lf)
	v, ok = lf.Get("log-key")
	assert.True(t, ok)
	assert.Equal(t, "log-val", v)
}

func TestRequestContext_LogFieldsFirst(t *testing.T) {
	ctx := context.Background()

	// Log fields first, then options
	ctx = AddToLogFields(ctx, "log-key", "log-val")
	ctx = AddToOptions(ctx, "opt-key", "opt-val")

	rc := RequestContextFromContext(ctx)
	require.NotNil(t, rc)
	assert.NotNil(t, rc.opts)
	assert.NotNil(t, rc.logFields)
}

func TestRequestContext_LazyInit(t *testing.T) {
	ctx := context.Background()
	ctx = AddToOptions(ctx, "k", "v")

	rc := RequestContextFromContext(ctx)
	require.NotNil(t, rc)
	assert.NotNil(t, rc.opts, "opts should be initialized")
	assert.Nil(t, rc.logFields, "logFields should be nil until used")
}

func TestRequestContext_EmptyContext(t *testing.T) {
	ctx := context.Background()
	assert.Nil(t, RequestContextFromContext(ctx))
	assert.Nil(t, FromContext(ctx))
	assert.Nil(t, LogFieldsFromContext(ctx))
}

func TestLogFieldsFromContext_NilCtx(t *testing.T) {
	assert.Nil(t, LogFieldsFromContext(nil))
}

func TestAddToLogFields_NilCtx(t *testing.T) {
	ctx := AddToLogFields(nil, "key", "val")
	assert.NotNil(t, ctx, "should create context from nil")
	lf := LogFieldsFromContext(ctx)
	require.NotNil(t, lf)
	v, ok := lf.Get("key")
	assert.True(t, ok)
	assert.Equal(t, "val", v)
}

func TestAddToLogFields_EmptyKey(t *testing.T) {
	ctx := context.Background()
	ctx = AddToLogFields(ctx, "", "val")
	assert.Nil(t, RequestContextFromContext(ctx), "empty key should not create RequestContext")
}

func TestRangeSlice(t *testing.T) {
	o := &Options{}
	o.Add("a", 1)
	o.Add("b", 2)
	o.Add("c", 3)

	collected := make(map[string]any)
	o.RangeSlice(func(key, value any) bool {
		collected[key.(string)] = value
		return true
	})
	assert.Len(t, collected, 3)
	assert.Equal(t, 1, collected["a"])
	assert.Equal(t, 2, collected["b"])
	assert.Equal(t, 3, collected["c"])
}

func TestRangeSlice_EarlyStop(t *testing.T) {
	o := &Options{}
	o.Add("a", 1)
	o.Add("b", 2)
	o.Add("c", 3)

	count := 0
	o.RangeSlice(func(key, value any) bool {
		count++
		return false // stop after first
	})
	assert.Equal(t, 1, count)
}

func TestRangeSlice_Empty(t *testing.T) {
	o := &Options{}
	called := false
	o.RangeSlice(func(key, value any) bool {
		called = true
		return true
	})
	assert.False(t, called)
}
