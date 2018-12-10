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
// changes by executing new commands. Each event message can be routed to many
// process types, but can target at most one instance of each type.
type ProcessMessageHandler interface {
	// New constructs a new process instance and returns its root.
	New() ProcessRoot

	// RouteEvent indicates whether a specific type or instance of a message
	// should be routed to this handler as an event.
	//
	// If p is false, then m is an event message that has been made available to
	// the application. If m should be routed to this handler, the implementation
	// sets ok to true and id to the ID of the process instance that the event
	// targets. The process instance need not already have begun in order for an
	// event to target it. id must not be empty if ok is true.
	//
	// If p is true, then the engine is performing a "routing probe". In this case
	// m is a non-nil, zero-value message. The implementation sets ok to true if
	// messages of the same type as m should be routed to this message handler when
	// they occur. The id output parameter is unused.
	RouteEvent(ctx context.Context, m Message, p bool) (id string, ok bool, err error)

	// HandleEvent handles an event message that has been routed to this
	// handler.
	//
	// Handling an event involves inspecting the state of the event's target
	// process instance to determine what commands, if any, should be executed.
	//
	// s provides access to the operations available within the scope of handling
	// m, such as beginning or ending the targetted instance, accessing its
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
	// application might use a timeout mark an invoice as overdue after some
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

// ProcessScope is an interface implemented by the engine and used by the
// application to perform operations within the context of handling an event or
// timeout message.
//
// In the context of this interface, "the message" refers to the message being
// handled and "the instance" refers to the process instance that is targetted
// by that message. This message may either be an event, or a timeout message.
type ProcessScope interface {
	// InstanceID is the ID of the targetted process instance.
	InstanceID() string

	// Begin starts the targetted process instance.
	//
	// It must be called before Root(), ExecuteCommand() or ScheduleTimeout() can
	// be called within this scope or the scope of any future event or timeout that
	// targets the same instance.
	//
	// It returns true if the targetted instance was begun, or false if
	// the instance had already begun.
	Begin() bool

	// End terminates the targetted process instance.
	//
	// After it has been called none of Root(), ExecuteCommand() or
	// ScheduleTimeout() can be called within this scope or the scope of any future
	// event or timeout that targets the same instance.
	//
	// It panics if the target instance does has not been begun.
	//
	// The precise semantics of ending a process instance are implementation
	// defined. The engine is not required to allow re-beginning a process
	// instance that has been ended.
	End()

	// Root returns the root of the targetted process instance.
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
	// Any pending timeout messages are cancelled with the instance is ended.
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
