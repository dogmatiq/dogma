package dogma

import "context"

// A CommandExecutor submits [Command] messages for execution.
//
// It's the primary way that code outside of the Dogma application interacts
// with the engine.
type CommandExecutor interface {
	// ExecuteCommand submits a [Command] for execution.
	//
	// It returns once the engine has taken ownership of the command. It doesn't
	// wait for handling to finish.
	//
	// The engine may invoke the command's handler more than once, but the
	// command's side-effects, such as the events it produces, occur exactly
	// once.
	//
	// If it returns a non-nil error, the engine may not have taken ownership of
	// message delivery, and the application should retry execution.
	//
	// See [WithIdempotencyKey].
	ExecuteCommand(context.Context, Command, ...ExecuteCommandOption) error
}

// ExecuteCommandOption is an option that modifies the behavior of
// [CommandExecutor].ExecuteCommand.
type ExecuteCommandOption interface {
	ApplyExecuteCommandOption(executeCommandOptionsBuilder)
}

// WithIdempotencyKey returns an [ExecuteCommandOption] that sets a unique
// identifier for the [Command].
//
// Use an idempotency key when retrying a failed [CommandExecutor].ExecuteCommand
// call to ensure that the engine doesn't execute the command multiple times.
func WithIdempotencyKey(key string) ExecuteCommandOption {
	if key == "" {
		panic("idempotency key cannot be empty")
	}
	return idempotencyKey{key}
}

type executeCommandOptionsBuilder interface {
	IdempotencyKey(string)
}

type idempotencyKey struct{ k string }

func (k idempotencyKey) ApplyExecuteCommandOption(b executeCommandOptionsBuilder) {
	b.IdempotencyKey(k.k)
}
