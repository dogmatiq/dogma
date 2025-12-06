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
	// The [ProjectionEventScope] s exposes the ID of the event stream to which
	// the incoming event belongs, the event's offset within that stream, and
	// the checkpoint offset at which the engine expects the handler to resume
	// consuming from the stream.
	//
	// If the engine and handler agree on the checkpoint offset, the handler
	// must atomically apply the event and update its checkpoint offset to one
	// greater than the incoming event's offset. Otherwise, the handler must not
	// modify any data.
	//
	// In either case, the method returns cp, the new checkpoint offset for this
	// stream. If cp is one greater than the offset of the incoming event, the
	// engine considers the event handled successfully. Otherwise, an OCC
	// conflict has occurred, and the engine resumes delivering events starting
	// at cp.
	//
	// A non-nil error indicates that the handler encountered a runtime problem
	// other than an OCC conflict.
	//
	// The engine arranges events on streams such that it delivers all events
	// within a single scope in the order they occurred. It also preserves the
	// order of events from a single aggregate instance, even across scopes. It
	// doesn't guarantee the relative delivery order of events from different
	// handlers or aggregate instances.
	//
	// See:
	//  - [ProjectionEventScope].StreamID
	//  - [ProjectionEventScope].Offset
	//  - [ProjectionEventScope].CheckpointOffset
	HandleEvent(
		ctx context.Context,
		s ProjectionEventScope,
		e Event,
	) (cp uint64, err error)

	// CheckpointOffset returns the offset at which the handler expects to
	// resume handling events from a specific stream.
	//
	// id is an RFC 9562 UUID that identifies the event stream, such as
	// "c50b5804-8312-4c61-b32c-9fbf49688db3". UUIDs are case-insensitive, but
	// the engine must use a lowercase representation.
	CheckpointOffset(ctx context.Context, id string) (uint64, error)

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
	//
	// Not all projections need compaction. Embed [NoCompactBehavior] in the
	// handler to indicate compaction not required.
	Compact(ctx context.Context, s ProjectionCompactScope) error

	// Reset clears all projection data and checkpoint offsets such that the
	// projection data is rebuilt by handling all historical events.
	//
	// Not all projections can be reset. Embed [NoResetBehavior] in the handler
	// to indicate that reset is not supported.
	Reset(ctx context.Context, s ProjectionResetScope) error
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

	// ConcurrencyPreference provides a hint to the engine as to the best way to
	// handle concurrent messages for this handler.
	//
	// The default is [MaximizeConcurrency].
	ConcurrencyPreference(ProjectionConcurrencyPreference)
}

// ProjectionEventScope represents the context within which a
// [ProjectionMessageHandler] handles an [Event] message.
type ProjectionEventScope interface {
	HandlerScope

	// RecordedAt returns the time at which the [Event] occurred.
	RecordedAt() time.Time

	// StreamID returns the RFC 9562 UUID that identifies the event stream to
	// which the [Event] belongs.
	StreamID() string

	// Offset returns the event's zero-based offset within the stream.
	Offset() uint64

	// CheckpointOffset returns the offset from which the handler should resume
	// handling events from this stream, according to the engine.
	//
	// It may be lower than the incoming event's offset when the stream contains
	// event types that the handler doesn't consume.
	CheckpointOffset() uint64
}

// ProjectionCompactScope represents the context within which a
// [ProjectionMessageHandler] compacts its data.
type ProjectionCompactScope interface {
	HandlerScope
}

// ProjectionResetScope represents the context within which a
// [ProjectionMessageHandler] resets its data.
type ProjectionResetScope interface {
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

// NoResetBehavior is an embeddable type for [ProjectionMessageHandler]
// implementations that don't support resetting their state.
//
// Embed this type in a [ProjectionMessageHandler] when resetting projection
// data isn't feasible or required.
type NoResetBehavior struct{}

// Reset returns an error indicating that reset is not supported.
func (NoResetBehavior) Reset(context.Context, ProjectionResetScope) error {
	return ErrNotSupported
}

// ProjectionConcurrencyPreference is a hint to the engine as to the best way to
// handle concurrent messages for a [ProjectionMessageHandler].
type ProjectionConcurrencyPreference = concurrencyPreference
