package dogma

// Application is an interface implemented by the application and used by the
// engine to describe the structure of an application.
type Application interface {
	// Configure configures the behavior of the engine as it relates to this
	// application.
	//
	// c provides access to the various configuration options, such as specifying
	// which message handlers the application contains.
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
	// key MUST NOT be changed. The RECOMMENDED key format is RFC 4122 UUID,
	// generated when the handler is first implemented.
	//
	// Both the name and the key MUST be non-empty UTF-8 strings consisting
	// solely of printable Unicode characters, excluding whitespace. A printable
	// character is any character from the Letter, Mark, Number, Punctuation or
	// Symbol categories.
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
