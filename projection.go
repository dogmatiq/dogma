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
// Each event belongs to an ordered event stream, identified by a UUID. The
// engine delivers events in the order they appear on the stream, using a
// numeric offset to represent each event's position.
//
// A single projection may receive events from multiple streams, which may
// belong to different applications.
//
// To ensure exactly-once event processing, the handler must implement
// optimistic concurrency control (OCC) based on each event's position within an
// event stream. The [github.com/dogmatiq/projectionkit] module provides
// adapters for popular databases, like PostgreSQL and DynamoDB, that handle the
// OCC details.
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
	// The handler must enforce exactly-once semantics by performing an
	// optimistic concurrency check based on the event's [EventStreamPosition],
	// available via [ProjectionEventScope].Position.
	//
	// To do this, the handler must persist an offset for each stream from which
	// it receives events. When handling an event, it compares the event’s
	// offset to the offset of the next unhandled event for that stream. If they
	// match, the handler applies the event and increments the stored offset in
	// a single atomic operation. If they don’t match, an OCC conflict has
	// occurred, and the handler must not apply the event.
	//
	// In either case, the return value n is the offset of the next unhandled
	// event in the stream. If n is the offset immediately after that of the
	// incoming event, the engine considers the event handled successfully.
	// Otherwise, an OCC conflict has occurred, and the engine resumes
	// delivering events from the stream starting at offset n.
	//
	// A non-nil error indicates that the handler encountered a runtime problem
	// other than an OCC conflict.
	//
	// The engine arranges events on streams such that it delivers all [Event]
	// messages recorded within a single scope in the order they occurred. It
	// also preserves the order of events from a single aggregate instance, even
	// across scopes. It doesn't guarantee the relative delivery order of events
	// from different handlers or aggregate instances.
	HandleEvent(
		ctx context.Context,
		s ProjectionEventScope,
		e Event,
	) (n uint64, err error)

	// StreamOffset returns the offset of the next unhandled event in a specific
	// event stream.
	//
	// s is the RFC 4122 UUID that identifies the event stream.
	//
	// The first event in any stream is at offset 0. Accordingly, if the handler
	// hasn’t handled any events from s, this method returns 0.
	StreamOffset(ctx context.Context, s string) (uint64, error)

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

	// Position returns the [EventStreamPosition] of the [Event].
	Position() EventStreamPosition

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
// Embed this type in a [ProjectionMessageHandler] when projection data doesn't
// grow unbounded or when an external system handles compaction.
type NoCompactBehavior struct{}

// Compact returns nil without performing any operations.
func (NoCompactBehavior) Compact(context.Context, ProjectionCompactScope) error {
	return nil
}
