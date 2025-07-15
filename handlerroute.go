package dogma

// ViaAggregate configures the [Application] to route messages to and from an
// [AggregateMessageHandler].
//
// Pass the returned [HandlerRoute] to [ApplicationConfigurer].Routes.
func ViaAggregate(h AggregateMessageHandler, _ ...ViaAggregateOption) HandlerRoute {
	return ViaAggregateRoute{Handler: h}
}

// ViaProcess configures the [Application] to route messages to and from a
// [ProcessMessageHandler].
//
// Pass the returned [HandlerRoute] to [ApplicationConfigurer].Routes.
func ViaProcess(h ProcessMessageHandler, _ ...ViaProcessOption) HandlerRoute {
	return ViaProcessRoute{Handler: h}
}

// ViaIntegration configures the [Application] to route messages to and from an
// [IntegrationMessageHandler].
//
// Pass the returned [HandlerRoute] to [ApplicationConfigurer].Routes.
func ViaIntegration(h IntegrationMessageHandler, _ ...ViaIntegrationOption) HandlerRoute {
	return ViaIntegrationRoute{Handler: h}
}

// ViaProjection configures the [Application] to route messages to a
// [ProjectionMessageHandler].
//
// Pass the returned [HandlerRoute] to [ApplicationConfigurer].Routes.
func ViaProjection(h ProjectionMessageHandler, _ ...ViaProjectionOption) HandlerRoute {
	return ViaProjectionRoute{Handler: h}
}

type (
	// HandlerRoute is an interface for types that represent a relationship
	// between the [Application] and a message handler.
	//
	// Use [ViaAggregate], [ViaProcess], [ViaIntegration], or [ViaProjection]
	// to create a HandlerRoute.
	HandlerRoute interface {
		ApplyHandlerRoute(HandlerRoutesBuilder)
	}

	// ViaAggregateRoute is a [HandlerRoute] that represents a relationship
	// between the [Application] and an [AggregateMessageHandler].
	//
	// Avoid constructing values of this type directly; use [ViaAggregate]
	// instead.
	ViaAggregateRoute struct {
		nocmp
		Handler AggregateMessageHandler
	}

	// ViaProcessRoute is a [HandlerRoute] that represents a relationship
	// between the [Application] and a [ProcessMessageHandler].
	//
	// Avoid constructing values of this type directly; use [ViaProcess]
	// instead.
	ViaProcessRoute struct {
		nocmp
		Handler ProcessMessageHandler
	}

	// ViaIntegrationRoute is a [HandlerRoute] that represents a relationship
	// between the [Application] and an [IntegrationMessageHandler].
	//
	// Avoid constructing values of this type directly; use [ViaIntegration]
	// instead.
	ViaIntegrationRoute struct {
		nocmp
		Handler IntegrationMessageHandler
	}

	// ViaProjectionRoute is a [HandlerRoute] that represents a relationship
	// between the [Application] and a [ProjectionMessageHandler].
	//
	// Avoid constructing values of this type directly; use [ViaProjection]
	// instead.
	ViaProjectionRoute struct {
		nocmp
		Handler ProjectionMessageHandler
	}
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

// HandlerRoutesBuilder is an interface for types that can build configuration
// from [HandlerRoute] values.
//
// This type is part of the engine configuration API. It's not intended for
// use during application development.
type HandlerRoutesBuilder interface {
	ViaAggregate(ViaAggregateRoute)
	ViaProcess(ViaProcessRoute)
	ViaIntegration(ViaIntegrationRoute)
	ViaProjection(ViaProjectionRoute)
}

// ApplyHandlerRoute passes r to [HandlerRoutesBuilder].ViaAggregate.
func (r ViaAggregateRoute) ApplyHandlerRoute(b HandlerRoutesBuilder) { b.ViaAggregate(r) }

// ApplyHandlerRoute passes r to [HandlerRoutesBuilder].ViaProcess.
func (r ViaProcessRoute) ApplyHandlerRoute(b HandlerRoutesBuilder) { b.ViaProcess(r) }

// ApplyHandlerRoute passes r to [HandlerRoutesBuilder].ViaIntegration.
func (r ViaIntegrationRoute) ApplyHandlerRoute(b HandlerRoutesBuilder) { b.ViaIntegration(r) }

// ApplyHandlerRoute passes r to [HandlerRoutesBuilder].ViaProjection.
func (r ViaProjectionRoute) ApplyHandlerRoute(b HandlerRoutesBuilder) { b.ViaProjection(r) }
