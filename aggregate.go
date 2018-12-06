package dogma

// AggregateMessageHandler is an interface implemented by the application and
// used by the engine to cause changes to an aggregate via command messages.
//
// An aggregate is a collection of objects that represent some domain state
// within the application.
//
// A request to change the state of an aggregate is represented by a command
// message. The changes caused by the command, if any, are represented by
// domain event messages.
//
// Each command message targets a single aggregate instance, represented by
// the AggregateRoot interface.
type AggregateMessageHandler interface {
	// New returns a new aggregate instance.
	New() AggregateRoot

	// RouteCommand indicates whether a specific type or instance of a command
	// message should be routed to this handler.
	//
	// If p is false, then m is a command message that has been sent to the
	// application. If m should be routed to this handler the implementation sets
	// ok to true, and id to a non-empty value indicating the ID of the aggregate
	// instance that the command targets.
	//
	// If p is true, then the engine is performing a "routing probe". In this case
	// m is a non-nil, zero-value message. The implementation sets ok to true if
	// messages of this type should be routed to this message handler when they
	// occur. The id output parameter is unused.
	RouteCommand(m Message, p bool) (id string, ok bool)

	// HandleCommand handles a command message that has been routed to this
	// handler.
	//
	// Handling a command involves inspecting the state of the command's target
	// aggregate instance to determine what changes, if any, should occur. Each
	// change is indicated by recording an event message.
	//
	// s provides access to the operations available within the scope of handling
	// m, such as loading the targetted instance and recording event messages.
	//
	// This method must not manipulate the targetted instance directly. Any such
	// manipulations must be applied by the instance's ApplyEvent() method,
	// which is called for each recorded event message.
	//
	// If m was not expected by the handler, as per the routes determined by calls
	// to RouteCommand(), it must panic with an UnexpectedMessage error.
	HandleCommand(s AggregateScope, m Message)
}

// AggregateRoot is an interface implemented by the application, and used by
// the engine to apply changes to an aggregate.
type AggregateRoot interface {
	// ApplyEvent updates the aggregate instance to reflect the fact that a
	// particular event has occurred.
	ApplyEvent(m Message)
}

// AggregateScope is an interface implemented by the engine, and used by the
// application to perform operations within the context of handling a command
// message.
type AggregateScope interface {
	// InstanceID is the ID of the aggregate instance targetted by the command
	// being handled.
	InstanceID() string

	// Instance loads and returns the aggregate instance targetted by the command
	// being handled.
	Instance() AggregateRoot

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
