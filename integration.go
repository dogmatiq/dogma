package dogma

type CommandHandler interface {
	Configure(c CommandHandlerConfigurer)
	HandleCommand(s CommandScope, m Message)
}

type CommandHandlerConfigurer interface {
	RouteCommandType(m Message)
}

type CommandScope interface {
	PublishEvent(m Message)
	Log(f string, v ...interface{})
}

type EventHandler interface {
	Configure(c EventHandlerConfigurer)
	HandleEvent(s EventScope, m Message)
}

type EventHandlerConfigurer interface {
	RouteEventType(m Message)
}

type EventScope interface {
	Log(f string, v ...interface{})
}
