package dogma

// AggregateMessageHandler is an interface implemented by the application and
// used by the engine to cause changes to an aggregate via command messages.
//
// Many instances of each aggregate type can be created. Each instance is a
// collection of objects that represent some domain state within the
// application. All manipulation of an aggregate instance is performed via one
// of its constituent objects, known as the "root", and represented by the
// AggregateRoot interface.
//
// A request to change the state of an aggregate instance is represented by a
// command message. The changes caused by the command, if any, are represented
// by domain event messages.
//
// Each command message targets a single aggregate instance of a specific type.
// A command can cause the creation or destruction of its target instance.
type AggregateMessageHandler interface {
	// New constructs a new aggregate instance and returns its root.
	New() AggregateRoot

	// Configure configures the behavior of the engine as it relates to this
	// handler.
	//
	// c provides access to the various configuration options, such as
	// specifying which command types are routed to this handler.
	Configure(c AggregateConfigurer)

	// RouteCommandToInstance returns the ID of the aggregate instance that is
	// targetted by m.
	//
	// It panics with the UnexpectedMessage value if m is not one of the command
	// types that is routed to this handler via Configure().
	RouteCommandToInstance(m Message) string

	// HandleCommand handles a command message that has been routed to this
	// handler.
	//
	// Handling a command involves inspecting the state of the command's target
	// aggregate instance to determine what changes, if any, should occur. Each
	// change is indicated by recording an event message.
	//
	// s provides access to the operations available within the scope of handling
	// m, such as creating or destroying the targeted instance, accessing its
	// state, and recording event messages.
	//
	// This method must not modify the targeted instance directly. All
	// modifications must be applied by the instance's ApplyEvent() method, which
	// is called for each recorded event message.
	//
	// It panics with the UnexpectedMessage value if m is not one of the command
	// types that is routed to this handler via Configure().
	HandleCommand(s AggregateScope, m Message)
}

// AggregateRoot is an interface implemented by the application and used by
// the engine to apply changes to an aggregate instance.
type AggregateRoot interface {
	// ApplyEvent updates the aggregate instance to reflect the fact that a
	// particular event has occurred.
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
	// RouteCommandType configures the engine to route commands of the same type as m
	// to the handler.
	RouteCommandType(m Message)
}

// AggregateScope is an interface implemented by the engine and used by the
// application to perform operations within the context of handling a command
// message.
//
// In the context of this interface, "the message" refers to the message being
// handled and "the instance" refers to the aggregate instance that is targeted
// by that message.
type AggregateScope interface {
	// InstanceID is the ID of the targeted aggregate instance.
	InstanceID() string

	// Create creates the targeted instance.
	//
	// It must be called before Root() or RecordEvent() can be called within this
	// scope or the scope of any future command that targets the same instance.
	//
	// It returns true if the targeted instance was created, or false if
	// the instance already exists.
	//
	// If it returns true, RecordEvent() must be called at least once within
	// the same scope. This guarantees that the creation of every instance is
	// represented by a domain event.
	Create() bool

	// Destroy destroys the targeted instance.
	//
	// After it has been called neither Root() nor RecordEvent() can be called
	// within this scope or the scope of any future command that targets the same
	// instance, unless Create() is called again first.
	//
	// It panics if the target instance does not currently exist.
	//
	// RecordEvent() must be called at least once within the same scope. This
	// guarantees that the destruction of every instance is represented by a domain
	// event.
	//
	// The precise semantics of destroy are implementation defined. The aggregate
	// data may be deleted or archived, for example.
	Destroy()

	// Root returns the root of the targeted aggregate instance.
	//
	// It panics if the instance has not been created, or was created but has
	// subsequently been destroyed.
	Root() AggregateRoot

	// RecordEvent records the occurrence of an event as a result of the command
	// being handled.
	//
	// It panics if the instance has not been created, or was created but has
	// subsequently been destroyed.
	//
	// The engine must call Instance().ApplyEvent(m) before returning, such that
	// the applied changes are visible to the handler.
	RecordEvent(m Message)

	// Log records an informational message within the context of the command being
	// handled.
	//
	// The log message should be worded such that it makes sense to anyone familiar
	// with the business domain.
	Log(f string, v ...interface{})
}
