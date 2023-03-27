package dogma

import "time"

// NoTimeoutHintBehavior can be embedded in message handler implementations
// to indicate the handler is unable to suggest a message handling timeout duration.
//
// It provides an implementation of TimeoutHint() method that always returns
// a zero value.
//
// The TimeoutHint() method is present in the ProcessMessageHandler,
// IntegrationMessageHandler and ProjectionMessageHandler interfaces.
type NoTimeoutHintBehavior struct{}

// TimeoutHint always returns a zero-value duration.
func (NoTimeoutHintBehavior) TimeoutHint(XMessage) time.Duration {
	return 0
}
