package dogma

import (
	"time"
)

// NoTimeoutHintBehavior is an embeddable type for [ProcessMessageHandler],
// [IntegrationMessageHandler] and [ProjectionMessageHandler] implementations
// that do not provide a message handling timeout hint.
type NoTimeoutHintBehavior struct{}

// TimeoutHint always returns zero.
func (NoTimeoutHintBehavior) TimeoutHint(Message) time.Duration {
	return 0
}
