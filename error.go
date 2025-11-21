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
)
