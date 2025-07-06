package dogma

// ViaAggregate configures the [Application] to route messages to and from an
// [AggregateMessageHandler].
//
// Pass the returned [HandlerRoute] to [ApplicationConfigurer].Routes.
func ViaAggregate(h AggregateMessageHandler, _ ...ViaAggregateOption) ViaAggregateRoute {
	return ViaAggregateRoute{h}
}

// ViaProcess configures the [Application] to route messages to and from a
// [ProcessMessageHandler].
//
// Pass the returned [HandlerRoute] to [ApplicationConfigurer].Routes.
func ViaProcess(h ProcessMessageHandler, _ ...ViaProcessOption) ViaProcessRoute {
	return ViaProcessRoute{h}
}

// ViaIntegration configures the [Application] to route messages to and from an
// [IntegrationMessageHandler].
//
// Pass the returned [HandlerRoute] to [ApplicationConfigurer].Routes.
func ViaIntegration(h IntegrationMessageHandler, _ ...ViaIntegrationOption) ViaIntegrationRoute {
	return ViaIntegrationRoute{h}
}

// ViaProjection configures the [Application] to route messages to a
// [ProjectionMessageHandler].
//
// Pass the returned [HandlerRoute] to [ApplicationConfigurer].Routes.
func ViaProjection(h ProjectionMessageHandler, _ ...ViaProjectionOption) ViaProjectionRoute {
	return ViaProjectionRoute{h}
}

type (
	// HandlerRoute is an interface for types that represent a relationship
	// between the [Application] and a message handler.
	//
	// Use [ViaAggregate], [ViaProcess], [ViaIntegration], or [ViaProjection]
	// to create a HandlerRoute.
	HandlerRoute interface{ isHandlerRoute() }

	// ViaAggregateRoute is a [HandlerRoute] that represents a relationship
	// between the [Application] and an [AggregateMessageHandler].
	//
	// Avoid constructing values of this type directly; use [ViaAggregate]
	// instead.
	ViaAggregateRoute struct{ Handler AggregateMessageHandler }

	// ViaProcessRoute is a [HandlerRoute] that represents a relationship
	// between the [Application] and a [ProcessMessageHandler].
	//
	// Avoid constructing values of this type directly; use [ViaProcess]
	// instead.
	ViaProcessRoute struct{ Handler ProcessMessageHandler }

	// ViaIntegrationRoute is a [HandlerRoute] that represents a relationship
	// between the [Application] and an [IntegrationMessageHandler].
	//
	// Avoid constructing values of this type directly; use [ViaIntegration]
	// instead.
	ViaIntegrationRoute struct{ Handler IntegrationMessageHandler }

	// ViaProjectionRoute is a [HandlerRoute] that represents a relationship
	// between the [Application] and a [ProjectionMessageHandler].
	//
	// Avoid constructing values of this type directly; use [ViaProjection]
	// instead.
	ViaProjectionRoute struct{ Handler ProjectionMessageHandler }
)

type (
	// ViaAggregateOption is an option that modifies the behavior of
	// [ViaAggregate].
	//
	// This type exists for forward compatibility.
	ViaAggregateOption interface {
		futureViaAggregateOption()
	}

	// ViaProcessOption is an option that modifies the behavior of
	// [ViaProcess].
	//
	// This type exists for forward compatibility.
	ViaProcessOption interface {
		futureViaProcessOption()
	}

	// ViaIntegrationOption is an option that modifies the behavior of
	// [ViaIntegration].
	//
	// This type exists for forward compatibility.
	ViaIntegrationOption interface {
		futureViaIntegrationOption()
	}

	// ViaProjectionOption is an option that modifies the behavior of
	// [ViaProjection].
	//
	// This type exists for forward compatibility.
	ViaProjectionOption interface {
		futureViaProjectionOption()
	}
)
