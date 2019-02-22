package dogma

// AggregateMessageHandler is an interface implemented by the application and
// used by the engine to cause changes to an aggregate instance via command
// messages.
//
// Many instances of each aggregate type can be created. Each instance is a
// collection of objects that represent some domain state within the
// application. All manipulation of an aggregate instance is performed via one
// of its constituent objects, known as the "root", and represented by the
// AggregateRoot interface.
//
// A request to change the state of an aggregate instance is represented by a
// command message. The changes caused by the command message, if any, are
// represented by event messages.
//
// Each command message targets a single aggregate instance of a specific type.
// A command message can cause the creation or destruction of its target
// instance.
type AggregateMessageHandler interface {
	// New constructs a new aggregate instance and returns its root.
	//
	// The return value MUST NOT be nil.
	New() AggregateRoot

	// Configure produces a configuration for this handler by calling methods on
	// the configurer c.
	//
	// The implementation MUST allow for multiple calls to Configure(). Each
	// call SHOULD produce the same configuration.
	//
	// The engine MUST call Configure() before calling HandleEvent(). It is
	// RECOMMENDED that the engine only call Configure() once per handler.
	Configure(c AggregateConfigurer)

	// RouteCommandToInstance returns the ID of the aggregate instance that is
	// targetted by m.
	//
	// The return value MUST be a non-empty string. The use of UUIDs for
	// instance identifiers is RECOMMENDED.
	//
	// The engine MUST NOT call RouteCommandToInstance() with any message of a
	// type that has not been configured for consumption by a prior call to
	// Configure(). If any such message is passed, the implementation MUST panic
	// with the UnexpectedMessage value.
	RouteCommandToInstance(m Message) string

	// HandleCommand handles a command message.
	//
	// Handling a command message involves inspecting the state of the target
	// aggregate instance to determine what changes, if any, should occur. Each
	// change is indicated by recording an event message.
	//
	// The targeted instance MUST NOT be modified directly. All modifications
	// must be applied by the instance's ApplyEvent() method, which is called
	// for each event message that is recorded via s.
	//
	// The engine SHOULD provide "at-least-once" delivery guarantees to the
	// handler. That is, the engine should call HandleCommand() with the same
	// command message until no panic occurs.
	//
	// The engine MUST NOT call HandleCommand() with any message of a type that
	// has not been configured for consumption by a prior call to Configure().
	// If any such message is passed, the implementation MUST panic with the
	// UnexpectedMessage value.
	//
	// The implementation MUST NOT assume that HandleCommand() will be called
	// with commands in the same order that they were executed.
	//
	// The engine MAY call HandleCommand() from multiple goroutines concurrently.
	HandleCommand(s AggregateCommandScope, m Message)
}

// AggregateRoot is an interface implemented by the application and used by
// the engine to apply changes to an aggregate instance.
type AggregateRoot interface {
	// ApplyEvent updates the aggregate instance to reflect the occurence of an
	// event that was recorded against this instance.
	//
	// It MUST NOT be called with a message of any type that has not been
	// configured for production by a prior call to Configure().
	//
	// It MUST accept all messages ofthe types  that have been configured for
	// production, though any given call MAY be a no-op.
	ApplyEvent(m Message)
}

// AggregateConfigurer is an interface implemented by the engine and used by
// the application to configure options related to an AggregateMessageHandler.
//
// It is passed to AggregateMessageHandler.Configure(), typically upon
// initialization of the engine.
//
// In the context of this interface, "the handler" refers to the handler on
// which Configure() has been called.
type AggregateConfigurer interface {
	// Name sets the name of the handler.
	//
	// It MUST be called exactly once within a single call to Configure().
	//
	// Each handler within an application MUST have a unique, non-empty name.
	Name(n string)

	// ConsumesCommandType configures the engine to route command messages of
	// the same type as m to the handler.
	//
	// It MUST be called at least once within a call to Configure(). It MUST NOT
	// be called more than once with a command message of the same type.
	//
	// A given command type MUST be routed to exactly one handler within an
	// application.
	//
	// The "content" of m MUST NOT be used, inspected, or treated as meaningful
	// in any way, only its runtime type information.
	ConsumesCommandType(m Message)

	// ProducesEventType instructs the engine that the handler records events of
	// the same type as m.
	//
	// It MUST be called at least once within a call to Configure(). It MUST NOT
	// be called more than once with an event message of the same type.
	//
	// A given event type MUST be produced by exactly one handler within an
	// application.
	//
	// The "content" of m MUST NOT be used, inspected, or treated as meaningful
	// in any way, only its runtime type information.
	ProducesEventType(m Message)
}

// AggregateCommandScope is an interface implemented by the engine and used by the
// application to perform operations within the context of handling a specific
// domain command message.
type AggregateCommandScope interface {
	// InstanceID returns the ID of the targeted aggregate instance.
	InstanceID() string

	// Create creates the targeted instance.
	//
	// It MUST be called before Root() or RecordEvent() can be called within
	// this scope or the scope of any future message that targets the same
	// instance.
	//
	// It returns true if the targeted instance was created, or false if
	// the instance already exists.
	//
	// If it returns true, RecordEvent() MUST be called at least once within the
	// same scope. This guarantees that the creation of every instance is
	// represented by an application-defined event.
	Create() bool

	// Destroy destroys the targeted instance.
	//
	// After it has been called, Root() and RecordEvent() MUST NOT be called
	// within this scope or the scope of any future message that targets the
	// same instance, unless Create() is called again first.
	//
	// It MUST NOT be called if the instance does not currently exist.
	//
	// RecordEvent() MUST be called at least once within the same scope. This
	// guarantees that the destruction of every instance is represented by an
	// application-defined event.
	//
	// The precise semantics of destroy are implementation defined. The
	// aggregate data MAY be deleted or archived, for example.
	Destroy()

	// Root returns the root of the targeted aggregate instance.
	//
	// It MUST NOT be called if the instance does not currently exist.
	Root() AggregateRoot

	// RecordEvent records the occurrence of an event as a result of the command
	// message that is being handled.
	//
	// It MUST NOT be called with a message of any type that has not been
	// configured for production by a prior call to Configure().
	//
	// It MUST NOT be called if the instance does not currently exist.
	//
	// The engine MUST call Root().ApplyEvent(m) before returning, such that the
	// applied changes are visible to the handler.
	RecordEvent(m Message)

	// Log records an informational message within the context of the command
	// message that is being handled.
	Log(f string, v ...interface{})
}

// StatelessAggregateBehavior can be embedded in AggregateMessageHandler
// implementations to indicate that no state is kept between messages.
//
// It provides an implementation of AggregateMessageHandler.New() that always
// returns a StatelessAggregateRoot.
type StatelessAggregateBehavior struct{}

// New returns StatelessAggregateRoot.
func (StatelessAggregateBehavior) New() AggregateRoot {
	return StatelessAggregateRoot
}

// StatelessAggregateRoot is an aggregate root with no state.
//
// It can be returned by an AggregateMessageHandler.New() implementation to
// indicate that no domain state is required beyond the existence/non-existence
// of the aggregate instance.
//
// See also StatelessAggregateBehavior, which provides an implementation of
// New() that returns this value.
//
// Engines may use this value as a sentinel, to provide an optimized code path
// when no state is required.
var StatelessAggregateRoot AggregateRoot = statelessAggregateRoot{}

type statelessAggregateRoot struct{}

func (statelessAggregateRoot) ApplyEvent(m Message) {
	if m == nil {
		panic("event must not be nil")
	}
}
