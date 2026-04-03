package dogma

import "errors"

var (
	// ErrNotSupported is the error returned when a feature is not supported
	// by a particular implementation.
	//
	// It may be returned by:
	//  - [AggregateRoot].MarshalBinary
	//  - [AggregateRoot].UnmarshalBinary
	//  - [ProjectionMessageHandler].Reset
	ErrNotSupported = errors.New("not supported")

	// ErrEventObserverNotSatisfied is returned when [WithEventObserver] is used
	// and the engine determines that no further relevant events can occur
	// before any [EventObserver] returned satisfied == true.
	//
	// It may be returned by:
	//  - [CommandExecutor].ExecuteCommand
	ErrEventObserverNotSatisfied = errors.New("event observer not satisfied")
)
