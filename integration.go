package dogma

import "context"

// IntegrationMessageHandler is an interface implemented by the application and
// used by the engine to integrate with non-message-based systems.
type IntegrationMessageHandler interface {
	// Configure configures the behavior of the engine as it relates to this
	// handler.
	//
	// c provides access to the various configuration options, such as specifying
	// which types of integration command messages are routed to this handler.
	Configure(c IntegrationConfigurer)

	// HandleCommand handles an integration command message that has been routed to
	// this handler.
	//
	// s provides access to the operations available within the scope of handling
	// m, such as publishing integration event messages.
	//
	// It panics with the UnexpectedMessage value if m is not one of the
	// integration command types that is routed to this handler via Configure().
	HandleCommand(ctx context.Context, s IntegrationCommandScope, m Message) error
}

// IntegrationConfigurer is an interface implemented by the engine and used
// by the application to configure options related to a IntegrationMessageHandler.
//
// It is passed to IntegrationMessageHandler.Configure(), typically upon
// initialization of the engine.
//
// In the context of this interface, "the handler" refers to the handler on
// which Configure() has been called.
type IntegrationConfigurer interface {
	// Name sets the name of the handler. Each handler within an application must
	// have a unique name.
	Name(n string)

	// AcceptsCommandType configures the engine to route command messages of the
	// same type as m to the handler.
	AcceptsCommandType(m Message)
}

// IntegrationCommandScope is an interface implemented by the engine and used by
// the application to perform operations within the context of handling a
// specific integration command message.
type IntegrationCommandScope interface {
	// RecordEvent records the occurrence of an integration event as a result of
	// the integration command message that is being handled.
	RecordEvent(m Message)

	// Log records an informational message within the context of the integration
	// command message that is being handled.
	Log(f string, v ...interface{})
}
