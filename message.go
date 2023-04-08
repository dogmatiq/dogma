package dogma

import "errors"

// A Message is an application-defined unit of data that describes a [Command],
// [Event], or [Timeout] within a message-based application.
type Message interface {
	// MessageDescription returns a human-readable description of the message.
	MessageDescription() string

	// Validate returns a non-nil error if the message is invalid.
	Validate() error
}

// A Command is a message that represents a request for a Dogma application to
// perform some action.
type Command = Message

// An Event is a message that indicates that some action has occurred within a
// Dogma application.
type Event = Message

// A Timeout is a message that represents a request for an action to be
// performed at a specific time.
type Timeout = Message

// UnexpectedMessage is a panic value used by a message handler when it receives
// a message of a type that it did not expect.
var UnexpectedMessage unexpectedMessage

type unexpectedMessage struct{}

// ValidateMessage returns an error if m is invalid.
func ValidateMessage(m Message) error {
	if m == nil {
		return errors.New("message must not be nil")
	}
	return m.Validate()
}
