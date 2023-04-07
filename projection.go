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
	Configure(c ProjectionConfigurer)

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
	// HandleEvent() concurrently from separate goroutines or operating system
	// processes.
	//
	// The implementation SHOULD NOT impose a context deadline. Instead, use the
	// TimeoutHint() method to provide the engine with a suitable timeout
	// duration.
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
	// See [NoTimeoutHintBehavior].
	TimeoutHint(m Message) time.Duration

	// Compact attempts to reduce the size of the projection.
	//
	// For example, it may delete unused data, or merge overly granular data.
	//
	// The handler SHOULD compact the projection incrementally such that it
	// makes some progress even if the context's deadline expires.
	//
	// See [NoCompactBehavior].
	Compact(ctx context.Context, s ProjectionCompactScope) error
}

// A ProjectionConfigurer configures the engine for use with a specific
// projection message handler.
//
// See [ProjectionMessageHandler.Configure]().
type ProjectionConfigurer interface {
	// Identity configures how the engine identifies the handler.
	//
	// The handler MUST call Identity().
	//
	// name is a human-readable identifier for the handler. Each handler within
	// an application MUST have a unique name. The name MAY change over time to
	// best reflect the purpose of the handler.
	//
	// name MUST be a non-empty UTF-8 string consisting solely of printable
	// Unicode characters, excluding whitespace. A printable character is any
	// character from the Letter, Mark, Number, Punctuation or Symbol
	// categories.
	//
	// key is an unique identifier for the handler that's used by the engine to
	// correlate its internal state with this handler. For that reason the key
	// SHOULD NOT change once in use.
	//
	// key MUST be an [RFC 4122] UUID expressed as a hyphen-separated, lowercase
	// hexadecimal string, such as "5195fe85-eb3f-4121-84b0-be72cbc5722f".
	//
	// [RFC 4122]: https://www.rfc-editor.org/rfc/rfc4122
	Identity(name string, key string)

	// ConsumesEventType configures the engine to route events of a specific
	// type to the handler.
	//
	// The handler MUST call ConsumesEventType() at least once.
	//
	// The event SHOULD be the zero-value of its type; the engine uses the type
	// information, but not the value itself.
	ConsumesEventType(e Event)

	// DeliveryPolicy configures how the engine delivers events to the handler.
	//
	// It accepts a list of candidate policies, in order of preference. It
	// returns the first candidate that the engine supports.
	//
	// The default policy is UnicastProjectionDeliveryPolicy.
	DeliveryPolicy(candidates ...ProjectionDeliveryPolicy) ProjectionDeliveryPolicy
}

// A ProjectionDeliveryPolicy describes how to deliver events to a projection
// message handler on engines that support concurrent or distributed execution
// of a single Dogma application.
type ProjectionDeliveryPolicy interface {
	isProjectionDeliveryPolicy()
}

// UnicastProjectionDeliveryPolicy is the default ProjectionDeliveryPolicy. It
// delivers each event to a single instance of the application.
type UnicastProjectionDeliveryPolicy struct {
}

func (UnicastProjectionDeliveryPolicy) isProjectionDeliveryPolicy() {}

// ProjectionEventScope performs operations within the context of a call to
// [ProjectionMessageHandler.HandleEvent]().
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
	Log(f string, v ...interface{})
}

// ProjectionCompactScope performs operations within the context of a call to
// [ProjectionMessageHandler.Compact]().
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
	Log(f string, v ...interface{})
}

// NoCompactBehavior can be embedded in [ProjectionMessageHandler]
// implementations to denote that the projection does not require compaction.
//
// It provides a no-op implementation of [ProjectionMessageHandler.Compact]()
// that always returns a nil error.
type NoCompactBehavior struct{}

// Compact does nothing.
func (NoCompactBehavior) Compact(ctx context.Context, s ProjectionCompactScope) error {
	return nil
}
