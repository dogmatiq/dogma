package dogma

import "context"

// A CommandExecutor executes a command from outside the context of any message
// handler.
//
// The CommandExecutor is the primary way that code outside of the Dogma
// application interacts with the Dogma engine.
type CommandExecutor interface {
	// ExecuteCommand executes or enqueues a command.
	//
	// If it returns nil, the engine has guaranteed execution of the command.
	// Otherwise, the it's the caller's responsibility to retry.
	//
	// The application SHOULD assume that the command is executed
	// asynchronously; it has not necessarily executed by the time the method
	// returns.
	ExecuteCommand(context.Context, Command) error
}

// EventRecorder records events outside of a Dogma message handler.
//
// Deprecated: No production engines implement this interface. To record
// arbitrary events implement an [IntegrationMessageHandler] instead.
type EventRecorder interface {
	// RecordEvent records the occurrence of an event.
	RecordEvent(context.Context, Event) error
}
