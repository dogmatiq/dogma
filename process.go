package dogma

import (
	"context"
	"time"
)

// A ProcessMessageHandler models a business process.
//
// Processes are useful for coordinating changes across aggregate instances and
// for modeling processes that include time-based logic.
//
// [Event] messages can begin, advance or end a process. The process causes
// changes within the application by executing [Command] messages.
//
// Processes are stateful. Am application typically uses multiple instances of a
// process, each with its own state. For example, an e-commerce application may
// use one instance of the "checkout" process for each customer's shopping cart.
//
// The state of each instance is application-defined. Often it's a tree of
// related entities and values. The [ProcessRoot] interface represents the "root"
// entity through which the handler accesses the instance's state.
//
// Processes can also schedule [Timeout] messages. Timeouts model time in the
// business domain. For example, a timeout could trigger an email to a customer
// who added a product to their shopping cart but did not pay within one hour.
//
// Process handlers SHOULD NOT directly perform write operations such as
// updating a database or invoking any API that does so.
type ProcessMessageHandler interface {
	// Configure describes the handler's configuration to the engine.
	Configure(ProcessConfigurer)

	// New returns a process root instance in its initial state.
	//
	// The return value MUST NOT be nil. It MAY be the zero-value of the root's
	// underlying type.
	//
	// Each call SHOULD return the same type and initial state.
	New() ProcessRoot

	// RouteEventToInstance returns the ID of the instance that handles a
	// specific event.
	//
	// If ok is false, the process ignores this event. Otherwise, id MUST not be
	// empty. [RFC 4122] UUIDs are the RECOMMENDED format for instance IDs.
	//
	// A process instance begins the first time it receives an event.
	RouteEventToInstance(context.Context, Event) (id string, ok bool, err error)

	// HandleEvent begins or continues the process in response to an event.
	//
	// The handler inspects the root to determine which [Command] messages to
	// execute, if any. It may also schedule [Timeout] messages to "wake" the
	// process at a later time.
	//
	// If this is the first event routed to this instance, the root is the
	// return value of [ProcessMessageHandler.New]. Otherwise, it's the value of
	// the root as it existed after handling the last event or timeout.
	//
	// The engine MAY provide specific guarantees about the order in which it
	// supplies events to the handler. To maximize portability across engines,
	// the handler SHOULD NOT assume any specific ordering. The engine MAY call
	// this method concurrently from separate goroutines or operating system
	// processes.
	//
	// The implementation SHOULD NOT impose a context deadline. Implement the
	// [ProjectionMessageHandler.TimeoutHint] method instead.
	HandleEvent(context.Context, ProcessRoot, ProcessEventScope, Event) error

	// HandleTimeout continues the process in response to a timeout.
	//
	// The handler inspects the root to determine which [Command] messages to
	// execute, if any. It may also schedule timeout messages to "wake" the
	// process at a later time.
	//
	// The engine MUST NOT call this method before the timeout's scheduled time.
	// The engine MAY call this method concurrently from separate goroutines or
	// operating system processes.
	//
	// The implementation SHOULD NOT impose a context deadline. Implement the
	// [ProjectionMessageHandler.TimeoutHint] method instead.
	HandleTimeout(context.Context, ProcessRoot, ProcessTimeoutScope, Timeout) error

	// TimeoutHint returns a suitable duration for handling the given message.
	//
	// The duration SHOULD be as short as possible. If no hint is available it
	// MUST be zero.
	//
	// In this context, "timeout" refers to a deadline, not a [Timeout] message.
	//
	// See [NoTimeoutHintBehavior].
	TimeoutHint(Message) time.Duration
}

// ProcessRoot is a "marker" interface for the domain-specific state of a
// specific process instance.
//
// The interface is empty to allow use of any types supported by the engine.
type ProcessRoot interface{}

// A ProcessConfigurer configures the engine for use with a specific process
// message handler.
//
// See [ProcessMessageHandler.Configure].
type ProcessConfigurer interface {
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
	// Process handlers support the [HandlesEvent], [ExecutesCommand] and
	// [SchedulesTimeout] route types.
	Routes(...ProcessRoute)

	// ConsumesEventType configures the engine to route events of a specific
	// type to the handler.
	//
	// The event SHOULD be the zero-value of its type; the engine uses the type
	// information, but not the value itself.
	//
	// Deprecated: Use [ProcessConfigurer.Routes] instead.
	ConsumesEventType(Event)

	// ProducesCommandType configures the engine to use the handler as the
	// source of events of a specific type.
	//
	// The command SHOULD be the zero-value of its type; the engine uses the
	// type information, but not the value itself.
	//
	// Deprecated: Use [ProcessConfigurer.Routes] instead.
	ProducesCommandType(Command)

	// SchedulesTimeoutType configures the engine to allow this handler to
	// schedule timeouts of a specific type.
	//
	// The timeout SHOULD be the zero-value of its type; the engine uses the
	// type information, but not the value itself.
	//
	// Deprecated: Use [ProcessConfigurer.Routes] instead.
	SchedulesTimeoutType(Timeout)
}

