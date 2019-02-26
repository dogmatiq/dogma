package dogma

import "context"

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
	// Name sets the name of the handler.
	//
	// It MUST be called exactly once within a single call to Configure().
	//
	// The name MUST be a non-empty UTF-8 string consisting solely of printable
	// Unicode characters. A printable character is any character from the
	// Letter, Mark, Number, Punctuation or Symbol categories.
	//
	// Each handler within an application MUST have a unique name. Although not
	// recommended, a handler MAY share its name with the application itself.
	Name(n string)

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
