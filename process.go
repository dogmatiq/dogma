package dogma

import (
	"context"
	"time"
)

// ProcessMessageHandler is an interface implemented by the application and used
// by the engine to model business processes.
//
// Many instances of each process type can be created. Each instance is a
// collection of objects that represent the state of the process within the
// application. All manipulation of a process instance is performed via one of
// its constituent objects, known as the "root", and represented by the
// ProcessRoot interface.
//
// Process instances are begun, updated and ended by event messages. The process
// can cause further changes by executing new commands.
//
// Each event message can be routed to many process types, but can target at
// most one instance of each type.
//
// Processes are used to coordinate changes to multiple aggregates, and to
// integrate the domain layer with non-domain concerns.
type ProcessMessageHandler interface {
	// New constructs a new process instance and returns its root.
	//
	// The return value MUST NOT be nil.
	New() ProcessRoot

	// Configure produces a configuration for this handler by calling methods on
	// the configurer c.
	//
	// The implementation MUST allow for multiple calls to Configure(). Each
	// call SHOULD produce the same configuration.
	//
	// The engine MUST call Configure() before calling HandleCommand(). It is
	// RECOMMENDED that the engine only call Configure() once per handler.
	Configure(c ProcessConfigurer)

	// RouteEventToInstance returns the ID of the process instance that is
	// targetted by m.
	//
	// If ok is false, the engine MUST NOT call HandleEvent() with this message.
	//
	// If ok is true, id MUST be a non-empty string. The use of UUIDs for
	// instance identifiers is RECOMMENDED.
	//
	// The engine MUST NOT call RouteEventToInstance() with any message of a
	// type that has not been configured for consumption by a prior call to
	// Configure(). If any such message is passed, the implementation MUST panic
	// with the UnexpectedMessage value.
	RouteEventToInstance(ctx context.Context, m Message) (id string, ok bool, err error)

	// HandleEvent handles an event message.
	//
	// Handling an event message involves inspecting the state of the target
	// process instance to determine what command messages, if any, should be
	// produced.
	//
	// The engine SHOULD provide "at-least-once" delivery guarantees to the
	// handler. That is, the engine should call HandleEvent() with the same
	// command message until a nil error is returned.
	//
	// The engine MUST NOT call HandleEvent() with any message of a type that
	// has not been configured for consumption by a prior call to Configure().
	// If any such message is passed, the implementation MUST panic with the
	// UnexpectedMessage value.
	//
	// The engine MAY provide guarantees about the order in which event messages
	// will be passed to HandleEvent(), however in the interest of engine
	// portability the implementation SHOULD NOT assume that HandleEvent() will
	// be called with events in the same order that they were recorded.
	//
	// The engine MAY call HandleEvent() from multiple goroutines concurrently.
	HandleEvent(ctx context.Context, s ProcessEventScope, m Message) error

	// HandleTimeout handles a timeout message that has been scheduled with
	// ProcessScope.ScheduleTimeout().
	//
	// Timeouts can be used to model time within the business domain. For
	// example, an application might use a timeout to mark an invoice as overdue
	// after some period of non-payment.
	//
	// The engine MUST NOT call HandleTimeout() with any message that was not
	// scheduled by this handler. If any such message is passed, the
	// implementation MUST panic with the UnexpectedMessage value.
	//
	// The engine SHOULD provide "at-least-once" delivery guarantees to the
	// handler. That is, the engine should call HandleTimeout() with the same
	// timeout message until a nil error is returned.
	//
	// The engine MUST NOT call HandleTimeout() before the time at which the
	// timeout message was scheduled. It SHOULD attempt to call HandleTimeout()
	// as soon as the scheduled time is reached.
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
	// Name sets the name of the handler.
	//
	// It MUST be called exactly once within a single call to Configure().
	//
	// Each handler within an application MUST have a unique, non-empty name.
	Name(n string)

	// ConsumesEventType configures the engine to route event messages of the
	// same type as m to the handler.
	//
	// It MUST be called at least once within a call to Configure(). It MUST NOT
	// be called more than once with an event message of the same type.
	//
	// Multiple handlers within an application MAY consume event messages of the
	// same type.
	//
	// The "content" of m MUST NOT be used, inspected, or treated as meaningful
	// in any way, only its runtime type information may be used.
	ConsumesEventType(m Message)

	// ProducesCommandType instructs the engine that the handler executes
	// commands of the same type as m.
	//
	// It MUST be called at least once within a call to Configure(). It MUST NOT
	// be called more than once with a command message of the same type.
	//
	// Multiple handlers within an application MAY produce command messages of
	// the same type.
	//
	// The "content" of m MUST NOT be used, inspected, or treated as meaningful
	// in any way, only its runtime type information may be used.
	ProducesCommandType(m Message)
}

