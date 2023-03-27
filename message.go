package dogma

import (
	"errors"
)

// A Message is an application-defined unit of data that encapsulates a
// "command" or "event" within a message-based application.
//
// The message implementations are provided by the application.
//
// Message implementations SHOULD implement ValidatableMessage in order to allow
// the engine to validate messages before they enter the application.
//
// Engine implementations MAY place further requirements upon message
// implementations.
type Message interface {
	// MessageDescription returns a human-readable description of the message.
	MessageDescription() string
}

// A Command is a message that represents a request for a Dogma application to
// perform some action.
type Command = Message

// An Event is a message that indicates that some action has occurred within a
// Dogma application.
type Event = Message

// A Timeout is a message that encapsulates information about an action that was
// scheduled to occur at a specific time.
type Timeout = Message

// DescribableMessage is a message that can provide a human-readable description
// of itself.
//
// Deprecated: All messages are now required to be describable.
type DescribableMessage interface {
	Message
}

// DescribeMessage returns a human-readable string representation of m.
//
// Deprecated: All messages are now required to be describable. Call
// c.MessageDescription() directly.
func DescribeMessage(m Message) string {
	return m.MessageDescription()
}

// ValidatableMessage is a message that can validate itself.
//
// This interface can be implemented to perform fine-grained validation of
// messages.
//
// Engine implementations SHOULD validate messages before allowing them to be
// produced in order to prevent "poison" messages from entering the application.
type ValidatableMessage interface {
	Message

	// Validate returns a non-nil error if the message is invalid.
	Validate() error
}

// ValidateMessage returns an error if m implements ValidatableMessage and is
// invalid.
//
// If m does not implement ValidatableMessage it returns nil.
func ValidateMessage(m Message) error {
	switch m := m.(type) {
	case ValidatableMessage:
		return m.Validate()
	case nil:
		return errors.New("message must not be nil")
	default:
		return nil
	}
}

// UnexpectedMessage is a panic value used by a message handler when it receives
// a message of a type that it did not expect.
var UnexpectedMessage unexpectedMessage

type unexpectedMessage struct{}
