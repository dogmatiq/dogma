package dogma

// ViaAggregate configures the [Application] to route messages to and from an
// [AggregateMessageHandler].
func ViaAggregate(h AggregateMessageHandler, _ ...ViaAggregateOption) ViaAggregateRoute {
	return ViaAggregateRoute{h}
}

// ViaProcess configures the [Application] to route messages to and from a
// [ProcessMessageHandler].
func ViaProcess(h ProcessMessageHandler, _ ...ViaProcessOption) ViaProcessRoute {
	return ViaProcessRoute{h}
}

// ViaIntegration configures the [Application] to route messages to and from an
// [IntegrationMessageHandler].
func ViaIntegration(h IntegrationMessageHandler, _ ...ViaIntegrationOption) ViaIntegrationRoute {
	return ViaIntegrationRoute{h}
}

// ViaProjection configures the [Application] to route messages to a
// [ProjectionMessageHandler].
func ViaProjection(h ProjectionMessageHandler, _ ...ViaProjectionOption) ViaProjectionRoute {
	return ViaProjectionRoute{h}
}

type (
	// HandlerRoute is an interface for types that represent a relationship
	// between the [Application] and a message handler.
	HandlerRoute interface {
		isHandlerRoute()
	}

	// ViaAggregateRoute is a [HandlerRoute] that represents a relationship
	// between the [Application] and an [AggregateMessageHandler].
	ViaAggregateRoute struct{ Handler AggregateMessageHandler }

	// ViaProcessRoute is a [HandlerRoute] that represents a relationship
	// between the [Application] and a [ProcessMessageHandler].
	ViaProcessRoute struct{ Handler ProcessMessageHandler }

	// ViaIntegrationRoute is a [HandlerRoute] that represents a relationship
	// between the [Application] and an [IntegrationMessageHandler].
	ViaIntegrationRoute struct{ Handler IntegrationMessageHandler }

	// ViaProjectionRoute is a [HandlerRoute] that represents a relationship
	// between the [Application] and a [ProjectionMessageHandler].
	ViaProjectionRoute struct{ Handler ProjectionMessageHandler }
)

type (
	// ViaAggregateOption is an option that affects the behavior of a call to
	// the RegisterAggregate() method of the [ApplicationConfigurer] interface.
	ViaAggregateOption interface {
		futureViaAggregateOption()
	}

	// ViaProcessOption is an option that affects the behavior of a call to
	// the RegisterProcess() method of the [ApplicationConfigurer] interface.
	ViaProcessOption interface {
		futureViaProcessOption()
	}

	// ViaIntegrationOption is an option that affects the behavior of a call to
	// the RegisterIntegration() method of the [ApplicationConfigurer] interface.
	ViaIntegrationOption interface {
		futureViaIntegrationOption()
	}

	// ViaProjectionOption is an option that affects the behavior of a call to
	// the RegisterProjection() method of the [ApplicationConfigurer] interface.
	ViaProjectionOption interface {
		futureViaProjectionOption()
	}
)
