package dogma

import (
	"errors"
	"fmt"
)

// A Message is an application-defined unit of data that encapsulates a
// "command" or "event" within a message-based application.
//
// The message implementations are provided by the application.
//
// Message implementations SHOULD implement fmt.Stringer or DescribableMessage
// in order to provide a human-readable description of every message.
//
// Message implementations SHOULD implement ValidatableMessage in order to
// allow the engine to validate messages before they enter the application.
//
// Engine implementations MAY place further requirements upon message
// implementations.
type XMessage interface {
}

// A Command is a message that represents a request for a Dogma application to
// perform some action.
type Command = XMessage

// An Event is a message that indicates that some action has occurred within a
// Dogma application.
type Event = XMessage

// A Timeout is a message that encapsulates information about an action that was
// scheduled to occur at a specific time.
type Timeout = XMessage

// DescribableMessage is a message that can provide a human-readable description
// of itself.
//
// This interface can be implemented to provide a more specific message
// description for message types that already implement fmt.Stringer in such a
// way that does not provide a useful human-readable description, such as when
// the message implementations are generated Protocol Buffers structs.
type DescribableMessage interface {
	XMessage

	// MessageDescription returns a human-readable description of the message.
	//
	// This method SHOULD NOT be called directly. Instead, obtain the
	// description using the DescribeMessage() function.
	MessageDescription() string
}

// DescribeMessage returns a human-readable string representation of m.
//
// If m implements DescribableMessage, it returns m.MessageDescription().
// Otherwise, if m implements fmt.Stringer, it returns m.String().
//
// Finally, if m does not implement either of these interfaces, it returns the
// standard Go "%v" representation of the message.
//
// Engine implementations SHOULD use the message description in logging and
// other tracing systems to provide contextual information to developers. The
// description SHOULD NOT be used by application code.
func DescribeMessage(m XMessage) string {
	switch m := m.(type) {
	case DescribableMessage:
		return m.MessageDescription()
	case fmt.Stringer:
		return m.String()
	default:
		return fmt.Sprintf("%v", m)
	}
}

// ValidatableMessage is a message that can validate itself.
//
// This interface can be implemented to perform fine-grained validation of
// messages.
//
// Engine implementations SHOULD validate messages before allowing them to be
// produced in order to prevent "poison" messages from entering the application.
type ValidatableMessage interface {
	XMessage

	// Validate returns a non-nil error if the message is invalid.
	Validate() error
}

// ValidateMessage returns an error if m implements ValidatableMessage and is
// invalid.
//
// If m does not implement ValidatableMessage it returns nil.
func ValidateMessage(m XMessage) error {
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
