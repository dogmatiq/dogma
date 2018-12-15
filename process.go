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
// Process instances are begun, updated and ended by domain event messages,
// typically those recorded by an aggregate. The process can cause further
// changes by executing new commands.
//
// Each event message can be routed to many process types, but can target at
// most one instance of each type.
type ProcessMessageHandler interface {
	// New constructs a new process instance and returns its root.
	New() ProcessRoot

	// Configure configures the behavior of the engine as it relates to this
	// handler.
	//
	// c provides access to the various configuration options, such as
	// specifying which event types are routed to this handler.
	Configure(c ProcessConfigurer)

	// RouteEventToInstance returns the ID of the process instance that is
	// targetted by m.
	//
	// It panics with the UnexpectedMessage value if m is not one of the command
	// types that is routed to this handler via Configure().
	//
	// If ok is false, the message is not routed to this handler at all.
	RouteEventToInstance(ctx context.Context, m Message) (id string, ok bool, err error)

	// HandleEvent handles an event message that has been routed to this
	// handler.
	//
	// Handling an event involves inspecting the state of the event's target
	// process instance to determine what commands, if any, should be executed.
	//
	// s provides access to the operations available within the scope of handling
	// m, such as beginning or ending the targeted instance, accessing its
	// state, executing commands or scheduling timeouts.
	//
	// This method may manipulate the process's state directly.
	//
	// If m was not expected by the handler the implementation must panic with an
	// UnexpectedMessage value.
	HandleEvent(ctx context.Context, s ProcessScope, m Message) error

	// HandleTimeout handles a timeout message that has been scheduled with
	// ProcessScope.ScheduleTimeout().
	//
	// Timeouts can be used to model time within the domain. For example, an
	// application might use a timeout to mark an invoice as overdue after some
	// period of non-payment.
	//
	// Handling a timeout is much like handling an event in that the same
	// operations are available to the handler via s.
	//
	// This method may manipulate the process's state directly.
	//
	// If m was not expected by the handler the implementation must panic with an
	// UnexpectedMessage value.
	HandleTimeout(ctx context.Context, s ProcessScope, m Message) error
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

// ProcessScope is an interface implemented by the engine and used by the
// application to perform operations within the context of handling an event or
// timeout message.
//
// In the context of this interface, "the message" refers to the message being
// handled and "the instance" refers to the process instance that is targeted
// by that message. This message may either be an event, or a timeout message.
type ProcessScope interface {
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

	// Log records an informational message within the context of the event or
	// timeout being handled.
	//
	// The log message should be worded such that it makes sense to anyone familiar
	// with the business domain.
	Log(f string, v ...interface{})
}

// StatelessProcess is a process root with no state.
//
// It can be returned by a ProcessMessageHandler.New() implementation to
// indicate that no domain state is required beyond the existence/non-existence
// of the process instance.
var StatelessProcess ProcessRoot = statelessProcess{}

type statelessProcess struct{}
