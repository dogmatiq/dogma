package dogma

// CommandHandler is an interface implemented by the application and
// used by the engine to handle integration commands.
type CommandHandler interface {
	// Configure configures the behavior of the engine as it relates to this
	// handler.
	//
	// c provides access to the various configuration options, such as specifying
	// which types of integration command messages are routed to this handler.
	Configure(c CommandHandlerConfigurer)

	// HandleCommand handles an integration command message that has been routed to
	// this handler.
	//
	// s provides access to the operations available within the scope of handling
	// m, such as publishing integration event messages.
	//
	// It panics with the UnexpectedMessage value if m is not one of the
	// integration command types that is routed to this handler via Configure().
	HandleCommand(s CommandScope, m Message)
}

// CommandHandlerConfigurer is an interface implemented by the engine and used
// by the application to configure options related to a CommandHandler.
//
// It is passed to CommandMessageHandler.Configure(), typically upon
// initialization of the engine.
//
// In the context of this interface, "the handler" refers to the handler on
// which Configure() has been called.
type CommandHandlerConfigurer interface {
	// RouteCommandType configures the engine to route integration command messages
	// of the same type as m to the handler.
	RouteCommandType(m Message)
}

// CommandScope is an interface implemented by the engine and used by the
// application to perform operations within the context of handling a specific
// integration command message.
type CommandScope interface {
	// RecordEvent records the occurrence of an integration event as a result of
	// the integration command message that is being handled.
	RecordEvent(m Message)

	// Log records an informational message within the context of the integration
	// command message that is being handled.
	Log(f string, v ...interface{})
}
