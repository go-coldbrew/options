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

	// AddToOptions with empty key should not allocate Options or store an entry
	ctx = AddToOptions(ctx, "", "should-not-store")
	opts := FromContext(ctx)
	assert.Nil(t, opts, "empty key should not create Options in context")

	// Direct Add on a standalone Options should also ignore empty keys
	standalone := &Options{}
	standalone.Add("", "also-should-not-store")
	_, found := standalone.Get("")
	assert.False(t, found, "empty key should not be retrievable via Add")
}
