package dogma

// A Message is an application-defined unit of data that describes a [Command],
// [Event], or [Timeout] within a message-based application.
type Message interface {
	// MessageDescription returns a human-readable description of the message.
	MessageDescription() string
}

// A Command is a message that represents a request for a Dogma application to
// perform some action.
type Command interface {
	// MessageDescription returns a human-readable description of the message.
	MessageDescription() string

	// Validate returns a non-nil error if the message is invalid.
	Validate(CommandValidationScope) error
}

// An Event is a message that indicates that some action has occurred within a
// Dogma application.
type Event interface {
	// MessageDescription returns a human-readable description of the message.
	MessageDescription() string

	// Validate returns a non-nil error if the message is invalid.
	Validate(EventValidationScope) error
}

// A Timeout is a message that represents a request for an action to be
// performed at a specific time.
type Timeout interface {
	// MessageDescription returns a human-readable description of the message.
	MessageDescription() string

	// Validate returns a non-nil error if the message is invalid.
	Validate(TimeoutValidationScope) error
}

// UnexpectedMessage is a panic value used by a message handler when it receives
// a message of a type that it did not expect.
var UnexpectedMessage unexpectedMessage

type unexpectedMessage struct{}

// CommandValidationScope provides information about the context in which a
// [Command] is being validated.
type CommandValidationScope interface {
	future()
}

// EventValidationScope provides information about the context in which an
// [Event] is being validated.
type EventValidationScope interface {
	future()
}

// TimeoutValidationScope provides information about the context in which a
// [Timeout] is being validated.
type TimeoutValidationScope interface {
	future()
}