// ProcessEventScope performs engine operations within the context of a call
// to [ProcessMessageHandler.HandleEvent].
type ProcessEventScope interface {
	// InstanceID returns the ID of the process instance.
	InstanceID() string

	// End signals the end of the process.
	//
	// Ending a process instance destroys its state and cancels any pending
	// timeouts.
	//
	// The process instance ends once [ProcessMessageHandler.HandleEvent]
	// returns. Any future call to [ProcessEventScope.ExecuteCommand] or
	// [ProcessEventScope.ScheduleTimeout] on this scope prevents the process
	// from ending.
	//
	// "Re-beginning" a process instance that has ended has undefined behavior
	// and is NOT RECOMMENDED.
	End()

	// ExecuteCommand executes a command as a result of the event.
	//
	// Executing a command cancels any prior call to [ProcessEventScope.End]
	// on this scope.
	ExecuteCommand(Command)

	// ScheduleTimeout schedules a timeout to occur at a specific time.
	//
	// Ending the process cancels any pending timeouts. Scheduling a timeout
	// cancels any prior call to [ProcessEventScope.End] on this scope.
	ScheduleTimeout(Timeout, time.Time)

	// RecordedAt returns the time at which the event occurred.
	RecordedAt() time.Time

	// Log records an informational message using [fmt.Printf] formatting.
	Log(format string, args ...any)
}

// ProcessTimeoutScope performs engine operations within the context of a call
// to [ProcessMessageHandler.HandleTimeout].
type ProcessTimeoutScope interface {
	// InstanceID returns the ID of the process instance.
	InstanceID() string

	// End signals the end of the process.
	//
	// Ending a process instance destroys its state and cancels any pending
	// timeouts.
	//
	// The process instance ends once [ProcessMessageHandler.HandleTimeout]
	// returns. Any future call to [ProcessTimeoutScope.ExecuteCommand] or
	// [ProcessTimeoutScope.ScheduleTimeout] on this scope prevents the process
	// from ending.
	//
	// "Re-beginning" a process instance that has ended has undefined behavior
	// and is NOT RECOMMENDED.
	End()

	// ExecuteCommand executes a command as a result of the timeout.
	//
	// Executing a command cancels any prior call to [ProcessTimeoutScope.End]
	// on this scope.
	ExecuteCommand(Command)

	// ScheduleTimeout schedules a timeout to occur at a specific time.
	//
	// Ending the process cancels any pending timeouts. Scheduling a timeout
	// cancels any prior call to [ProcessTimeoutScope.End] on this scope.
	ScheduleTimeout(Timeout, time.Time)

	// ScheduledFor returns the time at which the timeout occured.
	//
	// The time may be before the current time. For example, the engine may
	// deliver timeouts that were "missed" after recovering from downtime.
	ScheduledFor() time.Time

	// Log records an informational message using [fmt.Printf] formatting.
	Log(format string, args ...any)
}

// StatelessProcessRoot is an implementation of [ProcessRoot] for processes that
// do not require any domains-specific state.
//
// [StatelessProcessBehavior] provides an implementation of
// [ProcessMessageHandler.New] that returns this value.
//
// Engines MAY use this value as a sentinel to provide an optimized code path
// when no state is required.
var StatelessProcessRoot ProcessRoot = statelessProcessRoot{}

type statelessProcessRoot struct{}

// StatelessProcessBehavior is an embeddable type for [ProcessMessageHandler]
// that do not have any domain-specific state.
//
// It provides an implementation of [ProcessMessageHandler.New] that
// always returns [StatelessProcessRoot].
type StatelessProcessBehavior struct{}

// New returns [StatelessProcessRoot].
func (StatelessProcessBehavior) New() ProcessRoot {
	return StatelessProcessRoot
}

// NoTimeoutMessagesBehavior is an embeddable type for [ProcessMessageHandler]
// implementations that do not use [Timeout] messages.
//
// It provides an implementation of [ProcessMessageHandler.HandleTimeout] that
// always panics with the [UnexpectedMessage] value.
type NoTimeoutMessagesBehavior struct{}

// HandleTimeout panics with the [UnexpectedMessage] value.
func (NoTimeoutMessagesBehavior) HandleTimeout(
	context.Context,
	ProcessRoot,
	ProcessTimeoutScope,
	Timeout,
) error {
	panic(UnexpectedMessage)
}

// ProcessRoute describes a message type that's routed to or from a
// [ProcessMessageHandler].
type ProcessRoute interface {
	applyToProcess(ProcessRouteConfigurer)
}

// ProcessRouteConfigurer configures the engine to route messages for a
// [ProcessMessageHandler].
//
// The engine uses this interface configure its internal routing system. Process
// handlers should use [ProcessConfigurer.Routes] to configure their routes.
type ProcessRouteConfigurer interface {
	HandlesEvent(HandlesEventRoute)
	ExecutesCommand(ExecutesCommandRoute)
	SchedulesTimeout(SchedulesTimeoutRoute)
}

func (r HandlesEventRoute) applyToProcess(v ProcessRouteConfigurer)     { v.HandlesEvent(r) }
func (r ExecutesCommandRoute) applyToProcess(v ProcessRouteConfigurer)  { v.ExecutesCommand(r) }
func (r SchedulesTimeoutRoute) applyToProcess(v ProcessRouteConfigurer) { v.SchedulesTimeout(r) }
