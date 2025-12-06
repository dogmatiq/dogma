package dogma

const (
	// MaximizeConcurrency is instructs the engine to process as many messages
	// concurrently as possible, maximizing throughput.
	MaximizeConcurrency concurrencyPreference = iota

	// MinimizeConcurrency is instructs the engine to attempt to process
	// messages one at a time, minimizing conflicts and/or contention.
	MinimizeConcurrency
)

// concurrencyPreference is the underlying type for
// [IntegrationConcurrencyPreference] and [ProjectionConcurrencyPreference].
type concurrencyPreference int
