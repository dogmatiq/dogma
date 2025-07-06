package dogma

// An AggregateMessageHandler is an application-defined message handler that
// handles [Command] messages and records [Event] messages that represent
// changes to application state.
//
// An aggregate is a collection of related business entities that behave as a
// cohesive whole, such as a shopping cart and the items within it. The
// aggregate message handler manages the behavior and state of such aggregates.
//
// Each aggregate message handler typically manages multiple instances, where
// each instance represents a separate occurrence of the aggregate. For example,
// a shopping cart aggregate message handler may manage one instance per
// customer.
type AggregateMessageHandler interface {
	// Configure declares the handler's configuration by calling methods on c.
	//
	// The configuration includes the handler's identity and message routes.
	//
	// The engine calls this method at least once during startup. It must
	// produce the same configuration each time it's called.
	Configure(AggregateConfigurer)

	// New returns a new [AggregateRoot] for an aggregate instance.
	//
	// The engine calls this method to get a "blank slate" when handling the
	// first command for a new instance or when reconstructing an existing
	// instance from its historical events.
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

// An AggregateRoot is an interface for an application's working representation
// of an aggregate instance used within [AggregateMessageHandler]
// implementations.
//
// The aggregate root encapsulates business logic and provides a way to inspect
// the current state when making decisions about which events to record. The
// recorded events are the authoritative source of truth, not the AggregateRoot.
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