// ProcessEventScope is an interface implemented by the engine and used by the
// application to perform operations within the context of handling an event
// message.
type ProcessEventScope interface {
	// InstanceID returns the ID of the targeted process instance.
	InstanceID() string

	// Begin starts the targeted process instance.
	//
	// It MUST be called before Root(), ExecuteCommand() or ScheduleTimeout()
	// can be called within this scope or the scope of any future message that
	// targets the same instance.
	//
	// It returns true if the targeted instance was begun, or false if
	// the instance had already begun.
	Begin() bool

	// End terminates the targeted process instance.
	//
	// After it has been called Root(), ExecuteCommand() and ScheduleTimeout()
	// MUST NOT be called within this scope or the scope of any future message
	// that targets the same instance.
	//
	// It MUST NOT be called if the instance has not begun.
	//
	// The engine MUST discard any timeout messages associated with this
	// instance.
	//
	// The engine MAY allow re-beginning a process instance that has ended.
	// Callers SHOULD assume that such behavior is unavailable.
	End()

	// Root returns the root of the targeted process instance.
	//
	// It MUST NOT be called if the instance has not begun, or has ended.
	Root() ProcessRoot

	// ExecuteCommand executes a command as a result of the event or timeout
	// message being handled.
	//
	// It MUST NOT be called with a message of any type that has not been
	// configured for production by a prior call to Configure().
	//
	// It MUST NOT be called if the instance has not begun, or has ended.
	ExecuteCommand(m Message)

	// ScheduleTimeout schedules a timeout message to be returned to the process
	// at a specific time.
	//
	// Any pending timeout messages are cancelled when the instance is ended.
	//
	// It MUST NOT be called if the instance has not begun, or has ended.
	ScheduleTimeout(m Message, t time.Time)

	// Log records an informational message within the context of the message
	// that is being handled.
	Log(f string, v ...interface{})
}

// ProcessTimeoutScope is an interface implemented by the engine and used by the
// application to perform operations within the context of handling a timeout
// message.
type ProcessTimeoutScope interface {
	// InstanceID returns the ID of the targeted process instance.
	InstanceID() string

	// End terminates the targeted process instance.
	//
	// After it has been called Root(), ExecuteCommand() and ScheduleTimeout()
	// MUST NOT be called within this scope or the scope of any future message
	// that targets the same instance.
	//
	// It MUST NOT be called if the instance has not begun.
	//
	// The engine MUST discard any timeout messages associated with this
	// instance.
	//
	// The engine MAY allow re-beginning a process instance that has ended.
	// Callers SHOULD assume that such behavior is unavailable.
	End()

	// Root returns the root of the targeted process instance.
	//
	// It MUST NOT be called if the instance has not begun, or has ended.
	Root() ProcessRoot

	// ExecuteCommand executes a command as a result of the event or timeout
	// message being handled.
	//
	// It MUST NOT be called with a message of any type that has not been
	// configured for production by a prior call to Configure().
	//
	// It MUST NOT be called if the instance has not begun, or has ended.
	ExecuteCommand(m Message)

	// ScheduleTimeout schedules a timeout message to be returned to the process
	// at a specific time.
	//
	// Any pending timeout messages are cancelled when the instance is ended.
	//
	// It MUST NOT be called if the instance has not begun, or has ended.
	ScheduleTimeout(m Message, t time.Time)

	// Log records an informational message within the context of the message
	// that is being handled.
	Log(f string, v ...interface{})
}

// StatelessProcessBehavior can be embedded in ProcessMessageHandler
// implementations to indicate that no state is kept between messages.
//
// It provides an implementation of ProcessMessageHandler.New() that always
// returns a StatelessProcessRoot.
type StatelessProcessBehavior struct{}

// New returns StatelessProcessRoot.
func (StatelessProcessBehavior) New() ProcessRoot {
	return StatelessProcessRoot
}

// StatelessProcessRoot is a process root with no state.
//
// It can be returned by a ProcessMessageHandler.New() implementation to
// indicate that no domain state is required beyond the existence/non-existence
// of the process instance.
//
// See also StatelessProcessBehavior, which provides an implementation of
// New() that returns this value.
//
// Engines may use this value as a sentinel, to provide an optimized code path
// when no state is required.
var StatelessProcessRoot ProcessRoot = statelessProcessRoot{}

type statelessProcessRoot struct{}

// NoTimeoutBehavior can be embedded in ProcessMessageHandler implementations to
// indicate that no timeout messages are used.
//
// It provides an implementation of ProcessMessageHandler.HandleTimeout() that always
// panics with the UnexpectedMessage value.
type NoTimeoutBehavior struct{}

// HandleTimeout panic with the UnexpectedMessage value.
func (NoTimeoutBehavior) HandleTimeout(context.Context, ProcessTimeoutScope, Message) error {
	panic(UnexpectedMessage)
}
