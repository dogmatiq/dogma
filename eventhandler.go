package dogma

// EventHandler is an interface implemented by the application and
// used by the engine to handle domain or integration event messages.
type EventHandler interface {
	// Configure configures the behavior of the engine as it relates to this
	// handler.
	//
	// c provides access to the various configuration options, such as specifying
	// which types of event messages are routed to this handler.
	Configure(c EventHandlerConfigurer)

	// HandleEvent handles an event message that has been routed to this handler.
	//
	// s provides access to the operations available within the scope of handling
	// m.
	//
	// It panics with the UnexpectedMessage value if m is not one of the event
	// types that is routed to this handler via Configure().
	HandleEvent(s EventScope, m Message)
}

// EventHandlerConfigurer is an interface implemented by the engine and used
// by the application to configure options related to an EventHandler.
//
// It is passed to EventMessageHandler.Configure(), typically upon
// initialization of the engine.
//
// In the context of this interface, "the handler" refers to the handler on
// which Configure() has been called.
type EventHandlerConfigurer interface {
	// RouteEventType configures the engine to route events of the same type as m
	// to the handler.
	RouteEventType(m Message)
}

// EventScope is an interface implemented by the engine and used by the
// application to perform operations within the context of handling a specific
// event message.
type EventScope interface {
	// Log records an informational message within the context of the event
	// message that is being handled.
	Log(f string, v ...interface{})
}
