package options_test

import (
	"context"
	"fmt"

	"github.com/go-coldbrew/options"
)

func ExampleAddToOptions() {
	ctx := context.Background()

	// Add request-scoped metadata to context
	ctx = options.AddToOptions(ctx, "tenant", "acme-corp")
	ctx = options.AddToOptions(ctx, "region", "us-west-2")

	// Retrieve values downstream
	opts := options.FromContext(ctx)
	if tenant, ok := opts.Get("tenant"); ok {
		fmt.Println("tenant:", tenant)
	}
	if region, ok := opts.Get("region"); ok {
		fmt.Println("region:", region)
	}
	// Output:
	// tenant: acme-corp
	// region: us-west-2
}

func ExampleFromContext() {
	ctx := context.Background()

	// Without any options set, FromContext returns an empty Options
	opts := options.FromContext(ctx)
	_, found := opts.Get("missing-key")
	fmt.Println("found:", found)
	// Output: found: false
}
