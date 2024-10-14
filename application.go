package dogma

// An Application is a collection of message handlers that model a single
// logical business domain.
type Application interface {
	// Configure describes the application's configuration to the engine.
	Configure(ApplicationConfigurer)
}

// An ApplicationConfigurer configures the engine for use with a specific
// application.
type ApplicationConfigurer interface {
	// Identity configures the application's identity.
	//
	// n is a short human-readable name. It MAY change over the application's
	// lifetime. It MUST contain solely printable, non-space UTF-8 characters.
	// It must be between 1 and 255 bytes (not characters) in length.
	//
	// k is a unique key used to associate engine state with the application.
	// The key SHOULD NOT change over the application's lifetime. k MUST be an
	// RFC 4122 UUID, such as "5195fe85-eb3f-4121-84b0-be72cbc5722f".
	//
	// Use of hard-coded literals for both values is RECOMMENDED.
	Identity(n string, k string)

	// Handlers configures the application to route messages to and from
	// specific message handlers.
	Handlers(...HandlerRegistration)

	// RegisterAggregate configures the engine to route messages for an
	// aggregate.
	//
	// Deprecated: Pass the result of [RegisterAggregate] to the Handlers()
	// method instead.
	RegisterAggregate(AggregateMessageHandler, ...RegisterAggregateOption)

	// RegisterProcess configures the engine to route messages for a process.
	//
	// Deprecated: Pass the result of [RegisterProcess] to the Handlers()
	// method instead.
	RegisterProcess(ProcessMessageHandler, ...RegisterProcessOption)

	// RegisterIntegration configures the engine to route messages for an
	// integration.
	//
	// Deprecated: Pass the result of [RegisterIntegration] to the Handlers()
	// method instead.
	RegisterIntegration(IntegrationMessageHandler, ...RegisterIntegrationOption)

	// RegisterProjection configures the engine to route messages for a
	// projection.
	//
	// Deprecated: Pass the result of [RegisterProjection] to the Handlers()
	// method instead.
	RegisterProjection(ProjectionMessageHandler, ...RegisterProjectionOption)
}

// RegisterAggregate registers an [AggregateMessageHandler] with an
// [Application].
//
// It is used as an argument to the Handlers() method of
// [ApplicationConfigurer].
func RegisterAggregate(h AggregateMessageHandler, _ ...RegisterAggregateOption) AggregateRegistration {
	return AggregateRegistration{h}
}

// RegisterProcess registers a [ProcessMessageHandler] with an [Application].
//
// It is used as an argument to the Handlers() method of
// [ApplicationConfigurer].
func RegisterProcess(h ProcessMessageHandler, _ ...RegisterProcessOption) ProcessRegistration {
	return ProcessRegistration{h}
}

// RegisterIntegration registers an [IntegrationMessageHandler] with an
// [Application].
//
// It is used as an argument to the Handlers() method of
// [ApplicationConfigurer].
func RegisterIntegration(h IntegrationMessageHandler, _ ...RegisterIntegrationOption) IntegrationRegistration {
	return IntegrationRegistration{h}
}

// RegisterProjection registers a [ProjectionMessageHandler] with an
// [Application].
//
// It is used as an argument to the Handlers() method of
// [ApplicationConfigurer].
func RegisterProjection(h ProjectionMessageHandler, _ ...RegisterProjectionOption) ProjectionRegistration {
	return ProjectionRegistration{h}
}

type (
	// HandlerRegistration is an interface implemented by all handler
	// registration types.
	HandlerRegistration interface{ isHandlerRegistration() }

	// AggregateRegistration describes an [AggregateMessageHandler] that is to be
	// registered with an [Application].
	AggregateRegistration struct{ Handler AggregateMessageHandler }

	// ProcessRegistration describes a [ProcessMessageHandler] that is to be
	// registered with an [Application].
	ProcessRegistration struct{ Handler ProcessMessageHandler }

	// IntegrationRegistration describes an [IntegrationMessageHandler] that is
	// to be registered with an [Application].
	IntegrationRegistration struct{ Handler IntegrationMessageHandler }

	// ProjectionRegistration describes a [ProjectionMessageHandler] that is to
	// be registered with an [Application].
	ProjectionRegistration struct{ Handler ProjectionMessageHandler }
)

type (
	// RegisterAggregateOption is an option that affects the behavior of a call to
	// the RegisterAggregate() method of the [ApplicationConfigurer] interface.
	RegisterAggregateOption struct{}

	// RegisterProcessOption is an option that affects the behavior of a call to
	// the RegisterProcess() method of the [ApplicationConfigurer] interface.
	RegisterProcessOption struct{}

	// RegisterIntegrationOption is an option that affects the behavior of a call to
	// the RegisterIntegration() method of the [ApplicationConfigurer] interface.
	RegisterIntegrationOption struct{}

	// RegisterProjectionOption is an option that affects the behavior of a call to
	// the RegisterProjection() method of the [ApplicationConfigurer] interface.
	RegisterProjectionOption struct{}
)
