package dogma_test

import (
	"time"

	. "github.com/dogmatiq/dogma"
)

// scopeWithNow is a compile-time test to ensure all scope interfaces have a Now() method
type scopeWithNow interface {
	Now() time.Time
}

// Compile-time interface checks to ensure all scope interfaces implement Now()
var (
	_ scopeWithNow = (AggregateCommandScope)(nil)
	_ scopeWithNow = (ProcessEventScope)(nil)
	_ scopeWithNow = (ProcessTimeoutScope)(nil)
	_ scopeWithNow = (IntegrationCommandScope)(nil)
	_ scopeWithNow = (ProjectionEventScope)(nil)
	_ scopeWithNow = (ProjectionCompactScope)(nil)
)
