package dogma

// An Application is a collection of message handlers that model a single
// logical business domain.
type Application interface {
	// Configure describes the application's configuration to the engine.
	Configure(ApplicationConfigurer)
}

// An ApplicationConfigurer configures the engine for use with a specific
// application.
//
// See [ApplicationMessageHandler.Configure]().
type ApplicationConfigurer interface {
	// Identity configures the application's identity.
	//
	// n is a short human-readable name. The name MAY change over the
	// application's lifetime. n MUST contain solely printable, non-space UTF-8
	// characters.
	//
	// k is a unique key used to associate engine state with the application.
	// The key SHOULD NOT change over the application's lifetime. k MUST be a an
	// [RFC 4122] UUID, such as "5195fe85-eb3f-4121-84b0-be72cbc5722f".
	//
	// Use of hard-coded literals for both values is RECOMMENDED.
	Identity(n string, k string)

	// RegisterAggregate configures the engine to route messages for an
	// aggregate.
	RegisterAggregate(AggregateMessageHandler)

	// RegisterProcess configures the engine to route messages for a process.
	RegisterProcess(ProcessMessageHandler)

	// RegisterIntegration configures the engine to route messages for an
	// integration.
	RegisterIntegration(IntegrationMessageHandler)

	// RegisterProjection configures the engine to route messages for a
	// projection.
	RegisterProjection(ProjectionMessageHandler)
}
