package dogma

import (
	"context"
)

// An IntegrationMessageHandler is an application-defined message handler that
// performs actions outside the Dogma application, such as sending emails or
// processing payments using a third-party API.
//
// It handles [Command] messages, produces external side-effects, and optionally
// records [Event] messages.
type IntegrationMessageHandler interface {
	// Configure declares the handler's configuration by calling methods on c.
	//
	// The configuration includes the handler's identity and message routes.
	//
	// The engine calls this method at least once during startup. It must
	// produce the same configuration each time it's called.
	Configure(c IntegrationConfigurer)

	// HandleCommand handles a [Command] message by performing an action outside
	// the Dogma application.
	//
	// It may cause side-effects in external systems, such as invoking a
	// third-party API. It may use s to record one or more [Event] messages that
	// describe the outcome.
	//
	// The engine atomically persists the events recorded by exactly one
	// successful invocation of this method for each command message. It doesn't
	// guarantee the order, number, or concurrency of those attempts. The
	// implementation must ensure that the command's external side-effects are
	// idempotent and safe for concurrent execution.
	HandleCommand(
		ctx context.Context,
		s IntegrationCommandScope,
		c Command,
	) error
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

// IntegrationCommandScope represents the context within which the engine
// invokes [IntegrationMessageHandler].HandleCommand.
type IntegrationCommandScope interface {
	HandlerScope

	// RecordEvent records an [Event] that results from handling the [Command].
	//
	// The engine doesn't persist the event until
	// [IntegrationMessageHandler].HandleCommand returns successfully. It
	// persists all events recorded using the same scope atomically.
	RecordEvent(Event)
}

// IntegrationRoute is an interface for types that represent a relationship
// between an [IntegrationMessageHandler] and a message type.
//
// Use [HandlesCommand] or [RecordsEvent] to create an IntegrationRoute.
type IntegrationRoute interface {
	MessageRoute
	isIntegrationRoute()
}
