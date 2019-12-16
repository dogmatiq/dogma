package dogma

import (
	"context"
	"time"
)

// IntegrationMessageHandler is an interface implemented by the application and
// used by the engine to integrate with non-message-based systems.
//
// Integration message handlers consume command messages, and optionally produce
// event messages.
type IntegrationMessageHandler interface {
	// Configure produces a configuration for this handler by calling methods on
	// the configurer c.
	//
	// The implementation MUST allow for multiple calls to Configure(). Each
	// call SHOULD produce the same configuration.
	//
	// The engine MUST call Configure() before calling HandleCommand(). It is
	// RECOMMENDED that the engine only call Configure() once per handler.
	Configure(c IntegrationConfigurer)

	// HandleCommand handles a command message.
	//
	// If nil is returned, the command has been handled successfully.
	//
	// The engine SHOULD provide "at-least-once" delivery guarantees to the
	// handler. That is, the engine should call HandleCommand() with the same
	// command message until a nil error is returned.
	//
	// The supplied context parameter SHOULD have a deadline. The implementation
	// SHOULD NOT impose its own deadline. Instead a suitable timeout duration
	// can be suggested to the engine via the handler's TimeoutHint() method.
	//
	// The engine MUST NOT call HandleCommand() with any message of a type that
	// has not been configured for consumption by a prior call to Configure().
	// If any such message is passed, the implementation MUST panic with the
	// UnexpectedMessage value.
	//
	// The implementation MUST NOT assume that HandleCommand() will be called
	// with commands in the same order that they were executed.
	//
	// The engine MAY call HandleCommand() from multiple goroutines concurrently.
	HandleCommand(ctx context.Context, s IntegrationCommandScope, m Message) error

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

// IntegrationConfigurer is an interface implemented by the engine and used
// by the application to configure options related to a IntegrationMessageHandler.
//
// It is passed to IntegrationMessageHandler.Configure(), typically upon
// initialization of the engine.
//
// In the context of this interface, "the handler" refers to the handler on
// which Configure() has been called.
type IntegrationConfigurer interface {
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
	// key MUST NOT be changed. The RECOMMENDED key format is an RFC 4122 UUID
	// represented as a hyphen-separated, lowercase hexadecimal string, such as
	// "5195fe85-eb3f-4121-84b0-be72cbc5722f".
	//
	// Both identifiers MUST be non-empty UTF-8 strings consisting solely of
	// printable Unicode characters, excluding whitespace. A printable character
	// is any character from the Letter, Mark, Number, Punctuation or Symbol
	// categories.
	//
	// The engine MUST NOT perform any case-folding or normalization of
	// identifiers. Therefore, two identifiers compare as equivalent if and only
	// if they consist of the same sequence of bytes.
	Identity(name string, key string)

	// ConsumesCommandType configures the engine to route command messages of
	// the same type as m to the handler.
	//
	// It MUST be called at least once within a call to Configure(). It MUST NOT
	// be called more than once with a command message of the same type.
	//
	// A given command type MUST be routed to exactly one handler within an
	// application.
	//
	// The "content" of m MUST NOT be used, inspected, or treated as meaningful
	// in any way, only its runtime type information may be used.
	ConsumesCommandType(m Message)

	// ProducesEventType instructs the engine that the handler records events of
	// the same type as m.
	//
	// It MUST NOT be called more than once with an event message of the same
	// type.
	//
	// A given event type MUST be produced by exactly one handler within an
	// application.
	//
	// The "content" of m MUST NOT be used, inspected, or treated as meaningful
	// in any way, only its runtime type information may be used.
	ProducesEventType(m Message)
}

// IntegrationCommandScope is an interface implemented by the engine and used by
// the application to perform operations within the context of handling a
// specific integration command message.
type IntegrationCommandScope interface {
	// RecordEvent records the occurrence of an event as a result of the command
	// message that is being handled.
	//
	// It MUST NOT be called with a message of any type that has not been
	// configured for production by a prior call to Configure().
	RecordEvent(m Message)

	// Log records an informational message within the context of the message
	// that is being handled.
	Log(f string, v ...interface{})
}
