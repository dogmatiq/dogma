package dogma

import "context"

// ProjectionMessageHandler is an interface implemented by the application and
// used by the engine to build a "projection" (also known as a "read model", or
// "query model") from events that occur within the application.
type ProjectionMessageHandler interface {
	// Configure produces a configuration for this handler by calling methods on
	// the configurer c.
	//
	// The implementation MUST allow for multiple calls to Configure(). Each
	// call SHOULD produce the same configuration.
	//
	// The engine MUST call Configure() before calling HandleEvent(). It is
	// RECOMMENDED that the engine only call Configure() once per handler.
	Configure(c ProjectionConfigurer)

	// HandleEvent updates the projection to reflect the occurrence of event
	// message m.
	//
	// If nil is returned, the projection has been updated successfully.
	//
	// If an error is returned, the projection SHOULD be left in the state it
	// was before HandleEvent() was called.
	//
	// The engine SHOULD provide "at-least-once" delivery gaurantees to the
	// handler. That is, the engine should call HandleEvent() with the same
	// event message until a nil error is returned.
	//
	// The engine MUST NOT call HandleEvent() with any message of a type that
	// has not been routed to this handler by a prior call to Configure(). If
	// any such message is passed, the implementation MUST panic with the
	// UnexpectedMessage value.
	//
	// The engine MAY provide gaurantees about the order in which event messages
	// will be passed to HandleEvent(), however in the interest of engine
	// portability the implementation SHOULD NOT assume that HandleEvent() will
	// be called with events in the same order that they were recorded.
	//
	// The engine MAY call HandleEvent() from multiple goroutines concurrently.
	HandleEvent(ctx context.Context, s ProjectionEventScope, m Message) error
}

// ProjectionConfigurer is an interface implemented by the engine and used by
// the application to configure options related to a ProjectionMessageHandler.
//
// It is passed to ProjectionMessageHandler.Configure(), typically upon
// initialization of the engine.
//
// In the context of this interface, "the handler" refers to the handler on
// which Configure() has been called.
type ProjectionConfigurer interface {
	// Name sets the name of the handler.
	//
	// It MUST be called exactly once within a single call to Configure().
	//
	// Each handler within an application MUST have a unique, non-empty name.
	Name(n string)

	// RouteEventType configures the engine to route event messages of the same
	// type as m to the handler.
	//
	// It MUST be called at least once within a call to Configure(). It MUST NOT
	// be called more than once with an event message of the same type.
	//
	// Multiple handlers within a single application MAY receive event messages
	// of the same type.
	//
	// The "content" of m MUST NOT be used, inspected, or treated as meaningful
	// in any way, only its runtime type information.
	RouteEventType(m Message)
}

// ProjectionEventScope is an interface implemented by the engine and used by
// the application to perform operations within the context of handling a
// specific event message.
type ProjectionEventScope interface {
	// Log records an informational message within the context of the event
	// message that is being handled.
	Log(f string, v ...interface{})
}
