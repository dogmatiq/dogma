package dogma

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
	// The return value MUST not be empty. [RFC 4122] UUIDs are the RECOMMENDED
	// format for instance IDs.
	RouteCommandToInstance(Command) string

	// HandleCommand executes business logic in response to a command.
	//
	// The handler inspects the root to determine which [Event] messages to
	// record, if any.
	//
	// The handle SHOULD NOT have any side-effects beyond recording events.
	// Specifically, the implementation MUST NOT modify the root directly. Use
	// [AggregateCommandScope.RecordEvent] to record an event that represents
	// the state change. See also [AggregateRoot.ApplyEvent].
	//
	// If this is the first command routed to this instance, the root is the
	// return value of [AggregateMessageHandler.New]. Otherwise, it's the value
	// of the root as it existed after handling the command.
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
	// this aggregate type. This is typically required by event-sourcing engines
	// that sometimes load aggregates into memory by applying their entire
	// history.
	ApplyEvent(Event)
}

// An AggregateConfigurer configures the engine for use with a specific
// aggregate message handler.
//
// See [AggregateMessageHandler.Configure].
type AggregateConfigurer interface {
	// Identity configures the handler's identity.
	//
	// n is a short human-readable name. It MUST be unique within the
	// application. The name MAY change over the handler's lifetime. n MUST
	// contain solely printable, non-space UTF-8 characters.
	//
	// k is a unique key used to associate engine state with the handler. The
	// key SHOULD NOT change over the handler's lifetime. k MUST be a an [RFC
	// 4122] UUID, such as "5195fe85-eb3f-4121-84b0-be72cbc5722f".
	//
	// Use of hard-coded literals for both values is RECOMMENDED.
	Identity(n string, k string)

	// Routes configures the engine to route certain message types to and from
	// the handler.
	//
	// Aggregate handlers support the [HandlesCommand] and [RecordsEvent]
	// route types.
	Routes(...ProcessRoute)

	// ConsumesCommandType configures the engine to route commands of a specific
	// type to the handler.
	//
	// The application's configuration MUST route each command type to a single
	// handler.
	//
	// The command SHOULD be the zero-value of its type; the engine uses the
	// type information, but not the value itself.
	//
	// Deprecated: Use [AggregateConfigurer.Routes] instead.
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
	// Deprecated: Use [AggregateConfigurer.Routes] instead.
	ProducesEventType(Event)
}

// AggregateCommandScope performs engine operations within the context of a call
// to [AggregateMessageHandler.HandleCommand].
type AggregateCommandScope interface {
	// InstanceID returns the ID of the aggregate instance.
	InstanceID() string

	// RecordEvent records the occurrence of an event.
	//
	// It applies the event to the root via [AggregateRoot.ApplyEvent], such
	// that the applied changes are visible to the handler after this method
	// returns.
	//
	// Recording an event cancels any prior call to
	// [AggregateCommandScope.Destroy] on this scope.
	RecordEvent(Event)

	// Destroy signals destruction of the aggregate instance.
	//
	// Destroying a process discards its state. The first command to target a
	// destroyed instance operates on a new root, as returned by
	// [AggregateMessageHandler.New].
	//
	// Destruction occurs once [AggregateCommandScope.HandleCommand] returns.
	// Any future call to [AggregateCommandScope.RecordEvent] on this scope
	// prevents destruction.
	//
	// The precise destruction semantics are engine defined. For example,
	// event-sourcing engines typically do not destroy the record of the
	// aggregate's historical events.
	Destroy()

	// Log records an informational message using [fmt.Printf] formatting.
	Log(format string, args ...any)
}

// AggregateRouteConfigurer configures the engine to route messages for a
// [AggregateMessageHandler].
//
// The engine uses this interface configure its internal routing system.
// Aggregate handlers should use [AggregateConfigurer.Routes] to configure
// their routes.
type AggregateRouteConfigurer interface {
	HandlesCommand(HandlesCommandRoute)
	RecordsEvent(RecordsEventRoute)
}

func (r HandlesCommandRoute) applyToAggregate(v AggregateRouteConfigurer) { v.HandlesCommand(r) }
func (r RecordsEventRoute) applyToAggregate(v AggregateRouteConfigurer)   { v.RecordsEvent(r) }
