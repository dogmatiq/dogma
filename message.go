package dogma

// Message is an application-defined unit of data.
type Message interface {
}

// UnexpectedMessage is a panic value used by a message handler when it receives
// a message that should not have been routed to it.
var UnexpectedMessage unexpectedMessage

type unexpectedMessage struct{}
