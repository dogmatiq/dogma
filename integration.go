package dogma

import (
	"context"
	"time"
)

// An IntegrationMessageHandler integrates a Dogma application with external and
// non-message-based systems.
//
// The engine does not keep any state for integration handlers.
type IntegrationMessageHandler interface {
	// Configure describes the handler's configuration to the engine.
	Configure(IntegrationConfigurer)

	// HandleCommand handles a command, typically by invoking some external API.
	//
	// It MAY optionally record events that describe the outcome of the command.
	//
	// The engine MAY call this method concurrently from separate goroutines or
	// operating system processes.
	HandleCommand(context.Context, IntegrationCommandScope, Command) error
}

// A IntegrationConfigurer configures the engine for use with a specific
// integration message handler.
type IntegrationConfigurer interface {
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
	// Integration handlers support the HandlesCommand() and RecordsEvent()
	// route types.
	Routes(...IntegrationRoute)

	// Disable prevents the handler from receiving any messages.
	//
	// The engine MUST NOT call any methods other than Configure() on a disabled
	// handler.
	//
	// Disabling a handler is useful when the handler's configuration prevents
	// it from operating, such as when it's missing a required dependency,
	// without requiring the user to conditionally register the handler with the
	// application.
	Disable(...DisableOption)
}

// IntegrationCommandScope performs engine operations within the context of a
// call to the HandleCommand() method of an [IntegrationMessageHandler].
type IntegrationCommandScope interface {
	// RecordEvent records the occurrence of an event.
	RecordEvent(Event)

	// Now returns the current engine time.
	//
	// The handler SHOULD use the returned time instead of calling time.Now()
	// directly to ensure compatibility with testing frameworks that manipulate
	// time.
	//
	// Under normal operating conditions the engine SHOULD return the current
	// local time. The engine MAY return a different time under some
	// circumstances, such as when executing tests.
	Now() time.Time

	// Log records an informational message.
	Log(format string, args ...any)
}

// IntegrationRoute describes a message type that's routed to or from a
// [IntegrationMessageHandler].
type IntegrationRoute interface {
	Route
	isIntegrationRoute()
}
