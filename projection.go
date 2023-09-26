package dogma

import (
	"context"
	"time"
)

// A ProjectionMessageHandler builds a projection from events.
//
// The term "read-model" is often used interchangeably with "projection".
//
// Projections use an optimistic concurrency control (OCC) protocol to ensure
// that the engine applies each event to the projection exactly once.
//
// The OCC protocol uses a key/value store that associates engine-defined
// "resources" with their "version". These are they keys and values,
// respectively.
//
// The OCC store can be challenging to implement. The [projectionkit] module
// provides adaptors that implement the OCC protocol using various popular
// database systems.
//
// [projectionkit]: github.com/dogma/projectionkit
type ProjectionMessageHandler interface {
	// Configure describes the handler's configuration to the engine.
	Configure(ProjectionConfigurer)

	// HandleEvent updates the projection to reflect the occurrence of an event.
	//
	// r, c and n are the inputs to the OCC store.
	//
	//   - r is a key that identifies some engine-defined resource
	//   - c is engine's perception of the current version of r
	//   - n is the next version of r, made by handling this event
	//
	// If c is the current version of r in the OCC store, the method MUST
	// attempt to atomically update the projection and the version of r to be n.
	// On success, ok is true and err is nil.
	//
	// If c is not the current version of r an OCC conflict has occurred. The
	// method MUST return with ok set to false and without updating the
	// projection.
	//
	// r, c and n are engine-defined; the application SHOULD NOT infer any
	// meaning from their content. The "current" version of a new resource is
	// the empty byte-slice. nil and empty slices are interchangeable.
	//
	// The engine MAY provide specific guarantees about the order in which it
	// supplies events to the handler. To maximize portability across engines,
	// the handler SHOULD NOT assume any specific ordering. The engine MAY call
	// this method concurrently from separate goroutines or operating system
	// processes.
	//
	// The implementation SHOULD NOT impose a context deadline. Implement the
	// TimeoutHint() method instead.
	HandleEvent(
		ctx context.Context,
		r, c, n []byte,
		s ProjectionEventScope,
		e Event,
	) (ok bool, err error)

	// ResourceVersion returns the current version of a resource.
	//
	// It returns an empty slice if r is not in the OCC store.
	ResourceVersion(ctx context.Context, r []byte) ([]byte, error)

	// CloseResource informs the handler that the engine has no further use for
	// a resource.
	//
	// If r is present in the OCC store the handler SHOULD remove it.
	CloseResource(ctx context.Context, r []byte) error

	// TimeoutHint returns a suitable duration for handling the given event.
	//
	// The duration SHOULD be as short as possible. If no hint is available it
	// MUST be zero.
	//
	// In this context, "timeout" refers to a deadline, not a timeout message.
	TimeoutHint(Message) time.Duration

	// Compact attempts to reduce the size of the projection.
	//
	// For example, it may delete unused data, or merge overly granular data.
	//
	// The handler SHOULD compact the projection incrementally such that it
	// makes some progress even if the context's deadline expires.
	Compact(context.Context, ProjectionCompactScope) error
}

// A ProjectionConfigurer configures the engine for use with a specific
// projection message handler.
type ProjectionConfigurer interface {
	// Identity configures the handler's identity.
	//
	// n is a short human-readable name. It MUST be unique within the
	// application at any given time, but MAY change over the handler's
	// lifetime. It MUST contain solely printable, non-space UTF-8 characters.
	// It must be between 1 and 255 bytes (not characters) in length.
	//
	// k is a unique key used to associate engine state with the handler. The
	// key SHOULD NOT change over the handler's lifetime. k MUST be an RFC 4122
	// UUID, such as "5195fe85-eb3f-4121-84b0-be72cbc5722f".
	//
	// Use of hard-coded literals for both values is RECOMMENDED.
	Identity(n string, k string)

	// Routes configures the engine to route certain message types to and from
	// the handler.
	//
	// Projection handlers support the HandlesEvent() route type.
	Routes(...ProjectionRoute)

	// DeliveryPolicy configures how the engine delivers events to the handler.
	//
	// The default policy is UnicastProjectionDeliveryPolicy.
	DeliveryPolicy(ProjectionDeliveryPolicy)
}

// ProjectionEventScope performs engine operations within the context of a call
// to the HandleEvent() method of a [ProjectionMessageHandler].
type ProjectionEventScope interface {
	// RecordedAt returns the time at which the event occurred.
	RecordedAt() time.Time

	// IsPrimaryDelivery returns true on one of the application instances that
	// receive the event, and false on all other instances.
	//
	// This method is useful when the projection must perform some specific
	// operation once per event, such as updating a shared resource that's used
	// by all applications, while still delivering the event to all instances of
	// the application.
	IsPrimaryDelivery() bool

	// Log records an informational message.
	Log(format string, args ...any)
}

// ProjectionCompactScope performs engine operations within the context of a
// call to the Compact() method of a [ProjectionMessageHandler].
type ProjectionCompactScope interface {
	// Now returns the current engine time.
	//
	// The handler SHOULD use the returned time to implement compaction logic
	// that has some time-based component, such as removing data older than a
	// certain age.
	//
	// Under normal operating conditions the engine SHOULD return the current
	// local time. The engine MAY return a different time under some
	// circumstances, such as when executing tests.
	Now() time.Time

	// Log records an informational message.
	Log(format string, args ...any)
}

// NoCompactBehavior is an embeddable type for [ProjectionMessageHandler]
// implementations that do not require compaction.
type NoCompactBehavior struct{}

// Compact does nothing.
func (NoCompactBehavior) Compact(context.Context, ProjectionCompactScope) error {
	return nil
}

type (
	// A ProjectionDeliveryPolicy describes how to deliver events to a
	// projection message handler on engines that support concurrent or
	// distributed execution of a single Dogma application.
	ProjectionDeliveryPolicy interface{ isProjectionDeliveryPolicy() }

	// UnicastProjectionDeliveryPolicy is the default
	// [ProjectionDeliveryPolicy]. It delivers each event to a single instance
	// of the application.
	UnicastProjectionDeliveryPolicy struct{}

	// BroadcastProjectionDeliveryPolicy is a [ProjectionDeliveryPolicy] that
	// delivers each event to a all instance of the application.
	BroadcastProjectionDeliveryPolicy struct {
		// PrimaryFirst defers "secondary delivery" of events until after the
		// "primary delivery" has completed.
		PrimaryFirst bool
	}
)

func (UnicastProjectionDeliveryPolicy) isProjectionDeliveryPolicy()   {}
func (BroadcastProjectionDeliveryPolicy) isProjectionDeliveryPolicy() {}

// ProjectionRoute describes a message type that's routed to a
// [ProjectionMessageHandler].
type ProjectionRoute interface {
	Route
	isProjectionRoute()
}

func (HandlesEventRoute) isProjectionRoute() {}
