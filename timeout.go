package dogma

import (
	"time"
)

// NoTimeoutHintBehavior can be embedded in message handler implementations to
// denote that the handler is unable to suggest a message handling timeout
// duration.
//
// It provides an implementation of the TimeoutHint() method that always returns
// zero.
//
// The TimeoutHint() method is present in the [ProcessMessageHandler],
// [IntegrationMessageHandler] and [ProjectionMessageHandler] interfaces.
type NoTimeoutHintBehavior struct{}

// TimeoutHint always returns zero.
func (NoTimeoutHintBehavior) TimeoutHint(Message) time.Duration {
	return 0
}
