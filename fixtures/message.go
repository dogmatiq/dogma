package fixtures

import (
	"errors"
	"fmt"
)

// Command is an implementation of dogma.Command used for testing.
type Command[T any] struct {
	Content T
	Invalid string
}

// MessageDescription returns a description of the command.
func (c Command[T]) MessageDescription() string {
	validity := "valid"
	if c.Invalid != "" {
		validity = "invalid: " + c.Invalid
	}
	return fmt.Sprintf(
		"command(%T:%v, %s)",
		c.Content,
		c.Content,
		validity,
	)
}

// Validate returns a non-nil error if c.Invalid is not empty.
func (c Command[T]) Validate() error {
	if c.Invalid != "" {
		return errors.New(c.Invalid)
	}
	return nil
}

// Event is an implementation of dogma.Event used for testing.
type Event[T any] struct {
	Content T
	Invalid string
}

// MessageDescription returns a description of the command.
func (e Event[T]) MessageDescription() string {
	validity := "valid"
	if e.Invalid != "" {
		validity = "invalid: " + e.Invalid
	}
	return fmt.Sprintf(
		"event(%T:%v, %s)",
		e.Content,
		e.Content,
		validity,
	)
}

// Validate returns a non-nil error if c.Invalid is not empty.
func (e Event[T]) Validate() error {
	if e.Invalid != "" {
		return errors.New(e.Invalid)
	}
	return nil
}

// Timeout is an implementation of dogma.Timeout used for testing.
type Timeout[T any] struct {
	Content T
	Invalid string
}

// MessageDescription returns a description of the command.
func (t Timeout[T]) MessageDescription() string {
	validity := "valid"
	if t.Invalid != "" {
		validity = "invalid: " + t.Invalid
	}
	return fmt.Sprintf(
		"timeout(%T:%v, %s)",
		t.Content,
		t.Content,
		validity,
	)
}

// Validate returns a non-nil error if c.Invalid is not empty.
func (t Timeout[T]) Validate() error {
	if t.Invalid != "" {
		return errors.New(t.Invalid)
	}
	return nil
}
