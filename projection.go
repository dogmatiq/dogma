package dogma

import (
	"context"
	"time"
)

// A ProjectionMessageHandler builds a denormalized view of the application's
// state that's optimized for querying.
//
// It handles [Event] messages to build and update projections - partial
// representations of application state tailored to specific use cases. For
// example, a projection might contain a list of customers with active shopping
// carts.
//
// Projection message handlers use an optimistic concurrency control (OCC)
// protocol to ensure exactly-once event processing. The protocol coordinates
// between the engine and the handler using a key/value store that tracks the
// version of engine-defined "resources". This prevents duplicate updates when
// the engine retries operations or when multiple engine instances process the
// same events.
//
// The OCC protocol can be challenging to implement correctly. The
// [github.com/dogmatiq/projectionkit] module provides adapters for popular
// databases like PostgreSQL and DynamoDB that handle the OCC details.
type ProjectionMessageHandler interface {
	// Configure declares the handler's configuration by calling methods on c.
	//
	// The configuration includes the handler's identity and message routes.
	//
	// The engine calls this method at least once during startup. It must
	// produce the same configuration each time it's called.
	Configure(c ProjectionConfigurer)

	// HandleEvent updates the projection to reflect the occurrence of an
	// [Event].
	//
	// The handler must implement the OCC protocol to ensure exactly-once
	// processing. The engine provides three values that coordinate this
	// protocol:
	//
	//   - r identifies an engine-defined resource
	//   - c is the engine's view of r's current version
	//   - n is the next version of r after handling this event
	//
	// The handler must atomically:
	//   1. Verify that c matches r's actual version in the OCC store
	//   2. Update the projection to reflect the event
	//   3. Update r's version to n in the OCC store
	//
	// If all operations succeed, the handler returns ok as true. If c doesn't
	// match r's actual version (an OCC conflict), the handler returns ok as
	// false without modifying the projection. The engine retries conflicts
	// automatically.
	//
	// The handler must treat r, c, and n as opaque values. The engine defines
	// their meaning and content. The "current" version of a new resource is an
	// empty byte slice. The handler must treat nil and empty slices as
	// equivalent.
	//
	// The engine doesn't guarantee any specific event delivery order except
	// that it preserves the order of events from a single aggregate instance.
	// The engine may call this method concurrently from multiple goroutines or
	// processes.
	HandleEvent(
		ctx context.Context,
		r, c, n []byte,
		s ProjectionEventScope,
		e Event,
	) (ok bool, err error)

	// ResourceVersion returns the current version of the resource r from the
	// projection's OCC store.
	//
	// It returns an empty slice if r doesn't exist in the store. The handler
	// must treat nil and empty slices as equivalent.
	//
	// The engine uses this method to determine r's current version before
	// calling [ProjectionMessageHandler].HandleEvent.
	ResourceVersion(ctx context.Context, r []byte) ([]byte, error)

	// CloseResource removes the resource r from the OCC store.
	//
	// The engine calls this method when it no longer needs to track r. The
	// handler should remove r and its version from the store if present. It
	// should succeed even if r doesn't exist.
	CloseResource(ctx context.Context, r []byte) error

	// Compact reduces the projection's size by removing or consolidating data.
	//
	// The handler might delete obsolete entries, merge fine-grained data into
	// summaries, or move old data to archival storage. The specific strategy
	// depends on the projection's purpose and access patterns.
	//
	// The implementation should perform compaction incrementally to make some
	// progress even if ctx reaches its deadline.
	//
	// The engine may call this method at any time, including in parallel with
	// handling an event.
	Compact(ctx context.Context, s ProjectionCompactScope) error
}

// ProjectionConfigurer is the interface that a [ProjectionMessageHandler] uses
// to declare its configuration.
//
// The engine provides the implementation to
// [ProjectionMessageHandler].Configure during startup.
type ProjectionConfigurer interface {
	HandlerConfigurer

	// Routes declares the message types that the handler consumes.
	//
	// It accepts routes created by [HandlesEvent].
	Routes(...ProjectionRoute)
}

// ProjectionEventScope represents the context within which a
// [ProjectionMessageHandler] handles an [Event] message.
type ProjectionEventScope interface {
	HandlerScope

	// RecordedAt returns the time at which the [Event] occurred.
	RecordedAt() time.Time
}

// ProjectionCompactScope represents the context within which a
// [ProjectionMessageHandler] compacts its data.
type ProjectionCompactScope interface {
	HandlerScope
}

// ProjectionRoute describes a message type that's routed to a
// [ProjectionMessageHandler].
type ProjectionRoute interface {
	MessageRoute
	isProjectionRoute()
}

// NoCompactBehavior is an embeddable type for [ProjectionMessageHandler]
// implementations that don't require compaction.
//
// Embed this type in a [ProjectionMessageHandler] to when projection data
// doesn't grow unbounded or when an external system handles compaction.
type NoCompactBehavior struct{}

// Compact returns nil without performing any operations.
func (NoCompactBehavior) Compact(context.Context, ProjectionCompactScope) error {
	return nil
}
