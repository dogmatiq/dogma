package dogma

import "fmt"

// A Message is an application-defined unit of data that encapsulates a
// "command" or "event" within a message-based application.
//
// Command messages represent a request for the application to perform some
// action, whereas event messages indicate that some action has already
// occurred.
//
// Additionally, a "timeout" message can be used to perform actions within an
// application at specific wall-clock times.
//
// The message implementations are provided by the application. The interface is
// intentionally empty, allowing the use of any Go type as a message.
//
// Message implementations SHOULD implement fmt.Stringer or DescribableMessage
// in order to provide a human-readable description of every message.
//
// Engine implementations MAY place further requirements upon message
// implementations.
type Message interface {
}

// DescribableMessage is a message that can provide a human-readable description
// of itself.
//
// This interface can be implemented to provide a more specific message
// description for message types that already implement fmt.Stringer in such a
// way that does not provide a useful human-readable description, such as when
// the message implementations are generated Protocol Buffers structs.
type DescribableMessage interface {
	Message

	// MessageDescription returns a human-readable description of the message.
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
// other tracing systems to provide contextual information to developers.
func DescribeMessage(m Message) string {
	if s, ok := m.(DescribableMessage); ok {
		return s.MessageDescription()
	}

	if s, ok := m.(fmt.Stringer); ok {
		return s.String()
	}

	return fmt.Sprintf("%v", m)
}

// UnexpectedMessage is a panic value used by a message handler when it receives
// a message of a type that it did not expect.
var UnexpectedMessage unexpectedMessage

type unexpectedMessage struct{}
