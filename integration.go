package dogma

import (
	"context"
)

// An IntegrationMessageHandler connects a Dogma application to external systems
// by handling [Command] messages and optionally recording [Event] messages.
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

// IntegrationConfigurer is the interface an [IntegrationMessageHandler] uses to
// declare its configuration.
//
// The engine provides the implementation to
// [IntegrationMessageHandler].Configure during startup.
type IntegrationConfigurer interface {
	HandlerConfigurer

	// Routes associates message types with the handler, indicating which types
	// it consumes and produces.
	//
	// It accepts routes created by [HandlesCommand] and [RecordsEvent].
	Routes(...IntegrationRoute)
}

// IntegrationCommandScope performs engine operations within the context of a
// call to the HandleCommand() method of an [IntegrationMessageHandler].
type IntegrationCommandScope interface {
	// RecordEvent records the occurrence of an event.
	RecordEvent(Event)

	// Now returns the current local time, according to the engine.
	//
	// Handlers should call this method instead of [time.Now]. It may return a
	// time different to that returned by [time.Now] under some circumstances,
	// such as when executing tests or when accounting for clock skew in a
	// distributed system.
	Now() time.Time

	// Log records an informational message.
	Log(format string, args ...any)
}

// IntegrationRoute describes a message type that's routed to or from a
// [IntegrationMessageHandler].
type IntegrationRoute interface {
	MessageRoute
	isIntegrationRoute()
}
