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
	ExecuteCommand(context.Context, Command, ...ExecuteCommandOption) error
}

// ExecuteCommandOption is an option that affects the behavior of a call to the
// ExecuteCommand() method of the [CommandExecutor] interface.
type ExecuteCommandOption struct {
	idempotencyKey string
}

// WithIdempotencyKey returns an ExecuteCommandOption that specifies an
// idempotency key for command execution.
//
// The idempotency key allows applications to provide strict idempotency when
// using CommandExecutor. This is not a substitute for domain-level idempotency,
// but it can help a client safely retry when the domain-level idempotency is
// poorly implemented.
func WithIdempotencyKey(key string) ExecuteCommandOption {
	return ExecuteCommandOption{
		idempotencyKey: key,
	}
}

// IdempotencyKey returns the idempotency key specified for this option.
// It returns an empty string if no idempotency key was specified.
func (o ExecuteCommandOption) IdempotencyKey() string {
	return o.idempotencyKey
}
