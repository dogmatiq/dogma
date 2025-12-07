package dogma

import (
	"context"
)

// An IntegrationMessageHandler connects a Dogma application to an external
// system.
//
// It handles [Command] messages, and may record [Event] messages to describe
// the outcome.
//
// Integration message handlers encapsulate integration logic such as sending an
// email, charging a credit card, or updating a third-party API. They should
// implement only the "glue code". For example, an integration with a payment
// gateway might debit a customer's credit card, but shouldn't apply discounts
// or calculate sales tax.
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

	// Routes declares the message types that the handler consumes and produces.
	//
	// It accepts routes created by [HandlesCommand] and [RecordsEvent].
	Routes(...IntegrationRoute)

	// ConcurrencyPreference provides a hint to the engine as to the best way to
	// handle concurrent messages for this handler.
	//
	// The default is [MaximizeConcurrency].
	ConcurrencyPreference(ConcurrencyPreference)
}

// IntegrationCommandScope represents the context within which an
// [IntegrationMessageHandler] handles a [Command] message.
type IntegrationCommandScope interface {
	HandlerScope

	// RecordEvent records an [Event] that results from handling the [Command].
	//
	// The engine persists all events recorded within this scope in a single
	// atomic operation after the [IntegrationMessageHandler] finishes handling
	// the inbound command. If the handler returns a non-nil error, the engine
	// discards the events.
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
