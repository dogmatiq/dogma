package dogma

import "time"

// A Message is an application-defined unit of data that describes a [Command],
// [Event], or [Timeout] within a message-based application.
//
// Typically, application developers should use one of the more specific
// [Command], [Event], or [Timeout] interface instead.
type Message interface {
	// MessageDescription returns a human-readable description of the message.
	MessageDescription() string
}

// A Command is a message that represents a request for a Dogma application to
// perform some action.
type Command interface {
	Message

	// Validate returns a non-nil error if the message is invalid.
	Validate(CommandValidationScope) error
}

// An Event is a message that indicates that some action has occurred within a
// Dogma application.
type Event interface {
	Message

	// Validate returns a non-nil error if the message is invalid.
	Validate(EventValidationScope) error
}

// A Timeout is a message that represents a request for an action to be
// performed at a specific time.
type Timeout interface {
	Message

	// Validate returns a non-nil error if the message is invalid.
	Validate(TimeoutValidationScope) error
}

// UnexpectedMessage is a panic value used by a message handler when it receives
// a message of a type that it did not expect.
var UnexpectedMessage unexpectedMessage

type unexpectedMessage struct{}

// CommandValidationScope provides information about the scope under which a
// [Command] is being validated.
type CommandValidationScope interface {
	reserved()
}

// EventValidationScope provides information about the scope under which an
// [Event] is being validated.
type EventValidationScope interface {
	// RecordedAt returns the time at which the event being validated was
	// recorded.
	//
	// If the event has not yet been recorded, ok is false and the value of t is
	// undefined.
	RecordedAt() (t time.Time, ok bool)
}

// TimeoutValidationScope provides information about the scope under which a
// [Timeout] is being validated.
type TimeoutValidationScope interface {
	reserved()
}
