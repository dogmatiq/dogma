package dogma

// AggregateMessageHandler is an interface implemented by the application and
// used by the engine to cause changes to an aggregate via command messages.
//
// An aggregate is a collection of objects that represent some domain state
// within the application. All manipulation of the aggregate is performed via
// one of its constituent object, known as the "root", and represented by the
// AggregateRoot interface.
//
// A request to change the state of an aggregate is represented by a command
// message. The changes caused by the command, if any, are represented by domain
// event messages. Each command message targets a single aggregate instance.
type AggregateMessageHandler interface {
	// New constructs a new aggregate instance, and returns its root.
	New() AggregateRoot

	// RouteCommand indicates whether a specific type or instance of a message
	// should be routed to this handler as a command.
	//
	// If p is false, then m is a command message that has been sent to the
	// application. If m should be routed to this handler, the implementation sets
	// ok to true and id to the ID of the aggregate instance that the command
	// targets. id must not be empty if ok is true.
	//
	// If p is true, then the engine is performing a "routing probe". In this case
	// m is a non-nil, zero-value message. The implementation sets ok to true if
	// messages of the same type as m should be routed to this message handler when
	// they occur. The id output parameter is unused.
	RouteCommand(m Message, p bool) (id string, ok bool)

	// HandleCommand handles a command message that has been routed to this
	// handler.
	//
	// Handling a command involves inspecting the state of the command's target
	// aggregate instance to determine what changes, if any, should occur. Each
	// change is indicated by recording an event message.
	//
	// s provides access to the operations available within the scope of handling
	// m, such as accessing the targetted instance and recording event messages.
	//
	// This method must not modify the targetted instance directly. All
	// modifications must be applied by the instance's ApplyEvent() method, which
	// is called for each recorded event message.
	//
	// If m was not expected by the handler the implementation must panic with an
	// UnexpectedMessage value.
	HandleCommand(s AggregateScope, m Message)
}

// AggregateRoot is an interface implemented by the application and used by
// the engine to apply changes to an aggregate instance.
type AggregateRoot interface {
	// ApplyEvent updates the aggregate instance to reflect the fact that a
	// particular event has occurred.
	ApplyEvent(m Message)
}

// AggregateScope is an interface implemented by the engine and used by the
// application to perform operations within the context of handling a command
// message.
//
// In the context of this interface, "the message" refers to the message being
// handled, and "the instance" refers to the aggregate instance that is
// targetted by that message.
type AggregateScope interface {
	// InstanceID is the ID of the targetted aggregate instance.
	InstanceID() string

	// Root returns the root of the targetted aggregate instance.
	Root() AggregateRoot

	// RecordEvent records the occurrence of an event as a result of the command
	// being handled.
	//
	// The engine must call Instance().ApplyEvent(m) before returning, such that
	// the applied changes are visible to the handler.
	RecordEvent(m Message)

	// Log logs an informational message within the context of the command being
	// handled.
	//
	// The log message should be worded such that it makes sense to anyone familiar
	// with the business domain.
	Log(f string, v ...interface{})
}
