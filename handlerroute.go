package dogma

// ViaAggregate configures an [Application] to route messages to and from the
// specified [AggregateMessageHandler]. It is used as an argument to the
// Routes() method of [ApplicationConfigurer].
//
// [Command] messages executed using a [CommandExecutor], [ProcessEventScope] or
// [ProcessTimeoutScope] are routed to h if it has a [HandlesCommandRoute] for
// that command type.
//
// [Event] messages recorded by h using an [AggregateCommandScope] are routed to
// other handlers according to their route configurations.
func ViaAggregate(h AggregateMessageHandler, _ ...ViaAggregateOption) ViaAggregateRoute {
	return ViaAggregateRoute{h}
}

// ViaProcess configures an [Application] to route messages to and from the
// specified [ProcessMessageHandler]. It is used as an argument to the Routes()
// method of [ApplicationConfigurer].
//
// [Event] messages recorded using an [AggregateCommandScope],
// [IntegrationCommandScope], or [IntegrationEventScope] are routed to h if it
// has a [HandlesEvent] route for that event type.
//
// [Command] messages executed by h using a [ProcessEventScope] or
// [ProcessTimeoutScope] are routed to other handlers according to their route
// configurations.
//
// [Timeout] messages are always routed back to h itself.
func ViaProcess(h ProcessMessageHandler, _ ...ViaProcessOption) ViaProcessRoute {
	return ViaProcessRoute{h}
}

// ViaIntegration configures an [Application] to route messages to and from the
// specified [IntegrationMessageHandler]. It is used as an argument to the
// Routes() method of [ApplicationConfigurer].
//
// [Command] messages executed using a [CommandExecutor], [ProcessEventScope] or
// [ProcessTimeoutScope] are routed to h if it has a [HandlesCommandRoute] for
// that command type.
//
// [Event] messages recorded using an [AggregateCommandScope] or
// [IntegrationCommandScope] are routed to h if it has a [HandlesEvent] route
// for that event type.
//
// [Event] messages recorded by h using an [IntegrationCommandScope] or
// [IntegrationEventScope] are routed to other handlers according to their route
// configurations.
func ViaIntegration(h IntegrationMessageHandler, _ ...ViaIntegrationOption) ViaIntegrationRoute {
	return ViaIntegrationRoute{h}
}

// ViaProjection configures an [Application] to route messages to the specified
// [ProjectionMessageHandler]. It is used as an argument to the Routes() method
// of [ApplicationConfigurer].
//
// [Event] messages recorded using an [AggregateCommandScope],
// [IntegrationCommandScope], or [IntegrationEventScope] are routed to h if it
// has a [HandlesEvent] route for that event type.
func ViaProjection(h ProjectionMessageHandler, _ ...ViaProjectionOption) ViaProjectionRoute {
	return ViaProjectionRoute{h}
}

type (
	// HandlerRoute is an interface for all types that describe a relationship
	// between an [Application] and the a handler.
	HandlerRoute interface {
		isHandlerRoute()
	}

	// ViaAggregateRoute describes an [AggregateMessageHandler] that is to be
	// registered with an [Application].
	ViaAggregateRoute struct{ Handler AggregateMessageHandler }

	// ViaProcessRoute describes a [ProcessMessageHandler] that is to be
	// registered with an [Application].
	ViaProcessRoute struct{ Handler ProcessMessageHandler }

	// ViaIntegrationRoute describes an [IntegrationMessageHandler] that is
	// to be registered with an [Application].
	ViaIntegrationRoute struct{ Handler IntegrationMessageHandler }

	// ViaProjectionRoute describes a [ProjectionMessageHandler] that is to be
	// registered with an [Application].
	ViaProjectionRoute struct{ Handler ProjectionMessageHandler }
)

type (
	// ViaAggregateOption is an option that affects the behavior of a call to
	// the RegisterAggregate() method of the [ApplicationConfigurer] interface.
	ViaAggregateOption struct{}

	// ViaProcessOption is an option that affects the behavior of a call to
	// the RegisterProcess() method of the [ApplicationConfigurer] interface.
	ViaProcessOption struct{}

	// ViaIntegrationOption is an option that affects the behavior of a call to
	// the RegisterIntegration() method of the [ApplicationConfigurer] interface.
	ViaIntegrationOption struct{}

	// ViaProjectionOption is an option that affects the behavior of a call to
	// the RegisterProjection() method of the [ApplicationConfigurer] interface.
	ViaProjectionOption struct{}
)
