package dogma

// An AggregateMessageHandler models business logic and state within a Dogma
// application by handling [Command] messages and recording [Event] messages.
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

// An AggregateRoot is an interface for an application's in-memory
// representation of an aggregate instance.
type AggregateRoot interface {
	// ApplyEvent updates the aggregate instance to reflect the occurrence of an
	// event.
	//
	// The engine calls this method when loading the instance from historical
	// events or recording a new event. It must handle all historical event
	// types, including those no longer routed to this aggregate.
	ApplyEvent(Event)
}

// AggregateConfigurer is the interface an [AggregateMessageHandler] uses to
// declare its configuration.
//
// The engine provides the implementation to
// [AggregateMessageHandler].Configure during startup.
type AggregateConfigurer interface {
	HandlerConfigurer

	// Routes associates message types with the handler, indicating which types
	// it consumes and produces.
	//
	// It accepts routes created by [HandlesCommand] and [RecordsEvent].
	Routes(...AggregateRoute)
}

// AggregateCommandScope represents the context within which the engine
// invokes [AggregateMessageHandler].HandleCommand.
type AggregateCommandScope interface {
	HandlerScope

	// InstanceID returns the ID of the aggregate instance that the command
	// targets, as returned by [AggregateMessageHandler].RouteCommandToInstance.
	InstanceID() string

	// RecordEvent records an [Event] that results from handling the [Command].
	//
	// The engine applies the event to the aggregate root by calling
	// [AggregateRoot].ApplyEvent, making the state changes visible to the
	// handler immediately.
	//
	// The engine doesn't persist the event until
	// [AggregateMessageHandler].HandleCommand completes successfully. It
	// persists all events recorded using the same scope atomically.
	RecordEvent(Event)
}

// AggregateRoute describes a message type that's routed to or from a
// [AggregateMessageHandler].
type AggregateRoute interface {
	MessageRoute
	isAggregateRoute()
}
