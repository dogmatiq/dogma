package dogma

// HandlesCommand configures an [AggregateMessageHandler] or
// [IntegrationMessageHandler] as a consumer of [Command] messages of type T.
//
// It panics if T isn't in the message registry, see [RegisterCommand].
//
// Pass the returned [MessageRoute] to [AggregateConfigurer].Routes or
// [IntegrationConfigurer].Routes.
//
// The engine panics if the application has multiple handlers that handle T.
func HandlesCommand[T Command](...HandlesCommandOption) interface {
	AggregateRoute
	IntegrationRoute
} {
	return HandlesCommandRoute{typ: registeredMessageTypeFor[T]()}
}

// ExecutesCommand configures a [ProcessMessageHandler] as a producer of
// [Command] messages of type T.
//
// It panics if T isn't in the message registry, see [RegisterCommand].
//
// Pass the returned [MessageRoute] to [ProcessConfigurer].Routes.
//
// The application may have multiple handlers that execute T.
func ExecutesCommand[T Command](...ExecutesCommandOption) interface {
	ProcessRoute
} {
	return ExecutesCommandRoute{typ: registeredMessageTypeFor[T]()}
}

// RecordsEvent configures an [AggregateMessageHandler] or
// [IntegrationMessageHandler] as a producer of [Event] messages of type T.
//
// It panics if T isn't in the message registry, see [RegisterEvent].
//
// Pass the returned [MessageRoute] to [AggregateConfigurer].Routes or
// [IntegrationConfigurer].Routes.
//
// The engine panics if the application has multiple handlers that record T.
func RecordsEvent[T Event](...RecordsEventOption) interface {
	AggregateRoute
	IntegrationRoute
} {
	return RecordsEventRoute{typ: registeredMessageTypeFor[T]()}
}

// HandlesEvent configures a [ProcessMessageHandler] or [ProjectionMessageHandler] as a
// consumer of [Event] messages of type T.
//
// It panics if T isn't in the message registry, see [RegisterEvent].
//
// Pass the returned [MessageRoute] to [ProcessConfigurer].Routes or
// [ProjectionConfigurer].Routes.
//
// An application may have multiple handlers that handle T.
func HandlesEvent[T Event](...HandlesEventOption) interface {
	ProcessRoute
	ProjectionRoute
} {
	return HandlesEventRoute{typ: registeredMessageTypeFor[T]()}
}

// SchedulesTimeout configures a [ProcessMessageHandler] as a scheduler of
// [Timeout] messages of type T.
//
// It panics if T isn't in the message registry, see [RegisterTimeout].
//
// Pass the returned [MessageRoute] to [ProcessConfigurer].Routes.
//
// The application may have multiple handlers that schedule T.
func SchedulesTimeout[T Timeout](...SchedulesTimeoutOption) interface {
	ProcessRoute
} {
	return SchedulesTimeoutRoute{typ: registeredMessageTypeFor[T]()}
}

type (
	// MessageRoute is an interface for types that describe a relationship between a
	// message handler and a specific message type.
	MessageRoute interface {
		isMessageRoute()
		Type() RegisteredMessageType
	}

	// HandlesCommandRoute is a [HandlerRoute] that represents a handler's
	// ability to handle [Command] messages of a specific type.
	//
	// Use [HandlesCommand] to construct values of this type.
	HandlesCommandRoute struct {
		nocmp
		typ RegisteredMessageType
	}

	// ExecutesCommandRoute is a [HandlerRoute] that represents a handler's
	// ability to execute [Command] messages of a specific type.
	//
	// Use [ExecutesCommand] to construct values of this type.
	ExecutesCommandRoute struct {
		nocmp
		typ RegisteredMessageType
	}

	// HandlesEventRoute is a [HandlerRoute] that represents a handler's
	// ability to handle [Event] messages of a specific type.
	//
	// Use [HandlesEvent] to construct values of this type.
	HandlesEventRoute struct {
		nocmp
		typ RegisteredMessageType
	}

	// RecordsEventRoute is a [HandlerRoute] that represents a handler's
	// ability to record [Event] messages of a specific type.
	//
	// Use [RecordsEvent] to construct values of this type.
	RecordsEventRoute struct {
		nocmp
		typ RegisteredMessageType
	}

	// SchedulesTimeoutRoute is a [HandlerRoute] that represents a handler's
	// ability to schedule [Timeout] messages of a specific type.
	//
	// Use [SchedulesTimeout] to construct values of this type.
	SchedulesTimeoutRoute struct {
		nocmp
		typ RegisteredMessageType
	}
)

type (
	// HandlesCommandOption is an option that modifies the behavior of
	// [HandlesCommand].
	//
	// This type exists for forward-compatibility.
	HandlesCommandOption interface {
		futureHandlesCommandOption()
	}

	// ExecutesCommandOption is an option that modifies the behavior of
	// [ExecutesCommand].
	//
	// This type exists for forward-compatibility.
	ExecutesCommandOption interface {
		futureExecutesCommandOption()
	}

	// HandlesEventOption is an option that modifies the behavior of
	// [HandlesEvent].
	//
	// This type exists for forward-compatibility.
	HandlesEventOption interface {
		futureHandlesEventOption()
	}

	// RecordsEventOption is an option that modifies the behavior of
	// [RecordsEvent].
	//
	// This type exists for forward-compatibility.
	RecordsEventOption interface {
		futureRecordsEventOption()
	}

	// SchedulesTimeoutOption is an option that modifies the behavior of
	// [SchedulesTimeout].
	//
	// This type exists for forward-compatibility.
	SchedulesTimeoutOption interface {
		futureSchedulesTimeoutOption()
	}
)

// Type returns the [RegisteredMessageType] for r.
func (r HandlesCommandRoute) Type() RegisteredMessageType {
	return r.typ
}

// Type returns the [RegisteredMessageType] for r.
func (r ExecutesCommandRoute) Type() RegisteredMessageType {
	return r.typ
}

// Type returns the [RegisteredMessageType] for r.
func (r HandlesEventRoute) Type() RegisteredMessageType {
	return r.typ
}

// Type returns the [RegisteredMessageType] for r.
func (r RecordsEventRoute) Type() RegisteredMessageType {
	return r.typ
}

// Type returns the [RegisteredMessageType] for r.
func (r SchedulesTimeoutRoute) Type() RegisteredMessageType {
	return r.typ
}
