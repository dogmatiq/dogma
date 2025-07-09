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
	// empty. RFC 4122 UUIDs are the RECOMMENDED format for instance IDs.
	//
	// A process instance begins the first time it receives an event.
	RouteEventToInstance(context.Context, Event) (id string, ok bool, err error)

	// HandleEvent begins or continues the process in response to an event.
	//
	// The handler inspects the root to determine which commands to execute, if
	// any. It may also schedule timeouts to "wake" the process at a later time.
	//
	// If this is the first event routed to this instance, the root is the
	// return value of New(). Otherwise, it's the value of the root as it
	// existed after handling the last event or timeout.
	//
	// The engine MAY provide specific guarantees about the order in which it
	// supplies events to the handler. To maximize portability across engines,
	// the handler SHOULD NOT assume any specific ordering. The engine MAY call
	// this method concurrently from separate goroutines or operating system
	// processes.
	HandleEvent(context.Context, ProcessRoot, ProcessEventScope, Event) error

	// HandleTimeout continues the process in response to a timeout.
	//
	// The handler inspects the root to determine which commands to execute, if
	// any. It may also schedule timeout messages to "wake" the process at a
	// later time.
	//
	// The engine MUST NOT call this method before the timeout's scheduled time.
	// The engine MAY call this method concurrently from separate goroutines or
	// operating system processes.
	HandleTimeout(context.Context, ProcessRoot, ProcessTimeoutScope, Timeout) error
}

// A ProcessRoot is an interface for an application's representation of a
// process instance used within [ProcessMessageHandler] implementations.
//
// It encapsulates workflow logic and provides a way to inspect the current
// state when making decisions about which commands to execute and which
// timeouts to schedule.
//
// This interface is currently equivalent to [any], but is a distinct type to
// allow future extensions without breaking compatibility.
type ProcessRoot any

// ProcessConfigurer is the interface that a [ProcessMessageHandler] uses to
// declare its configuration.
//
// The engine provides the implementation to [ProcessMessageHandler].Configure
// during startup.
type ProcessConfigurer interface {
	HandlerConfigurer

	// Routes declares the message types that the handler consumes and produces.
	//
	// It accepts routes created by [HandlesEvent], [ExecutesCommand], and
	// [SchedulesTimeout].
	Routes(...ProcessRoute)
}

// ProcessScope represents the context within which a [ProcessMessageHandler]
// handles a message.
//
// Each kind of message handled by a process message handler has a corresponding
// scope type that extends this interface:
//
//   - [ProcessEventScope]
//   - [ProcessTimeoutScope]
type ProcessScope interface {
	HandlerScope

	// InstanceID returns the ID of the process instance that the message
	// targets.
	//
	// When handling an [Event] message, it returns the ID produced by
	// [ProcessMessageHandler].RouteEventToInstance during routing.
	//
	// When handling a [Timeout] message, it returns the ID of the instance that
	// scheduled the timeout.
	InstanceID() string

	// End signals the end of a process.
	//
	// The engine discards the instance's state, cancels any pending [Timeout]
	// messages. It ignores any future messages that target the ended instance.
	End()

	// ExecuteCommand submits a [Command] for execution.
	//
	// The engine persists all commands and timeouts produced within this scope
	// in a single atomic operation after the [ProcessMessageHandler] finishes
	// handling the inbound message. If the handler returns a non-nil error, the
	// engine discards the messages.
	//
	// This method panics if the process instance has ended.
	ExecuteCommand(Command)

	// ScheduleTimeout schedules a [Timeout] message to occur at the specified
	// time.
	//
	// The engine persists all commands and timeouts produced within this scope
	// in a single atomic operation after the [ProcessMessageHandler] finishes
	// handling the inbound message. If the handler returns a non-nil error, the
	// engine discards the messages.
	//
	// This method panics if the process instance has ended.
	ScheduleTimeout(Timeout, time.Time)
}

// ProcessEventScope represents the context within which a
// [ProcessMessageHandler] handles an [Event] message.
type ProcessEventScope interface {
	ProcessScope

	// RecordedAt returns the time at which the inbound [Event] occurred.
	RecordedAt() time.Time
}

// ProcessTimeoutScope represents the context within which a
// [ProcessMessageHandler] handles a [Timeout] message.
type ProcessTimeoutScope interface {
	ProcessScope

	// ScheduledFor returns the time at which the timeout occurred.
	//
	// Even though the engine attempts to deliver timeouts at their scheduled
	// time, it may deliver them later when recovering from downtime or retrying
	// after a failure.
	ScheduledFor() time.Time
}

// ProcessRoute describes a message type that's routed to or from a
// [ProcessMessageHandler].
type ProcessRoute interface {
	MessageRoute
	isProcessRoute()
}

// NoTimeoutMessagesBehavior is an embeddable type for [ProcessMessageHandler]
// implementations that don't use [Timeout] messages.
//
// Embed this type in a [ProcessMessageHandler] to signal that the handler
// doesn't schedule timeouts and to avoid boilerplate code that's never
// used.
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

// StatelessProcessBehavior is an embeddable type for [ProcessMessageHandler]
// implementations that don't maintain any state.
//
// Embed this type in a [ProcessMessageHandler] to signal that the handler is
// stateless and to avoid boilerplate code that's never used.
type StatelessProcessBehavior struct{}

// New returns [StatelessProcessRoot].
func (StatelessProcessBehavior) New() ProcessRoot {
	return StatelessProcessRoot
}

// StatelessProcessRoot is an empty [ProcessRoot] for processes that don't
// maintain any state.
//
// Embed [StatelessProcessBehavior] in a [ProcessMessageHandler] to use this
// type as its [ProcessRoot] implementation.
//
// The engine may provide optimized persistence for stateless processes that use
// this type.
var StatelessProcessRoot ProcessRoot = statelessProcessRoot{}

type statelessProcessRoot struct{}
