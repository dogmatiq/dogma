package dogma

import "context"

// A CommandExecutor submits [Command] messages for execution.
//
// It's the primary way that code outside of the Dogma application interacts
// with the engine.
type CommandExecutor interface {
	// ExecuteCommand submits a [Command] for execution.
	//
	// It returns once the engine has taken ownership of the command. By
	// default, it doesn't wait for handling to finish.
	//
	// The engine may invoke the command's handler more than once, but the
	// command's side-effects, such as the events it produces, occur exactly
	// once.
	//
	// If it returns a non-nil error, the engine may not have taken ownership of
	// message delivery, and the application should retry execution.
	//
	// If the application intends to retry execution of the command, use
	// [WithIdempotencyKey] with the same key on the initial call and on every
	// retry. This is typically unnecessary if the application's logic already
	// guarantees that retrying cannot produce duplicate side-effects. When
	// using an idempotency key, the caller is responsible for retrying until
	// ExecuteCommand returns nil. Otherwise, the command may be lost.
	//
	// Pass [WithEventObserver] to wait until a specific event is recorded
	// during command execution. If the engine determines that no further
	// relevant events can occur and none of the observers returned satisfied ==
	// true, ExecuteCommand returns [ErrEventObserverNotSatisfied].
	ExecuteCommand(context.Context, Command, ...ExecuteCommandOption) error
}

// ExecuteCommandOption is an option that modifies the behavior of
// [CommandExecutor].ExecuteCommand.
type ExecuteCommandOption interface {
	isExecuteCommandOption()
}

// WithIdempotencyKey returns an [ExecuteCommandOption] that sets a stable
// identifier for the [Command].
//
// Pass the idempotency key on the initial [CommandExecutor].ExecuteCommand call
// and reuse the same key on every retry to ensure that the engine doesn't
// execute the command multiple times.
//
// By using an idempotency key, the caller accepts responsibility for retrying
// until [CommandExecutor].ExecuteCommand returns nil. Otherwise, the command
// may be lost.
func WithIdempotencyKey(key string) ExecuteCommandOption {
	if key == "" {
		panic("idempotency key cannot be empty")
	}

	return IdempotencyKeyOption{key: key}
}

// IdempotencyKeyOption is an [ExecuteCommandOption] that sets a unique
// identifier for a [Command].
//
// Use [WithIdempotencyKey] to create values of this type.
type IdempotencyKeyOption struct {
	nocmp
	key string
}

// Key returns the idempotency key.
func (o IdempotencyKeyOption) Key() string {
	return o.key
}

// EventObserver is a callback invoked by the engine for each event of type T
// recorded as a result of executing a command with
// [CommandExecutor].ExecuteCommand.
//
// See [WithEventObserver].
type EventObserver[T Event] func(ctx context.Context, event T) (satisfied bool, err error)

// WithEventObserver returns an [ExecuteCommandOption] that observes events of
// type T recorded while executing a command.
//
// It panics if T isn't in the message registry, see [RegisterEvent].
//
// [CommandExecutor].ExecuteCommand blocks until the observer returns
// satisfied == true, the observer returns a non-nil error, the caller's context
// ends, or the engine determines that no further relevant events can occur.
//
// Multiple WithEventObserver options may be passed to a single call.
// ExecuteCommand unblocks as soon as any one of the observers is satisfied.
func WithEventObserver[T Event](fn EventObserver[T]) ExecuteCommandOption {
	if fn == nil {
		panic("event observer cannot be nil")
	}

	typ := registeredMessageTypeFor[T]()
	return EventObserverOption{
		eventType: typ,
		observer: func(ctx context.Context, e Event) (bool, error) {
			if v, ok := e.(T); ok {
				return fn(ctx, v)
			}
			return false, nil
		},
	}
}

// EventObserverOption is an [ExecuteCommandOption] that observes events of a
// specific type recorded during command execution.
//
// Use [WithEventObserver] to construct values of this type.
type EventObserverOption struct {
	nocmp
	eventType RegisteredMessageType
	observer  func(context.Context, Event) (bool, error)
}

// EventType returns the event type that this option observes.
func (o EventObserverOption) EventType() RegisteredMessageType {
	return o.eventType
}

// Observer returns the type-erased observer function.
func (o EventObserverOption) Observer() func(context.Context, Event) (bool, error) {
	return o.observer
}
