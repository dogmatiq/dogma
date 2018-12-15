package dogma

// A Message is an application-defined unit of data.
//
// Messages are divided into two broad categories; "domain messages", which
// relate to the application's domain logic, and "integration messages" which
// are used to integrate the domain with "non-domain concerns", such as
// third-party APIs.
//
// Within each category, messages are further divided into "commands" and
// "events". Command messages represent a request to perform some action,
// whereas event messages represent some occurrance which has already taken
// place.
type Message interface {
}

// UnexpectedMessage is a panic value used by a message handler when it receives
// a message that should not have been routed to it.
var UnexpectedMessage unexpectedMessage

type unexpectedMessage struct{}
