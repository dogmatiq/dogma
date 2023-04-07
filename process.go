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
//
// Processes SHOULD NOT perform any kind of "write" operation directly, such as
// updating a database or invoking an API operation that causes a state change.
// Any such state changes should be communicated via a command message instead.
type ProcessMessageHandler interface {
	// New constructs a new process instance initialized with any
	// default values and returns the process root.
	//
	// Repeated calls SHOULD return a value that is of the same type and
	// initialized in the same way. The return value MUST NOT be nil.
	New() ProcessRoot

	// Configure describes the handler's configuration to the engine.
	Configure(c ProcessConfigurer)

	// RouteEventToInstance returns the ID of the process instance that is
	// targeted by e.
	//
	// If ok is false, the engine MUST NOT call HandleEvent() with this message.
	//
	// If ok is true, id MUST be a non-empty string. The use of UUIDs for
	// instance identifiers is RECOMMENDED.
	//
	// A process instance is considered to begin the first time an event is
	// routed to it.
	//
	// The engine MUST NOT call RouteEventToInstance() with any message of a
	// type that has not been configured for consumption by a prior call to
	// Configure(). If any such message is passed, the implementation MUST panic
	// with the UnexpectedMessage value.
	RouteEventToInstance(ctx context.Context, e Event) (id string, ok bool, err error)

	// HandleEvent handles an event message.
	//
	// Handling an event message involves inspecting the state of the target
	// process instance (via the process root r) to determine what command
	// messages, if any, should be produced.
	//
	// The engine MUST provide a ProcessRoot, r, the value of which is
	// equivalent to the value of r as it existed after the last call to
	// HandleEvent() or HandleTimeout() for the targeted instance.
	//
	// If this is the first event to target this instance (or the first event to
	// do so since s.End() was last used to end the instance), r MUST be
	// equivalent to the result of New().
	//
	// The engine SHOULD provide "at-least-once" delivery guarantees to the
	// handler. That is, the engine should call HandleEvent() with the same
	// event message until a nil error is returned.
	//
	// The supplied context parameter SHOULD have a deadline. The implementation
	// SHOULD NOT impose its own deadline. Instead a suitable timeout duration
	// can be suggested to the engine via the handler's TimeoutHint() method.
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
	HandleEvent(ctx context.Context, r ProcessRoot, s ProcessEventScope, e Event) error

	// HandleTimeout handles a timeout message that has been scheduled with
	// ProcessScope.ScheduleTimeout().
	//
	// The engine MUST provide a ProcessRoot, r, the value of which is
	// equivalent to the value of r as it existed after the last call to
	// HandleEvent() or HandleTimeout() for the targeted instance.
	//
	// Timeouts can be used to model time within the business domain. For
	// example, an application might use a timeout to mark an invoice as overdue
	// after some period of non-payment.
	//
	// The supplied context parameter SHOULD have a deadline. The implementation
	// SHOULD NOT impose its own deadline. Instead a suitable timeout duration
	// can be suggested to the engine via the handler's TimeoutHint() method.
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
	HandleTimeout(ctx context.Context, r ProcessRoot, s ProcessTimeoutScope, t Timeout) error

	// TimeoutHint returns a suitable duration for handling the given message.
	//
	// The duration SHOULD be as short as possible. If no hint is available it
	// MUST be zero.
	//
	// See [NoTimeoutHintBehavior].
	TimeoutHint(m Message) time.Duration
}

// ProcessRoot is an interface implemented by the application and used by
// the engine to represent the state of a process instance.
type ProcessRoot interface {
}

// A ProcessConfigurer configures the engine for use with a specific process
// message handler.
//
// See [ProcessMessageHandler.Configure]().
type ProcessConfigurer interface {
	// Identity configures how the engine identifies the handler.
	//
	// The handler MUST call Identity().
	//
	// name is a human-readable identifier for the handler. Each handler within
	// an application MUST have a unique name. The name MAY change over time to
	// best reflect the purpose of the handler.
	//
	// name MUST be a non-empty UTF-8 string consisting solely of printable
	// Unicode characters, excluding whitespace. A printable character is any
	// character from the Letter, Mark, Number, Punctuation or Symbol
	// categories.
	//
	// key is an unique identifier for the handler that's used by the engine to
	// correlate its internal state with this handler. For that reason the key
	// SHOULD NOT change once in use.
	//
	// key MUST be an [RFC 4122] UUID expressed as a hyphen-separated, lowercase
	// hexadecimal string, such as "5195fe85-eb3f-4121-84b0-be72cbc5722f".
	//
	// [RFC 4122]: https://www.rfc-editor.org/rfc/rfc4122
	Identity(name string, key string)

	// ConsumesEventType configures the engine to route events of a specific
	// type to the handler.
	//
	// The handler MUST call ConsumesEventType() at least once.
	//
	// The event SHOULD be the zero-value of its type; the engine uses the type
	// information, but not the value itself.
	ConsumesEventType(e Event)

	// ProducesCommandType configures the engine to use the handler as the
	// source of events of a specific type.
	//
	// The handler MUST call ProducesCommandType() at least once.
	//
	// The command SHOULD be the zero-value of its type; the engine uses the
	// type information, but not the value itself.
	ProducesCommandType(c Command)

	// SchedulesTimeoutType configures the engine to allow this handler to
	// schedule timeouts of a specific type.
	//
	// Several handlers MAY schedule timeout messages of the same type.
	//
	// The timeout SHOULD be the zero-value of its type; the engine uses the
	// type information, but not the value itself.
	SchedulesTimeoutType(t Timeout)
}

