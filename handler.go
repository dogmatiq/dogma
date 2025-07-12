package dogma

import "time"

// A HandlerConfigurer is the interface a handler uses to declare its
// configuration. It defines configuration methods common to all handler types.
//
// Each handler type has a corresponding configurer type that extends this
// interface.
//
// See:
//   - [AggregateConfigurer]
//   - [ProcessConfigurer]
//   - [IntegrationConfigurer]
//   - [ProjectionConfigurer].
type HandlerConfigurer interface {
	// Identity sets the handler's unique identity.
	//
	// n is a short human-readable name displayed wherever the handler's
	// identity appears, such as in logs and telemetry signals. The value must
	// be between 1 and 255 bytes in length, and contain only printable,
	// non-space UTF-8 characters. Changing the handler's name doesn't affect
	// its behavior. Each handler within an application must have a unique name.
	//
	// k is a key that uniquely identifies the handler. The engine uses the key
	// to associate handler state with the correct handler instance - it must
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

// HandlerScope represents the context within which the engine invokes a message
// handler.
//
// Each operation that a handler performs has a corresponding scope type that
// extends this interface:
//
//   - [AggregateCommandScope]
//   - [ProcessEventScope]
//   - [ProcessTimeoutScope]
//   - [IntegrationCommandScope]
//   - [ProjectionEventScope]
//   - [ProjectionCompactScope]
type HandlerScope interface {
	// Now returns the current local time according to the engine.
	//
	// Use this method in preference to [time.Now]; its return value may differ
	// in certain situations, such as when executing tests or when compensating
	// for clock skew in a distributed system.
	Now() time.Time

	// Log records an informational message using [fmt.Printf]-style formatting.
	//
	// The message should be clear and relevant to developers and non-technical
	// stakeholders familiar with the application's domain. It's not intended
	// for display to end users.
	//
	// Use lowercase sentences with no trailing punctuation and omit sensitive
	// information.
	//
	// Use this method to explain conditions or decisions that aren't captured
	// in a [Message]. For example, if a handler receives a command to cancel a
	// shopping cart order after shipping, it might log “cannot cancel order
	// #49412, it has already shipped”.
	Log(format string, args ...any)
}
