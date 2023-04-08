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
	//
	// The implementation SHOULD NOT impose a context deadline; implement the
	// TimeoutHint() method instead.
	HandleCommand(context.Context, IntegrationCommandScope, Command) error

	// TimeoutHint returns a suitable duration for handling the given message.
	//
	// The duration SHOULD be as short as possible. If no hint is available it
	// MUST be zero.
	//
	// In this context, "timeout" refers to a deadline, not a timeout message.
	TimeoutHint(Message) time.Duration
}

// A IntegrationConfigurer configures the engine for use with a specific
// integration message handler.
type IntegrationConfigurer interface {
	// Identity configures the handler's identity.
	//
	// n is a short human-readable name. It MUST be unique within the
	// application. The name MAY change over the handler's lifetime. n MUST
	// contain solely printable, non-space UTF-8 characters.
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

	// ConsumesCommandType configures the engine to route commands of a specific
	// type to the handler.
	//
	// The application's configuration MUST route each command type to a single
	// handler.
	//
	// The command SHOULD be the zero-value of its type; the engine uses the
	// type information, but not the value itself.
	//
	// Deprecated: Use IntegrationConfigurer.Routes() instead.
	ConsumesCommandType(Command)

	// ProducesEventType configures the engine to use the handler as the source
	// of events of a specific type.
	//
	// The application's configuration MUST source each event type from a single
	// handler.
	//
	// The event SHOULD be the zero-value of its type; the engine uses the type
	// information, but not the value itself.
	//
	// Deprecated: Use IntegrationConfigurer.Routes() instead.
	ProducesEventType(Event)
}

// IntegrationCommandScope performs engine operations within the context of a
// call to the HandleCommand() method of an [IntegrationMessageHandler].
type IntegrationCommandScope interface {
	// RecordEvent records the occurrence of an event.
	RecordEvent(Event)

	// Log records an informational message.
	Log(format string, args ...any)
}

// IntegrationRoute describes a message type that's routed to or from a
// [IntegrationMessageHandler].
type IntegrationRoute interface{ isIntegrationRoute() }

func (HandlesCommandRoute) isIntegrationRoute() {}
func (RecordsEventRoute) isIntegrationRoute()   {}