// ProcessEventScope is an interface implemented by the engine and used by the
// application to perform operations within the context of handling an event
// message.
type ProcessEventScope interface {
	// InstanceID returns the ID of the targeted process instance.
	InstanceID() string

	// End indicates to the engine that the process has ended, and therefore the
	// state of the process root is no longer meaningful.
	//
	// A call to Destroy() is negated by a subsequent call to ExecuteCommand()
	// or ScheduleTimeout() within the same scope.
	//
	// The engine MUST pass a newly initialized process root to the handler when
	// the next event message is handled.
	//
	// The engine MUST discard any timeout messages associated with this
	// instance.
	//
	// The engine MAY allow re-beginning a process instance that has ended.
	// Callers SHOULD assume that such behavior is unavailable.
	End()

	// ExecuteCommand executes a command as a result of the event or timeout
	// message being handled.
	//
	// It MUST NOT be called with a message of any type that has not been
	// configured for production by a prior call to Configure().
	//
	// Any prior call to End() within the same scope is negated.
	ExecuteCommand(c Command)

	// ScheduleTimeout schedules a timeout message to be handled by this process
	// instance at a specific time.
	//
	// Any pending timeout messages are cancelled when the instance is ended.
	//
	// It MUST NOT be called with a message of any type that has not been
	// configured for production by a prior call to Configure().
	//
	// Any prior call to End() within the same scope is negated.
	ScheduleTimeout(t Timeout, at time.Time)

	// RecordedAt returns the time at which the event occurred.
	RecordedAt() time.Time

	// Log records an informational message.
	Log(f string, v ...any)
}

// ProcessTimeoutScope is an interface implemented by the engine and used by the
// application to perform operations within the context of handling a timeout
// message.
type ProcessTimeoutScope interface {
	// InstanceID returns the ID of the targeted process instance.
	InstanceID() string

	// End indicates to the engine that the process has ended, and therefore the
	// state of the process root is no longer meaningful.
	//
	// A call to Destroy() is negated by a subsequent call to ExecuteCommand()
	// or ScheduleTimeout() within the same scope.
	//
	// The engine MUST pass a newly initialized process root to the handler when
	// the next event message is handled.
	//
	// The engine MUST discard any timeout messages associated with this
	// instance.
	//
	// The engine MAY allow re-beginning a process instance that has ended.
	// Callers SHOULD assume that such behavior is unavailable.
	End()

	// ExecuteCommand executes a command as a result of the event or timeout
	// message being handled.
	//
	// It MUST NOT be called with a message of any type that has not been
	// configured for production by a prior call to Configure().
	//
	// Any prior call to End() within the same scope is negated.
	ExecuteCommand(c Command)

	// ScheduleTimeout schedules a timeout message to be handled by this process
	// instance at a specific time.
	//
	// Any pending timeout messages are cancelled when the instance is ended.
	//
	// It MUST NOT be called with a message of any type that has not been
	// configured for production by a prior call to Configure().
	//
	// Any prior call to End() within the same scope is negated.
	ScheduleTimeout(t Timeout, at time.Time)

	// ScheduledFor returns the time at which the timeout was configured to
	// occur.
	ScheduledFor() time.Time

	// Log records an informational message.
	Log(f string, v ...any)
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

// NoTimeoutMessagesBehavior can be embedded in ProcessMessageHandler
// implementations to indicate that no timeout messages are used.
//
// It provides an implementation of ProcessMessageHandler.HandleTimeout() that
// always panics with the UnexpectedMessage value.
type NoTimeoutMessagesBehavior struct{}

// HandleTimeout panics with the UnexpectedMessage value.
func (NoTimeoutMessagesBehavior) HandleTimeout(
	context.Context,
	ProcessRoot,
	ProcessTimeoutScope,
	Timeout,
) error {
	panic(UnexpectedMessage)
}
