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
	// Name sets the name of the application.
	//
	// It MUST be called exactly once within a single call to Configure().
	//
	// The name MUST be a non-empty UTF-8 string consisting solely of printable
	// unicode characters. A printable character is any character from the
	// Letter, Mark, Number, Punctuation or Symbol categories.
	//
	// Although not recommended, the application MAY share its name with one of
	// its handlers.
	Name(n string)

	// RegisterAggregate configures the engine to route messages to h.
	RegisterAggregate(h AggregateMessageHandler)

	// RegisterProcess configures the engine to route messages to h.
	RegisterProcess(h ProcessMessageHandler)

	// RegisterIntegration configures the engine to route messages to h.
	RegisterIntegration(h IntegrationMessageHandler)

	// RegisterProjection configures the engine to route messages to h.
	RegisterProjection(h ProjectionMessageHandler)
}
