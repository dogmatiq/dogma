package dogma

// A HandlerConfigurer is the interface a handler uses to declare its
// configuration. It defines configuration methods common to all handler types.
//
// Each handler type has a corresponding configurer type that extends this
// interface. See [AggregateConfigurer], [ProcessConfigurer],
// [IntegrationConfigurer], and [ProjectionConfigurer].
type HandlerConfigurer interface {
	// Identity sets the handler's unique identity.
	//
	// n is a short human-readable name displayed wherever the handler's
	// identity appears, such as in logs and telemetry signals. The value must
	// be between 1 and 255 bytes in length, and contain only printable,
	// non-space UTF-8 characters. Changing the handler's name does not affect
	// its behavior. Each handler within an application must have a unique name.
	//
	// k is a key that uniquely identifies the handler. The engine uses the key
	// to associate handler state with the correct handler instance, so it must
	// not change. The value must be a canonical RFC 4122 UUID string, such as
	// "3a6da031-aa6c-406a-8453-73762f71f917", and is case-insensitive.
	Identity(n string, k string)

	// Disable excludes the handler from the application's runtime without
	// removing it from the application's configuration.
	//
	// Disabling a handler is appropriate when it cannot function in the current
	// execution environment but should still be present for documentation,
	// static analysis, and route discovery purposes.
	//
	// Prefer conditionally disabling a handler over conditionally adding it to
	// the application.
	Disable(...DisableOption)
}

// DisableOption is an option that modifies the behavior of the
// [HandlerConfigurer].Disable method.
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
