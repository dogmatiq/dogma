package dogma

// Application is an interface implemented by the application and used by the
// engine to describe the structure of an application.
type Application interface {
	// Configure describes the application's configuration to the engine.
	Configure(c ApplicationConfigurer)
}

// An ApplicationConfigurer configures the engine for use with a specific
// application.
//
// See [ApplicationMessageHandler.Configure]().
type ApplicationConfigurer interface {
	// Identity configures how the engine identifies the application.
	//
	// The application MUST call Identity().
	//
	// name is a human-readable identifier for the application. The name MAY
	// change over time to best reflect the purpose of the application.
	//
	// name MUST be a non-empty UTF-8 string consisting solely of printable
	// Unicode characters, excluding whitespace. A printable character is any
	// character from the Letter, Mark, Number, Punctuation or Symbol
	// categories.
	//
	// key is an unique identifier for the application that's used by the engine
	// to correlate its internal state with this application. For that reason
	// the key SHOULD NOT change once in use.
	//
	// key MUST be an [RFC 4122] UUID expressed as a hyphen-separated, lowercase
	// hexadecimal string, such as "5195fe85-eb3f-4121-84b0-be72cbc5722f".
	//
	// [RFC 4122]: https://www.rfc-editor.org/rfc/rfc4122
	Identity(name string, key string)

	// RegisterAggregate configures the engine to route messages to h.
	RegisterAggregate(h AggregateMessageHandler)

	// RegisterProcess configures the engine to route messages to h.
	RegisterProcess(h ProcessMessageHandler)

	// RegisterIntegration configures the engine to route messages to h.
	RegisterIntegration(h IntegrationMessageHandler)

	// RegisterProjection configures the engine to route messages to h.
	RegisterProjection(h ProjectionMessageHandler)
}
