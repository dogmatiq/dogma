package dogma

import (
	"context"
	"time"
)

// ProcessMessageHandler is an interface implemented by the application and
// used by the engine to model business processes.
//
// Many instances of each process type can be created. Each instance is a
// collection of objects that represent the state of the process within the
// application. All manipulation of a process instance is performed via one
// of its constituent objects, known as the "root", and represented by the
// ProcessRoot interface.
//
// Process instances are begun, updated and ended by event messages. The process
// can cause further changes by executing new commands.
//
// Each event message can be routed to many process types, but can target at
// most one instance of each type.
//
// Processes are often used to integrate the domain layer with non-domain
// concerns, and as such they often accept and produce both domain messages and
// integration messages.
type ProcessMessageHandler interface {
	// New constructs a new process instance and returns its root.
	New() ProcessRoot

	// Configure configures the behavior of the engine as it relates to this
	// handler.
	//
	// c provides access to the various configuration options, such as
	// specifying which types of event messages are routed to this handler.
	Configure(c ProcessConfigurer)

	// RouteEventToInstance returns the ID of the process instance that is
	// targetted by m.
	//
	// It panics with the UnexpectedMessage value if m is not one of the event
	// types that is routed to this handler via Configure().
	//
	// If ok is false, the message is not routed to this handler at all.
	RouteEventToInstance(ctx context.Context, m Message) (id string, ok bool, err error)

	// HandleEvent handles an event message that has been routed to this
	// handler.
	//
	// Handling an event message involves inspecting the state of the target
	// process instance to determine what command messages, if any, should be
	// produced.
	//
	// s provides access to the operations available within the scope of handling
	// m, such as beginning or ending the targeted instance, accessing its state,
	// sending command messages or scheduling timeouts.
	//
	// This method may manipulate the process's state directly.
	//
	// It panics with the UnexpectedMessage value if m is not one of the event
	// types that is routed to this handler via Configure().
	HandleEvent(ctx context.Context, s ProcessEventScope, m Message) error

	// HandleTimeout handles a timeout message that has been scheduled with
	// ProcessScope.ScheduleTimeout().
	//
	// Timeouts can be used to model time within the domain. For example, an
	// application might use a timeout to mark an invoice as overdue after some
	// period of non-payment.
	//
	// Handling a timeout is much like handling an event in that much the same
	// operations are available to the handler via s.
	//
	// This method may manipulate the process's state directly.
	//
	// If m was not expected by the handler the implementation must panic with an
	// UnexpectedMessage value.
	HandleTimeout(ctx context.Context, s ProcessTimeoutScope, m Message) error
}

// ProcessRoot is an interface implemented by the application and used by
// the engine to represent the state of a process instance.
type ProcessRoot interface {
}

// ProcessConfigurer is an interface implemented by the engine and used by
// the application to configure options related to a ProcessMessageHandler.
//
// It is passed to ProcessMessageHandler.Configure(), typically upon
// initialization of the engine.
//
// In the context of this interface, "the handler" refers to the handler on
// which Configure() has been called.
type ProcessConfigurer interface {
	// RouteEventType configures the engine to route events of the same type as m
	// to the handler.
	RouteEventType(m Message)
}

// ProcessEventScope is an interface implemented by the engine and used by the
// application to perform operations within the context of handling an event
// message.
type ProcessEventScope interface {
	// InstanceID is the ID of the targeted process instance.
	InstanceID() string

	// Begin starts the targeted process instance.
	//
	// It must be called before Root(), ExecuteCommand() or ScheduleTimeout() can
	// be called within this scope or the scope of any future event or timeout that
	// targets the same instance.
	//
	// It returns true if the targeted instance was begun, or false if
	// the instance had already begun.
	Begin() bool

	// End terminates the targeted process instance.
	//
	// After it has been called none of Root(), ExecuteCommand() or
	// ScheduleTimeout() can be called within this scope or the scope of any future
	// event or timeout that targets the same instance.
	//
	// It panics if the target instance has not been begun.
	//
	// The precise semantics of ending a process instance are implementation
	// defined. The engine is not required to allow re-beginning a process
	// instance that has been ended.
	End()

	// Root returns the root of the targeted process instance.
	//
	// It panics if the instance has not been begun, or was begun but has
	// subsequently been ended.
	Root() ProcessRoot

	// ExecuteCommand executes a command as a result of the event or timeout
	// message being handled.
	//
	// It panics if the instance has not been begun, or was begun but has
	// subsequently been ended.
	ExecuteCommand(m Message)

	// ScheduleTimeout schedules a timeout message to be returned to the process
	// at a specific time.
	//
	// Any pending timeout messages are cancelled when the instance is ended.
	//
	// It panics if the instance has not been begun, or was begun but has
	// subsequently been ended.
	ScheduleTimeout(m Message, t time.Time)

	// Log records an informational message within the context of the event being
	// handled.
	//
	// The log message should be worded such that it makes sense to anyone familiar
	// with the business domain.
	Log(f string, v ...interface{})
}

// ProcessTimeoutScope is an interface implemented by the engine and used by the
// application to perform operations within the context of handling a timeout
// message.
type ProcessTimeoutScope interface {
	// InstanceID is the ID of the targeted process instance.
	InstanceID() string

	// End terminates the targeted process instance.
	//
	// After it has been called none of Root(), ExecuteCommand() or
	// ScheduleTimeout() can be called within this scope or the scope of any future
	// event or timeout that targets the same instance.
	//
	// It panics if the target instance has not been begun.
	//
	// The precise semantics of ending a process instance are implementation
	// defined. The engine is not required to allow re-beginning a process
	// instance that has been ended.
	End()

	// Root returns the root of the targeted process instance.
	//
	// It panics if the instance has not been begun, or was begun but has
	// subsequently been ended.
	Root() ProcessRoot

	// ExecuteCommand executes a command as a result of the event or timeout
	// message being handled.
	//
	// It panics if the instance has not been begun, or was begun but has
	// subsequently been ended.
	ExecuteCommand(m Message)

	// ScheduleTimeout schedules a timeout message to be returned to the process
	// at a specific time.
	//
	// Any pending timeout messages are cancelled when the instance is ended.
	//
	// It panics if the instance has not been begun, or was begun but has
	// subsequently been ended.
	ScheduleTimeout(m Message, t time.Time)

	// Log records an informational message within the context of the timeout being
	// handled.
	Log(f string, v ...interface{})
}

// StatelessProcess is an embeddable type that provides an implementation of
// ProcessMessageHandler.New() that always returns a StatelessProcessRoot.
type StatelessProcess struct{}

// New returns StatelessProcessRoot.
func (StatelessProcess) New() ProcessRoot {
	return StatelessProcessRoot
}

// StatelessProcessRoot is a process root with no state.
//
// It can be returned by a ProcessMessageHandler.New() implementation to
// indicate that no domain state is required beyond the existence/non-existence
// of the process instance.
var StatelessProcessRoot ProcessRoot = statelessProcessRoot{}

type statelessProcessRoot struct{}

// NoTimeouts is an embeddable type that provides an implementation of
// Process.HandleTimeout() that always panics with the UnexpectedMessage value.
type NoTimeouts struct{}

// HandleTimeout panic with the UnexpectedMessage value.
func (NoTimeouts) HandleTimeout(context.Context, ProcessTimeoutScope, Message) error {
	panic(UnexpectedMessage)
}
