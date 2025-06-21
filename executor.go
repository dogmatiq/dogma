package dogma

import "context"

// A CommandExecutor submits [Command] messages for execution.
//
// It's the primary way that code outside of the Dogma application interacts
// with the engine.
type CommandExecutor interface {
	// ExecuteCommand submits a command for execution.
	//
	// The engine guarantees execution of the command at some point. The engine
	// may invoke the associated handler more than once, but the command's
	// side-effects, such as the events it produces, occur exactly once.
	//
	// If it returns a non-nil error, the engine may not have taken ownership of
	// message delivery, and the application should retry execution.
	ExecuteCommand(context.Context, Command, ...ExecuteCommandOption) error
}

// ExecuteCommandOption is an option that modifies the behavior of
// [CommandExecutor].ExecuteCommand.
//
// This type exists for forward-compatibility.
type ExecuteCommandOption struct{}
