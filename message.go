package dogma

import "time"

// A Message is an application-defined unit of data that describes a [Command],
// [Event], or [Timeout] within a message-based application.

// A Message describes something an [Application] can do or has done.
//
// Each message type implements either [Command], [Event], or [Timeout].
type Message interface {
	// MessageDescription returns a concise human-readable explanation of the
	// message's meaning for use in contexts such as logging and telemetry.
	//
	// The description should be clear and relevant to developers and
	// non-technical stakeholders familiar with the application's domain. It's
	// not intended for display to end users.
	//
	// Use lowercase sentences with no trailing punctuation. Omit sensitive
	// information and overly specific details that don't alter the message's
	// intent.
	//
	// Descriptions of [Command] messages should use present-continuous tense.
	// For example: "adding 10 widgets to Alex's shopping cart".
	//
	// Descriptions of [Event] messages should use past tense. For example:
	// "added 10 widgets to Alex's shopping cart".
	//
	// Descriptions of [Timeout] messages should read as though the timeout has
	// just elapsed. For example: "Alex's cart is now inactive" or "24 hours
	// elapsed since first item added to Alex's cart".
	//
	// Be wary of assuming a specific actor if the message doesn't explicitly
	// encode that information. For example, prefer "Alex's purchase completed"
	// over "Alex completed their purchase". This guidance is especially
	// relevant to [Event] messages, where each type should represent a specific
	// state change regardless of who initiated it.
	MessageDescription() string

	// MarshalBinary returns the message's binary representation, suitable for
	// persistence or transmission over the network.
	MarshalBinary() ([]byte, error)

	// UnmarshalBinary populates the message from its binary representation.
	//
	// The implementation must clone the data if it is used after returning.
	UnmarshalBinary([]byte) error
}

// UnexpectedMessage is a panic value used by a message handler when it receives
// a message of a type that it didn't expect.
var UnexpectedMessage unexpectedMessage

type unexpectedMessage struct{}

// MessageValidationScope provides context during [Message] validation.
//
// Each message type has a corresponding scope type that
// extends this interface:
//
//   - [CommandValidationScope]
//   - [EventValidationScope]
//   - [TimeoutValidationScope]
type MessageValidationScope interface {
	// IsNew returns true when a handler or [CommandExecutor] is creating a new
	// message, or false when the engine is re-validating a message that it has
	// already accepted.
	//
	// Use the distinction to apply different validation rules for new messages
	// versus existing messages while keeping all validation logic in one
	// location.
	IsNew() bool
}

// A Command is a [Message] that instructs an [Application] to perform a specific
// action immediately.
type Command interface {
	Message

	// Validate returns a non-nil error if the message isn't well-formed.
	//
	// A command is well-formed if all required information is present and
	// correctly encoded such that it represents a valid action that the
	// application can perform, if current state permits.
	//
	// The [CommandValidationScope] argument exists for forward-compatibility;
	// the interface is currently empty.
	Validate(CommandValidationScope) error
}

// CommandValidationScope provides context during [Command] validation.
//
// The engine provides the implementation to [Command].Validate.
//
// This type exists for forward-compatibility.
type CommandValidationScope interface {
	MessageValidationScope

	// ExecutedAt returns the time at which the application submitted the
	// command for execution by calling ExecuteCommand() on a [CommandExecutor]
	// or [ProcessScope].
	ExecutedAt() time.Time
}

// An Event is a [Message] that represents an action that an [Application] has
// performed.
//
// Events capture facts about what has happened within the application and serve
// as a permanent record of past activity.
type Event interface {
	Message

	// Validate returns a non-nil error if the message isn't well-formed.
	//
	// An event is well-formed if all required information is present and
	// correctly encoded such that it accurately represents an action that can
	// occur within the application.
	//
	// Validation requirements may change over time. Use the
	// [EventValidationScope] to access context that may affect the strictness
	// or criteria of the validation logic.
	Validate(EventValidationScope) error
}

// EventValidationScope provides context during [Event] validation.
//
// The engine provides the implementation to [Event].Validate.
type EventValidationScope interface {
	MessageValidationScope

	// RecordedAt returns the time at which the event occurred.
	RecordedAt() time.Time
}

// A Timeout is a [Message] that notifies an [Application], specifically a
// [ProcessMessageHandler] that some domain-relevant period of time has elapsed.
type Timeout interface {
	Message

	// Validate returns a non-nil error if the message isn't well-formed.
	//
	// A timeout is well-formed if all required information is present and
	// correctly encoded such that it accurately represents an action that can
	// occur within the process, if current state permits.
	//
	// Validation requirements may change over time. Use the
	// [TimeoutValidationScope] to access context that may affect the strictness
	// or criteria of the validation logic.
	Validate(TimeoutValidationScope) error
}

// TimeoutValidationScope provides context during [Timeout] validation.
//
// The engine provides the implementation to [Timeout].Validate.
type TimeoutValidationScope interface {
	MessageValidationScope

	// ScheduledAt returns the time at which the handler scheduled the timeout.
	ScheduledAt() time.Time

	// ScheduledFor returns the time at which the timeout occurs.
	ScheduledFor() time.Time
}
