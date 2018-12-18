package dogma

// A Message is an application-defined unit of data.
//
// Most messages fall into one of two broad categories; "domain messages", which
// relate to the application's domain logic, and "integration messages" which
// are used to integrate the domain with "non-domain concerns", such as
// third-party APIs. There are other kinds of messages, such as process timeout
// messages which do not clearly belong to one category or another.
//
// Within these categories, messages are further divided into "commands" and
// "events". Command messages represent a request to perform some action,
// whereas event messages represent some occurrance which has already taken
// place.
//
// These categorizations are largely conceptual. Within Dogma, they are all
// modeled by the Message interface, though engine implementations may require
// messages to implement more specific interfaces for each category.
type Message interface {
}

// UnexpectedMessage is a panic value used by a message handler when it receives
// a message that should not have been routed to it.
var UnexpectedMessage unexpectedMessage

type unexpectedMessage struct{}
