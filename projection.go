package dogma

import (
	"context"
	"time"
)

// ProjectionMessageHandler is an interface implemented by the application and
// used by the engine to build a "projection" (also known as a "read model", or
// "query model") from events that occur within the application.
//
// Projection message handlers consume event messages, and do not produce
// messages of any kind.
type ProjectionMessageHandler interface {
	// Configure produces a configuration for this handler by calling methods on
	// the configurer c.
	//
	// The implementation MUST allow for multiple calls to Configure(). Each
	// call SHOULD produce the same configuration.
	//
	// The engine MUST call Configure() before calling HandleEvent(). It is
	// RECOMMENDED that the engine only call Configure() once per handler.
	Configure(c ProjectionConfigurer)

	// HandleEvent updates the projection to reflect the occurrence of an event.
	//
	// If nil is returned, the projection has been updated successfully.
	//
	// If an error is returned, the projection SHOULD be left in the state it
	// was before HandleEvent() was called.
	//
	// The engine SHOULD provide "at-least-once" delivery guarantees to the
	// handler. That is, the engine should call HandleEvent() with the same
	// event message until a nil error is returned.
	//
	// The supplied context parameter SHOULD have a deadline. The implementation
	// SHOULD NOT impose its own deadline. Instead a suitable timeout duration
	// can be suggested to the engine via the handler's TimeoutHint() method.
	//
	// The engine MUST NOT call HandleEvent() with any message of a type that
	// has not been configured for consumption by a prior call to Configure().
	// If any such message is passed, the implementation MUST panic with the
	// UnexpectedMessage value.
	//
	// The engine MAY provide guarantees about the order in which event messages
	// will be passed to HandleEvent(), however in the interest of engine
	// portability the implementation SHOULD NOT assume that HandleEvent() will
	// be called with events in the same order that they were recorded.
	//
	// The engine MAY call HandleEvent() from multiple goroutines concurrently.
	HandleEvent(ctx context.Context, s ProjectionEventScope, m Message) error

	// TimeoutHint returns a duration that is suitable for computing a deadline
	// for the handling of the given message by this handler.
	//
	// The hint SHOULD be as short as possible. The implementation MAY return a
	// zero-value to indicate that no hint can be made.
	//
	// The engine SHOULD use a duration as close as possible to the hint. Use of
	// a duration shorter than the hint is NOT RECOMMENDED, as this will likely
	// lead to repeated message handling failures.
	TimeoutHint(m Message) time.Duration
}

// IdempotentProjectionMessageHandler is specialization of
// ProjectionMessageHandler with additional features that allow the engine to
// provide exactly-once delivery guarantees.
type IdempotentProjectionMessageHandler interface {
	ProjectionMessageHandler
}

// ProjectionConfigurer is an interface implemented by the engine and used by
// the application to configure options related to a ProjectionMessageHandler.
//
// It is passed to ProjectionMessageHandler.Configure(), typically upon
// initialization of the engine.
//
// In the context of this interface, "the handler" refers to the handler on
// which Configure() has been called.
type ProjectionConfigurer interface {
	// Identity sets unique identifiers for the handler.
	//
	// It MUST be called exactly once within a single call to Configure().
	//
	// The name is a human-readable identifier for the handler. Each handler
	// within an application MUST have a unique name. Handler names SHOULD be
	// distinct from the application's name. The name MAY be changed over time
	// to best reflect the purpose of the handler.
	//
	// The key is an immutable identifier for the handler. Its purpose is to
	// allow engine implementations to associate ancillary data with the
	// handler, such as application state or message routing information.
	//
	// The application and the handlers within it MUST have distinct keys. The
	// key MUST NOT be changed. The RECOMMENDED key format is RFC 4122 UUID,
	// generated when the handler is first implemented.
	//
	// Both the name and the key MUST be non-empty UTF-8 strings consisting
	// solely of printable Unicode characters, excluding whitespace. A printable
	// character is any character from the Letter, Mark, Number, Punctuation or
	// Symbol categories.
	Identity(name string, key string)

	// ConsumesEventType configures the engine to route event messages of the
	// same type as m to the handler.
	//
	// It MUST be called at least once within a call to Configure(). It MUST NOT
	// be called more than once with an event message of the same type.
	//
	// Multiple handlers within an application MAY consume event messages of the
	// same type.
	//
	// The "content" of m MUST NOT be used, inspected, or treated as meaningful
	// in any way, only its runtime type information may be used.
	ConsumesEventType(m Message)
}

// ProjectionEventScope is an interface implemented by the engine and used by
// the application to perform operations within the context of handling a
// specific event message.
type ProjectionEventScope interface {
	// MessageIdentity returns a key-value pair that uniquely identify the event
	// being handled.
	//
	// In combination, the returned values MUST be unique to the specific event
	// message within this projection handler. There is no guarantee that the
	// returned values will be globally unique.
	//
	// The returned values are engine-defined, and MUST be treated as opaque
	// data structures by the handler.
	MessageIdentity() (k, v []byte)

	// RecordedAt returns the time at which the event was recorded.
	RecordedAt() time.Time

	// Log records an informational message within the context of the message
	// that is being handled.
	Log(f string, v ...interface{})
}
