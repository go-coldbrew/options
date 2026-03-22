package options

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOptions(t *testing.T) {
	TestKey := "TestKey"
	TestValue := "TestValue"
	ctx := context.Background()
	// add
	ctx = AddToOptions(ctx, TestKey, TestValue)
	options := FromContext(ctx)
	// fetch
	value, found := options.Get(TestKey)
	assert.True(t, found, "key should be found")
	assert.Equal(t, TestValue, value, "values should be equal")

	//delete
	options.Del(TestKey)
	//fetch
	options2 := FromContext(ctx)
	value, found = options2.Get(TestKey)
	assert.False(t, found, "key should NOT be found")
	assert.NotEqual(t, TestValue, value, "values should NOT be equal")
}

func TestEmptyKeyIgnored(t *testing.T) {
	ctx := context.Background()

	// AddToOptions with empty key should not store an entry
	ctx = AddToOptions(ctx, "", "should-not-store")
	options := FromContext(ctx)
	assert.NotNil(t, options, "options should be initialized")
	_, found := options.Get("")
	assert.False(t, found, "empty key should not be retrievable via AddToOptions")

	// Direct Add with empty key should not store an entry
	options.Add("", "also-should-not-store")
	_, found = options.Get("")
	assert.False(t, found, "empty key should not be retrievable via Add")
}
