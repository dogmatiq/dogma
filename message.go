package dogma

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
// Engine implementations MAY place further requirements upon message
// implementations.
type Message interface {
}

// UnexpectedMessage is a panic value used by a message handler when it receives
// a message of a type that it did not expect.
var UnexpectedMessage unexpectedMessage

type unexpectedMessage struct{}
