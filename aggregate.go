package dogma

import "time"

// A AggregateMessageHandler models business logic and state.
//
// Aggregates are the primary building blocks of an application's domain logic.
// They enforce the domain's strict invariants.
//
// Aggregates use [Command] messages to represent requests to perform some
// specific business logic and change the state of the application
// accordingly. [Event] messages represent those changes.
//
// Aggregates are stateful. An application typically uses multiple instances of
// an aggregate, each with its own state. For example, a banking application may
// use one instance of the "account" aggregate for each bank account.
//
// The state of each instance is application-defined. Often it's a tree of
// related entities and values. The [AggregateRoot] interface represents the
// "root" entity through which the handler accesses the instance's state.
type AggregateMessageHandler interface {
	// Configure describes the handler's configuration to the engine.
	Configure(AggregateConfigurer)

	// New returns an aggregate root instance in its initial state.
	//
	// The return value MUST NOT be nil. It MAY be the zero-value of the root's
	// underlying type.
	//
	// Each call SHOULD return the same type and initial state.
	New() AggregateRoot

	// RouteCommandToInstance returns the ID of the instance that handles a
	// specific command.
	//
	// The return value MUST not be empty. RFC 4122 UUIDs are the RECOMMENDED
	// format for instance IDs.
	RouteCommandToInstance(Command) string

	// HandleCommand executes business logic in response to a command.
	//
	// The handler inspects the root to determine which events to record, if
	// any.
	//
	// The handler SHOULD NOT have any side-effects beyond recording events.
	// Specifically, the implementation MUST NOT modify the root directly. Use
	// [AggregateCommandScope.RecordEvent] to record an event that represents
	// the state change. See also [AggregateRoot.ApplyEvent].
	//
	// If this is the first command routed to this instance, the root is the
	// return value of New(). Otherwise, it's the value of the root as it
	// existed after handling the command.
	//
	// While the engine MAY call this method concurrently from separate
	// goroutines or operating system processes, the state changes and events
	// that represent them always appear to have occurred sequentially.
	HandleCommand(AggregateRoot, AggregateCommandScope, Command)
}

// AggregateRoot is an interface for the domain-specific state of a specific
// aggregate instance.
type AggregateRoot interface {
	// ApplyEvent updates aggregate instance to reflect the occurrence of an
	// event.
	//
	// This implementation of this method is the only code permitted to
	// modify the instance's state.
	//
	// The method SHOULD accept historical events that are no longer routed to
	// this aggregate type. This is typically required by event sourcing engines
	// that sometimes load aggregates into memory by applying their entire
	// history.
	ApplyEvent(Event)
}

// An AggregateConfigurer configures the engine for use with a specific
// aggregate message handler.
type AggregateConfigurer interface {
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
	// Aggregate handlers support the HandlesCommand() and RecordsEvent() route
	// types.
	Routes(...AggregateRoute)

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

// AggregateCommandScope performs engine operations within the context of a call
// to the HandleCommand() method of an [AggregateMessageHandler].
type AggregateCommandScope interface {
	// InstanceID returns the ID of the aggregate instance.
	InstanceID() string

	// RecordEvent records the occurrence of an event.
	//
	// It applies the event to the root such that the applied changes are
	// visible to the handler after this method returns.
	//
	// Recording an event cancels any prior call to Destroy() on this scope.
	RecordEvent(Event)

	// Destroy signals destruction of the aggregate instance.
	//
	// Destroying an aggregate instance discards its state. The first command to
	// target a destroyed instance operates on a new root.
	//
	// Destruction occurs once the HandleCommand() method returns. Any future
	// call to RecordEvent() on this scope prevents destruction.
	//
	// The precise destruction semantics are engine defined. For example,
	// event sourcing engines typically do not destroy the record of the
	// aggregate's historical events.
	Destroy()

	// Now returns the current local time, according to the engine.
	//
	// Use of this method is discouraged. It is preferrable to use information
	// contained within the message or the aggregate root, which provides
	// consistent behavior when message delivery is delayed or retried.
	//
	// If access to the system clock is absolutely necessary, handlers should
	// call this method instead of [time.Now]. It may return a time different to
	// that returned by [time.Now] under some circumstances, such as when
	// executing tests or when accounting for clock skew in a distributed
	// system.
	Now() time.Time

	// Log records an informational message.
	Log(format string, args ...any)
}

// AggregateRoute describes a message type that's routed to or from a
// [AggregateMessageHandler].
type AggregateRoute interface {
	Route
	isAggregateRoute()
}
