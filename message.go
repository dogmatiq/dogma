package dogma

// A Message is an application-defined unit of data that describes a [Command],
// [Event], or [Timeout] within a message-based application.

// A Message describes something an [Application] can do or has done.
//
// Each message type implements either [Command], [Event], or [Timeout].
type Message interface {
	// MessageDescription returns a concise human-readable explanation of the message's
	// meaning for use in contexts such as logging and telemetry.
	//
	// The description should be clear and relevant to developers and
	// non-technical people familiar with the application's domain, but is not
	// intended for end users.
	MessageDescription() string
}

// A Command is a [Message] that instructs an [Application] to perform a specific
// action immediately.
type Command interface {
	Message

	// Validate returns a non-nil error if the message is not well-formed.
	//
	// A command is well-formed if all required information is present and
	// correctly encoded such that it represents a valid action that the
	// application can perform, if current state permits.
	//
	// The [CommandValidationScope] argument exists for forward-compatibility;
	// the interface is currently empty.
	Validate(CommandValidationScope) error
}

// An Event is a [Message] that represents an action that an [Application] has
// performed.
//
// Events capture facts about what has happened within the application and serve
// as a permanent record of past activity.
type Event interface {
	Message

	// Validate returns a non-nil error if the message is not well-formed.
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

// A Timeout is a [Message] that notifies an [Application], specifically a
// [ProcessMessageHandler] that some domain-relevant period of time has elapsed.
type Timeout interface {
	Message

	// Validate returns a non-nil error if the message is not well-formed.
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

// UnexpectedMessage is a panic value used by a message handler when it receives
// a message of a type that it did not expect.
var UnexpectedMessage unexpectedMessage

type unexpectedMessage struct{}

// CommandValidationScope provides context during [Command] validation.
//
// The engine provides the implementation to [Command].Validate.
//
// This interface exists for forward-compatibility.
type CommandValidationScope interface {
	futureCommandValidationScope()
}

// EventValidationScope provides context during [Event] validation.
//
// The engine provides the implementation to [Event].Validate.
type EventValidationScope interface {
	// IsHistorical returns true if the event has already occurred, or false if
	// the application is recording a new event.
	IsHistorical() bool
}

// TimeoutValidationScope provides context during [Timeout] validation.
//
// The engine provides the implementation to [Timeout].Validate.
type TimeoutValidationScope interface {
	// IsScheduled returns true if the timeout is already scheduled to occur, or
	// false if the application is scheduling a new timeout.
	IsScheduled() bool
}
