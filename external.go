package dogma

import "context"

// CommandExecutor is an interface implemented by the engine and used by the
// application to execute commands outside of a Dogma message handler.
//
// Engine implementations MUST provide mechanisms for executing commands outside
// of any message handler. It is RECOMMENDED that one such mechanism adhere to
// this interface.
type CommandExecutor interface {
	// ExecuteCommand enqueues a command for execution.
	//
	// If nil is returned, the engine MUST provide reasonable guarantees that
	// the message will be dispatched to the appropriate handler. For example,
	// the engine is responsible for retrying the message in the face of a
	// transient error.
	//
	// If an error is returned, the command may not have been enqueued. The
	// application SHOULD retry executing the command.
	//
	// The engine MAY handle the message synchronously, however the application
	// SHOULD NOT assume that the message has been handled immediately.
	//
	// If ctx has a deadline, it MUST NOT be used by the engine as a mechanism
	// for message expiration.
	ExecuteCommand(ctx context.Context, m Message) error
}

// EventRecorder is an interface implemented by the engine and used by the
// application to record events outside of a Dogma message handler.
//
// Engine implementations MAY provide mechanisms for recording events outside of
// any message handler, in which case it is RECOMMENDED that one such mechanism
// adhere to this interface.
type EventRecorder interface {
	// RecordEvent records the occurrence of an event.
	//
	// If nil is returned, the engine MUST provide reasonable guarantees that
	// the message will be dispatched to the appropriate handler. For example,
	// the engine is responsible for retrying the message in the face of a
	// transient error.
	//
	// If an error is returned, the event may not have been recorded. The
	// application SHOULD retry recording the event.
	//
	// The engine MAY handle the message synchronously, however the application
	// SHOULD NOT assume that the message has been handled immediately.
	//
	// If ctx has a deadline, it MUST NOT be used by the engine as a mechanism
	// for message expiration.
	RecordEvent(ctx context.Context, m Message) error
}
