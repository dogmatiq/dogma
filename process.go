package dogma

import (
	"context"
	"errors"
	"time"
)

// A ProcessMessageHandler encapsulates an application's "workflow" logic -
// stateful decision-making that spans multiple [Command] messages.
//
// It handles [Event] messages and executes [Command] to enact further
// application state changes. It may also schedule [Timeout] messages to perform
// actions at specific times. For example, to send a reminder if a customer
// hasn't completed the checkout process within one hour.
//
// Each process message handler typically manages multiple instances, where each
// instance represents a distinct occurrence of the process. For example, a
// shopping cart checkout process might use one instance per customer.
//
// Process message handlers coordinate state changes that involve some
// combination of multiple aggregate instances, integrations with external
// systems, and time-based logic. As a general rule, they should implement only
// workflow logic. For example, a process might decide when to refund a
// customer's purchase, but shouldn't calculate the refund amount or interact
// directly with the payment processor.
type ProcessMessageHandler interface {
	// Configure declares the handler's configuration by calling methods on c.
	//
	// The configuration includes the handler's identity and message routes.
	//
	// The engine calls this method at least once during startup. It must
	// produce the same configuration each time it's called.
	Configure(c ProcessConfigurer)

	// New returns a new [ProcessRoot] representing the initial state of a
	// process instance.
	//
	// The engine calls this method to get a "blank slate" when handling the
	// first [Event] for a new instance. Unlike aggregates, the engine doesn't
	// reconstruct process state from historical events.
	//
	// Not all processes maintain state. Embed [StatelessProcessBehavior] in the
	// handler implementation to indicate that the process is stateless.
	New() ProcessRoot

	// RouteEventToInstance returns the ID of the process instance that e
	// targets.
	//
	// If ok is false, the handler ignores the event. Otherwise, the returned ID
	// must be a non-empty string that uniquely identifies the target instance.
	// For example, in a shopping cart checkout process, the instance ID might
	// be the customer's ID. RFC 9562 UUIDs are the recommended format.
	//
	// Events routed to the same instance operate on the same state. There's no
	// need to create an instance in advance - it "exists" once the handler
	// modifies its [ProcessRoot], executes a [Command], or schedules a
	// [Timeout] against it.
	//
	// The engine calls this method before handling the [Event]. The
	// implementation may query external data - such as the application's
	// projections - but this isn't recommended. Wherever possible, it should
	// derive the ID from information within e.
	//
	// If the process instance identified by the returned ID has ended, the
	// engine ignores the event.
	RouteEventToInstance(
		ctx context.Context,
		e Event,
	) (id string, ok bool, err error)

	// HandleEvent begins or advances a process in response to an [Event]
	// message.
	//
	// r is the [ProcessRoot] for the instance that the event targets, as
	// determined by [ProcessMessageHandler].RouteEventToInstance. It reflects
	// the state of the targeted instance after handling any prior [Event] or
	// [Timeout] messages.
	//
	// The implementation may update r directly, execute [Command] messages,
	// schedule [Timeout] messages, or end the process. It may query external
	// data - such as the application's projections - but this isn't
	// recommended. Wherever possible, logic should depend solely on information
	// within r, s, and e.
	//
	// The engine atomically persists the state changes, events, and timeouts
	// produced by exactly one successful invocation of this method for each
	// event message. It doesn't guarantee the order, number, or concurrency of
	// those attempts. Generally, the implementation doesn't need to perform any
	// synchronization or idempotency checks.
	//
	// The engine delivers all [Event] messages recorded within a single scope
	// in the order they occurred. It also preserves the order of events from a
	// single aggregate instance, even across scopes. It doesn't guarantee the
	// relative delivery order of events from different handlers or aggregate
	// instances.
	HandleEvent(
		ctx context.Context,
		r ProcessRoot,
		s ProcessEventScope,
		e Event,
	) error

	// HandleTimeout advances a process in response to a [Timeout] message.
	//
	// r is the [ProcessRoot] for the instance that scheduled the timeout. It
	// reflects the state of the targeted instance after handling any prior
	// [Event] or [Timeout] messages.
	//
	// The implementation may update r directly, execute [Command] messages,
	// schedule [Timeout] messages, or end the process. It may query external
	// data - such as the application's projections - but this isn't
	// recommended. Wherever possible, logic should depend solely on information
	// within r, s, and t.
	//
	// The engine atomically persists the state changes, events, and timeouts
	// produced by exactly one successful invocation of this method for each
	// timeout message. It doesn't guarantee the order, number, or concurrency
	// of those attempts. Generally, the implementation doesn't need to perform
	// any synchronization or idempotency checks.
	//
	// The engine attempts to deliver timeout messages at their scheduled time.
	// It may deliver them later when recovering from downtime or retrying after
	// a failure. It doesn't guarantee the relative delivery order of timeout
	// messages with the same scheduled time.
	//
	// Not all processes use timeouts. Embed [NoTimeoutMessagesBehavior] in the
	// handler implementation to indicate that timeout messages aren't used.
	HandleTimeout(
		ctx context.Context,
		r ProcessRoot,
		s ProcessTimeoutScope,
		t Timeout,
	) error
}

// A ProcessRoot is an interface for an application's representation of a
// process instance used within [ProcessMessageHandler] implementations.
//
// It encapsulates process logic and provides a way to inspect the current state
// when making decisions about which commands to execute and which timeouts to
// schedule.
//
// This interface is currently equivalent to [any], but is a distinct type to
// allow future extensions without breaking compatibility.
type ProcessRoot interface {
	// MarshalBinary returns a binary representation of the process instsance's
	// current state.
	MarshalBinary() ([]byte, error)

	// UnmarshalBinary populates the process instance's state from its binary
	// representation.
	//
	// The implementation must clone the data if it is used after returning.
	UnmarshalBinary(data []byte) error
}

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

func (statelessProcessRoot) MarshalBinary() ([]byte, error) {
	return nil, nil
}

func (statelessProcessRoot) UnmarshalBinary(data []byte) error {
	if len(data) != 0 {
		return errors.New("cannot unmarshal non-empty data into stateless process")
	}
	return nil
}
