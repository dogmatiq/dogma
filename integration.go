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
	// Configure describes the handler's configuration to the engine.
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
	// The engine MAY call HandleCommand() from multiple goroutines
	// concurrently.
	HandleCommand(ctx context.Context, s IntegrationCommandScope, c Command) error

	// TimeoutHint returns a suitable duration for handling the given message.
	//
	// The duration SHOULD be as short as possible. If no hint is available it
	// MUST be zero.
	//
	// See [NoTimeoutHintBehavior].
	TimeoutHint(m Message) time.Duration
}

// A IntegrationConfigurer configures the engine for use with a specific
// integration message handler.
//
// See [IntegrationMessageHandler.Configure]().
type IntegrationConfigurer interface {
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

	// ConsumesCommandType configures the engine to route commands of a specific
	// type to the handler.
	//
	// The handler MUST call ConsumesCommandType() at least once.
	//
	// The application's configuration MUST route each command type to a single
	// handler.
	//
	// The command SHOULD be the zero-value of its type; the engine uses the
	// type information, but not the value itself.
	ConsumesCommandType(c Command)

	// ProducesEventType configures the engine to use the handler as the source
	// of events of a specific type.
	//
	// The handler MUST call ProducesEventType() at least once.
	//
	// The application's configuration MUST source each event type from a single
	// handler.
	//
	// The event SHOULD be the zero-value of its type; the engine uses the type
	// information, but not the value itself.
	ProducesEventType(e Event)
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
	RecordEvent(e Event)

	// Log records an informational message.
	Log(f string, v ...interface{})
}
