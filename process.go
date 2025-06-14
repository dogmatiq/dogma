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

// ProcessRoot is a "marker" interface for the domain-specific state of a
// specific process instance.
//
// The interface is empty to allow use of any types supported by the engine.
type ProcessRoot interface{}

// A ProcessConfigurer configures the engine for use with a specific process
// message handler.
type ProcessConfigurer interface {
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
	// Process handlers support the HandlesEvent(), ExecutesCommand() and
	// SchedulesTimeout() route types.
	Routes(...ProcessRoute)

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

// ProcessEventScope performs engine operations within the context of a call
// to the HandleEvent() method of a [ProcessMessageHandler].
type ProcessEventScope interface {
	// InstanceID returns the ID of the process instance.
	InstanceID() string

	// End signals the end of the process.
	//
	// Ending a process instance destroys its state and cancels any pending
	// timeouts.
	//
	// The process instance ends once HandleEvent() returns. Any future call to
	// ExecuteCommand() or ScheduleTimeout() on this scope prevents the process
	// from ending.
	//
	// "Re-beginning" a process instance that has ended has undefined behavior
	// and is NOT RECOMMENDED.
	End()

	// ExecuteCommand executes a command as a result of the event.
	//
	// Executing a command cancels any prior call to End() on this scope.
	ExecuteCommand(Command)

	// ScheduleTimeout schedules a timeout to occur at a specific time.
	//
	// Ending the process cancels any pending timeouts. Scheduling a timeout
	// cancels any prior call to End() on this scope.
	ScheduleTimeout(Timeout, time.Time)

	// RecordedAt returns the time at which the event occurred.
	RecordedAt() time.Time

	// Now returns the current local time, according to the engine.
	//
	// Use of this method is discouraged. It is preferrable to use information
	// contained within the message, the process root, or the time returned by
	// [ProcessEventScope.RecordedAt], which provides consistent behavior when
	// message delivery is delayed or retried.
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

// ProcessTimeoutScope performs engine operations within the context of a call
// to the HandleTimeout() method of a [ProcessMessageHandler].
type ProcessTimeoutScope interface {
	// InstanceID returns the ID of the process instance.
	InstanceID() string

	// End signals the end of the process.
	//
	// Ending a process instance destroys its state and cancels any pending
	// timeouts.
	//
	// The process instance ends once HandleTimeout() returns. Any future call
	// to ExecuteCommand() or ScheduleTimeout() on this scope prevents the
	// process from ending.
	//
	// "Re-beginning" a process instance that has ended has undefined behavior
	// and is NOT RECOMMENDED.
	End()

	// ExecuteCommand executes a command as a result of the timeout.
	//
	// Executing a command cancels any prior call to End() on this scope.
	ExecuteCommand(Command)

	// ScheduleTimeout schedules a timeout to occur at a specific time.
	//
	// Ending the process cancels any pending timeouts. Scheduling a timeout
	// cancels any prior call to End() on this scope.
	ScheduleTimeout(Timeout, time.Time)

	// ScheduledFor returns the time at which the timeout occured.
	//
	// The time may be before the current time. For example, the engine may
	// deliver timeouts that were "missed" after recovering from downtime.
	ScheduledFor() time.Time

	// Now returns the current local time, according to the engine.
	//
	// Use of this method is discouraged. It is preferrable to use information
	// contained within the message, the process root, or the time returned by
	// [ProcessTimeoutScope.ScheduledFor], which provides consistent behavior
	// when message delivery is delayed or retried.
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

// StatelessProcessRoot is an implementation of [ProcessRoot] for processes that
// do not require any domains-specific state.
//
// [StatelessProcessBehavior] provides a partial implementation of
// [ProcessMessageHandler] that returns this value.
//
// Engines MAY use this value as a sentinel to provide an optimized code path
// when no state is required.
var StatelessProcessRoot ProcessRoot = statelessProcessRoot{}

type statelessProcessRoot struct{}

// StatelessProcessBehavior is an embeddable type for [ProcessMessageHandler]
// that do not have any domain-specific state.
type StatelessProcessBehavior struct{}

// New returns [StatelessProcessRoot].
func (StatelessProcessBehavior) New() ProcessRoot {
	return StatelessProcessRoot
}

// NoTimeoutMessagesBehavior is an embeddable type for [ProcessMessageHandler]
// implementations that do not use [Timeout] messages.
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
	Route
	isProcessRoute()
}
