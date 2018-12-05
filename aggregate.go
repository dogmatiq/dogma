package dogma

// AggregateMessageHandler is an interface for applying changes to an aggregate
// via command messages.
//
// An aggregate is a collection of objects that represent some domain state
// within the application.
//
// A request to change the state of an aggregate is represented by a command
// message. The changes the command induces are represented by domain event
// messages.
//
// Each command message is routed to a single aggregate instance, represented by
// the AggregateRoot interface.
type AggregateMessageHandler interface {
	// New returns a new aggregate instance.
	New() AggregateRoot

	// RouteCommand determines how to route a command message to this handler,
	// based on its type.
	//
	// If ok is true, m is routed to this handler, within the context of the
	// aggregate instance nominated by id. id must not be empty.
	//
	// If p is true, a "routing probe" is being performed, in which case m is a
	// zero-value and id is ignored. ok should be true if this handler should
	// receive messages of m's type.
	RouteCommand(m Message, p bool) (id string, ok bool)

	// HandleCommand handles a command message within the scope of a specific
	// aggregate instance.
	//
	// It inspects the state of the instance returned by s.Instance() to determine
	// what state changes, if any, the command should produce.
	//
	// Any such change is indicated by calling s.RecordEvent() with an event
	// message that represents that state change. The instance MUST NOT be modified
	// directly by this method. Instead, the modifications must be made by the
	// instance's ApplyEvent() method, which is called for each recorded event.
	//
	// If m is of an unsupported message type, it must panic with an
	// UnexpectedMessage error.
	HandleCommand(s AggregateScope, m Message)
}

// AggregateRoot is an interface to an aggregate instance.
type AggregateRoot interface {
	// ApplyEvent updates the aggregate instance to reflect the fact that a
	// particular event has occurred.
	ApplyEvent(m Message)
}

// AggregateScope is an interface used to access and manipulate an aggregate
// when handling a command message.
type AggregateScope interface {
	// InstanceID is the ID of the aggregate instance that the command message has
	// been routed to.
	InstanceID() string

	// Instance loads and returns the aggregate instance that the command message
	// has been routed to.
	Instance() AggregateRoot

	// RecordEvent records the occurrence of an event.
	RecordEvent(m Message)

	// Log logs an informational message within the context of this command.
	Log(f string, v ...interface{})
}
