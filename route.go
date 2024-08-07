package dogma

import "reflect"

// HandlesCommand routes command messages to an [AggregateMessageHandler] or
// [IntegrationMessageHandler].
//
// It's used as an argument to the Routes() method of [AggregateConfigurer] or
// [IntegrationConfigurer].
//
// An application MUST NOT route a single command type to more than one handler.
func HandlesCommand[T Command](...HandlesCommandOption) HandlesCommandRoute {
	return HandlesCommandRoute{typeOf[T]()}
}

// RecordsEvent routes event messages recorded by an [AggregateMessageHandler]
// or [IntegrationMessageHandler].
//
// It's used as an argument to the Routes() method of [AggregateConfigurer] or
// [IntegrationConfigurer].
//
// An application MUST NOT route a single event type from more than one handler.
func RecordsEvent[T Event](...RecordsEventOption) RecordsEventRoute {
	return RecordsEventRoute{typeOf[T]()}
}

// HandlesEvent routes event messages to a [ProcessMessageHandler] or
// [ProjectionMessageHandler].
//
// It's used as an argument to the Routes() method of [ProcessConfigurer] or
// [ProjectionConfigurer].
func HandlesEvent[T Event](...HandlesEventOption) HandlesEventRoute {
	return HandlesEventRoute{typeOf[T]()}
}

// ExecutesCommand routes command messages produced by a
// [ProcessMessageHandler].
//
// It's used as an argument to the Routes() method of [ProcessConfigurer].
func ExecutesCommand[T Command](...ExecutesCommandOption) ExecutesCommandRoute {
	return ExecutesCommandRoute{typeOf[T]()}
}

// SchedulesTimeout routes timeout messages scheduled by
// [ProcessMessageHandler].
//
// It's used as an argument to the Routes() method of [ProcessConfigurer].
//
// An application MAY use a single timeout type with more than one process.
func SchedulesTimeout[T Timeout](...SchedulesTimeoutOption) SchedulesTimeoutRoute {
	return SchedulesTimeoutRoute{typeOf[T]()}
}

type (
	// Route is an interface implemented by all route types.
	Route interface{ isRoute() }

	// HandlesCommandRoute describes a route for a handler that handles a
	// [Command] of a specific type.
	HandlesCommandRoute struct{ Type reflect.Type }

	// ExecutesCommandRoute describes a route for a handler that executes a
	// [Command] of a specific type.
	ExecutesCommandRoute struct{ Type reflect.Type }

	// HandlesEventRoute describes a route for a handler that handles an
	// [Event] of a specific type.
	HandlesEventRoute struct{ Type reflect.Type }

	// RecordsEventRoute describes a route for a handler that records an
	// [Event] of a specific type.
	RecordsEventRoute struct{ Type reflect.Type }

	// SchedulesTimeoutRoute describes a route for a handler that schedules a
	// [Timeout] of a specific type.
	SchedulesTimeoutRoute struct{ Type reflect.Type }
)

type (
	// HandlesCommandOption is an option that affects the behavior of the route
	// returned by [HandlesCommand].
	HandlesCommandOption struct{}

	// ExecutesCommandOption is an option that affects the behavior of the route
	// returned by [ExecutesCommand].
	ExecutesCommandOption struct{}

	// HandlesEventOption is an option that affects the behavior of the route
	// returned by [HandlesEvent].
	HandlesEventOption struct{}

	// RecordsEventOption is an option that affects the behavior of the route
	// returned by [RecordsEvent].
	RecordsEventOption struct{}

	// SchedulesTimeoutOption is an option that affects the behavior of the
	// route returned by [SchedulesTimeout].
	SchedulesTimeoutOption struct{}
)

func typeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}
