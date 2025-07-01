package dogma

// DisableOption is an option that modifies the behavior of the Disable() method
// on a handler configurer.
//
// See:
//   - [AggregateConfigurer].Disable()
//   - [ProcessConfigurer].Disable()
//   - [IntegrationConfigurer].Disable()
//   - [ProjectionConfigurer].Disable()
//
// This type exists for forward-compatibility.
type DisableOption interface {
	futureDisableOption()
}
