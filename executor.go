package dogma

import "context"

// A CommandExecutor submits [Command] messages for execution.
//
// It's the primary way that code outside of the Dogma application interacts
// with the engine.
type CommandExecutor interface {
	// ExecuteCommand submits a command for execution.
	//
	// The engine may invoke the associated handler more than once, but the
	// command's side-effects, such as the events it produces, occur exactly
	// once. The engine attempts to execute the command immediately, but there
	// is no guarantee that execution is complete by the time this method
	// returns.
	//
	// If it returns a non-nil error, the engine may not have taken ownership of
	// message delivery, and the application should retry execution.
	ExecuteCommand(context.Context, Command, ...ExecuteCommandOption) error
}

// ExecuteCommandOption is an option that modifies the behavior of
// [CommandExecutor].ExecuteCommand.
type ExecuteCommandOption interface {
	ApplyExecuteCommandOption(executeCommandOptionsBuilder)
}

// WithIdempotencyKey returns an [ExecuteCommandOption] that sets a unique
// identifier for the [Command].
func WithIdempotencyKey(key string) ExecuteCommandOption {
	if key == "" {
		panic("idempotency key cannot be empty")
	}
	return idempotencyKey(key)
}

type executeCommandOptionsBuilder interface {
	IdempotencyKey(string)
}

type idempotencyKey string

func (i idempotencyKey) ApplyExecuteCommandOption(v executeCommandOptionsBuilder) {
	v.IdempotencyKey(string(i))
}
