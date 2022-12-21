package dogma

// Application is an interface implemented by the application and used by the
// engine to describe the structure of an application.
type Application interface {
	// Configure configures the behavior of the engine as it relates to this
	// application.
	//
	// c provides access to the various configuration options, such as specifying
	// which message handlers the application contains.
	//
	// The implementation MUST allow for multiple calls to Configure(). Each
	// call SHOULD produce the same configuration.
	Configure(c ApplicationConfigurer)
}

// ApplicationConfigurer is an interface implemented by the engine and used by
// the application to configure options related to the application itself.
//
// It is passed to Application.Configure(), typically upon initialization of the
// engine.
type ApplicationConfigurer interface {
	// Identity sets unique identifiers for the application.
	//
	// It MUST be called exactly once within a single call to Configure().
	//
	// The name is a human-readable identifier for the application. The
	// application name SHOULD be distinct from that of any handlers within the
	// application. The name MAY be changed over time to best reflect the
	// purpose of the application.
	//
	// The key is an immutable identifier for the application. Its purpose is to
	// allow engine implementations to uniquely identify the application in a
	// multi-application environment, and to associate ancillary data with the
	// application, such as message routing information.
	//
	// The application and the handlers within it MUST have distinct keys. The
	// key MUST NOT be changed. The RECOMMENDED key format is an RFC 4122 UUID
	// represented as a hyphen-separated, lowercase hexadecimal string, such as
	// "5195fe85-eb3f-4121-84b0-be72cbc5722f".
	//
	// Both identifiers MUST be non-empty UTF-8 strings consisting solely of
	// printable Unicode characters, excluding whitespace. A printable character
	// is any character from the Letter, Mark, Number, Punctuation or Symbol
	// categories.
	//
	// The engine MUST NOT perform any case-folding or normalization of
	// identifiers. Therefore, two identifiers compare as equivalent if and only
	// if they consist of the same sequence of bytes.
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
