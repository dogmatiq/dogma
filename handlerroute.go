package dogma

// ViaAggregate configures the [Application] to route messages to and from an
// [AggregateMessageHandler].
//
// Pass the returned [HandlerRoute] to [ApplicationConfigurer].Routes.
func ViaAggregate(h AggregateMessageHandler, _ ...ViaAggregateOption) HandlerRoute {
	if h == nil {
		panic("handler cannot be nil")
	}
	return AggregateHandlerRoute{handler: h}
}

// ViaProcess configures the [Application] to route messages to and from a
// [ProcessMessageHandler].
//
// Pass the returned [HandlerRoute] to [ApplicationConfigurer].Routes.
func ViaProcess(h ProcessMessageHandler, _ ...ViaProcessOption) HandlerRoute {
	if h == nil {
		panic("handler cannot be nil")
	}
	return ProcessHandlerRoute{handler: h}
}

// ViaIntegration configures the [Application] to route messages to and from an
// [IntegrationMessageHandler].
//
// Pass the returned [HandlerRoute] to [ApplicationConfigurer].Routes.
func ViaIntegration(h IntegrationMessageHandler, _ ...ViaIntegrationOption) HandlerRoute {
	if h == nil {
		panic("handler cannot be nil")
	}
	return IntegrationHandlerRoute{handler: h}
}

// ViaProjection configures the [Application] to route messages to a
// [ProjectionMessageHandler].
//
// Pass the returned [HandlerRoute] to [ApplicationConfigurer].Routes.
func ViaProjection(h ProjectionMessageHandler, _ ...ViaProjectionOption) HandlerRoute {
	if h == nil {
		panic("handler cannot be nil")
	}
	return ProjectionHandlerRoute{handler: h}
}

type (
	// HandlerRoute is an interface for types that represent a relationship
	// between the [Application] and a message handler.
	//
	// Use [ViaAggregate], [ViaProcess], [ViaIntegration], or [ViaProjection]
	// to create a HandlerRoute.
	HandlerRoute interface {
		isHandlerRoute()
	}

	// AggregateHandlerRoute is a [HandlerRoute] that represents a relationship
	// between the [Application] and an [AggregateMessageHandler].
	//
	// Use [ViaAggregate] to construct values of this type.
	AggregateHandlerRoute struct {
		nocmp
		handler AggregateMessageHandler
	}

	// ProcessHandlerRoute is a [HandlerRoute] that represents a relationship
	// between the [Application] and a [ProcessMessageHandler].
	//
	// See [ViaProcess] to construct values of this type.
	ProcessHandlerRoute struct {
		nocmp
		handler ProcessMessageHandler
	}

	// IntegrationHandlerRoute is a [HandlerRoute] that represents a
	// relationship between the [Application] and an
	// [IntegrationMessageHandler].
	//
	// Use [ViaIntegration] to construct values of this type.
	IntegrationHandlerRoute struct {
		nocmp
		handler IntegrationMessageHandler
	}

	// ProjectionHandlerRoute is a [HandlerRoute] that represents a relationship
	// between the [Application] and a [ProjectionMessageHandler].
	//
	// Use [ViaProjection] to construct values of this type.
	ProjectionHandlerRoute struct {
		nocmp
		handler ProjectionMessageHandler
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

// Handler returns the [AggregateMessageHandler] for r.
func (r AggregateHandlerRoute) Handler() AggregateMessageHandler {
	return r.handler
}

// Handler returns the [ProcessMessageHandler] for r.
func (r ProcessHandlerRoute) Handler() ProcessMessageHandler {
	return r.handler
}

// Handler returns the [IntegrationMessageHandler] for r.
func (r IntegrationHandlerRoute) Handler() IntegrationMessageHandler {
	return r.handler
}

// Handler returns the [ProjectionMessageHandler] for r.
func (r ProjectionHandlerRoute) Handler() ProjectionMessageHandler {
	return r.handler
}
