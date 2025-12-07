package dogma

const (
	// MaximizeConcurrency is instructs the engine to process as many messages
	// concurrently as possible, maximizing throughput.
	MaximizeConcurrency ConcurrencyPreference = iota

	// MinimizeConcurrency is instructs the engine to attempt to process
	// messages one at a time, minimizing conflicts and/or contention.
	MinimizeConcurrency
)

// ConcurrencyPreference is a hint to the engine as to the best way to handle
// concurrent messages for a message handler.
//
// [IntegrationMessageHandler] and [ProjectionMessageHandler] and support
// configuring their concurrency preference via their respective configurer
// interfaces.
type ConcurrencyPreference int
