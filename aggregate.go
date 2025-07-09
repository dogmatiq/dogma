package dogma

// An AggregateMessageHandler encapsulates an application's "invariant" logic -
// conditions that must always hold true.
//
// It handles [Command] messages and records [Event] messages that represent
// changes to application state.
//
// An "aggregate" is a collection of related business entities that behave as a
// cohesive whole, such as a shopping cart and the items within it. The
// aggregate message handler manages the behavior and state of such aggregates.
// The term doesn't refer to data aggregation or summarization.
//
// Each aggregate message handler typically manages multiple instances, where
// each instance represents a distinct occurrence of the aggregate. For example,
// a shopping cart aggregate might use one instance per customer.
//
// Aggregates enforce business rules that must always hold true for a specific
// instance. For example, a shopping cart aggregate might prevent checkout if
// the cart is empty, or limit the number of items to 10.
type AggregateMessageHandler interface {
	// Configure declares the handler's configuration by calling methods on c.
	//
	// The configuration includes the handler's identity and message routes.
	//
	// The engine calls this method at least once during startup. It must
	// produce the same configuration each time it's called.
	Configure(c AggregateConfigurer)

	// New returns a new [AggregateRoot] representing the initial state of an
	// aggregate instance.
	//
	// The engine calls this method to get a "blank slate" when handling the
	// first [Command] for a new instance or when reconstructing an existing
	// instance from its historical [Event] messages.
	New() AggregateRoot

	// RouteCommandToInstance returns the ID of the aggregate instance that c
	// targets.
	//
	// The return value must be a non-empty string that uniquely identifies the
	// target instance. For example, in a shopping cart aggregate, the instance
	// ID might be the customer's ID. RFC 4122 UUIDs are the recommended format.
	//
	// Commands routed to the same instance operate on the same state. There's
	// no need to create an instance in advance - it "exists" once the handler
	// records an [Event] against it.
	//
	// The engine calls this method before handling the [Command]. The
	// implementation must derive the ID from information within c.
	RouteCommandToInstance(c Command) string

	// HandleCommand updates an aggregate instance's state by recording [Event]
	// messages that represent the effects of a [Command].
	//
	// r is the [AggregateRoot] for the instance that the command targets, as
	// determined by [AggregateMessageHandler].RouteCommandToInstance. It
	// reflects the state of the targeted instance after applying its historical
	// events.
	//
	// The implementation must not cause external side-effects or modify r
	// directly. Logic must depend only on information within r, s, and c.
	//
	// The engine atomically persists the events recorded by exactly one
	// successful invocation of this method for each command message. It doesn't
	// guarantee the order, number, or concurrency of those attempts. The
	// implementation doesn't need to perform any synchronization or idempotency
	// checks.
	HandleCommand(
		r AggregateRoot,
		s AggregateCommandScope,
		c Command,
	)
}

// An AggregateRoot is an interface for an application's working representation
// of an aggregate instance used within [AggregateMessageHandler]
// implementations.
//
// It encapsulates business logic and provides a way to inspect the current
// state when making decisions about which events to record. The recorded events
// are the authoritative source of truth, not the AggregateRoot.
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

	// Routes declares the message types that the handler consumes and produces.
	//
	// It accepts routes created by [HandlesCommand] and [RecordsEvent].
	Routes(...AggregateRoute)
}

// AggregateCommandScope represents the context within which an
// [AggregateMessageHandler] handles a [Command] message.
type AggregateCommandScope interface {
	HandlerScope

	// InstanceID returns the ID of the aggregate instance that the [Command]
	// targets, as returned by [AggregateMessageHandler].RouteCommandToInstance.
	InstanceID() string

	// RecordEvent records an [Event] that results from handling the [Command].
	//
	// It applies the event to the aggregate root by calling
	// [AggregateRoot].ApplyEvent, making the state changes visible to the
	// handler immediately.
	//
	// The engine persists all events recorded within this scope in a single
	// atomic operation after the [AggregateMessageHandler] finishes handling
	// the inbound command.
	RecordEvent(Event)
}

// AggregateRoute describes a message type that's routed to or from a
// [AggregateMessageHandler].
type AggregateRoute interface {
	MessageRoute
	isAggregateRoute()
}
