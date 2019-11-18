package fixtures

import (
	"context"
	"time"

	"github.com/dogmatiq/dogma"
)

// ProcessRoot is a test implementation of dogma.ProcessRoot.
type ProcessRoot struct {
	Value interface{}
}

var _ dogma.ProcessRoot = &ProcessRoot{}

// ProcessMessageHandler is a test implementation of dogma.ProcessMessageHandler.
type ProcessMessageHandler struct {
	NewFunc                  func() dogma.ProcessRoot
	ConfigureFunc            func(dogma.ProcessConfigurer)
	RouteEventToInstanceFunc func(context.Context, dogma.Message) (string, bool, error)
	HandleEventFunc          func(context.Context, dogma.ProcessEventScope, dogma.Message) error
	HandleTimeoutFunc        func(context.Context, dogma.ProcessTimeoutScope, dogma.Message) error
	TimeoutHintFunc          func(m dogma.Message) time.Duration
}

var _ dogma.ProcessMessageHandler = &ProcessMessageHandler{}

// New constructs a new process instance and returns its root.
//
// If h.NewFunc is nil, it returns a new empty fixtures.ProcessRoot, otherwise
// it calls h.NewFunc().
func (h *ProcessMessageHandler) New() dogma.ProcessRoot {
	if h.NewFunc != nil {
		return h.NewFunc()
	}

	return &ProcessRoot{}
}

// Configure configures the behavior of the engine as it relates to this
// handler.
//
// c provides access to the various configuration options, such as
// specifying which types of event messages are routed to this handler.
//
// If h.ConfigureFunc is non-nil, it calls h.ConfigureFunc(c)
func (h *ProcessMessageHandler) Configure(c dogma.ProcessConfigurer) {
	if h.ConfigureFunc != nil {
		h.ConfigureFunc(c)
	}
}

// RouteEventToInstance returns the ID of the process instance that is
// targetted by m.
//
// It panics with the UnexpectedMessage value if m is not one of the event
// types that is routed to this handler via Configure().
//
// If ok is false, the message is not routed to this handler at all.
//
// If h.RouteEventToInstance is non-nil it returns the result of
// h.RouteEventToInstance(ctx, m), otherwise it panics.
func (h *ProcessMessageHandler) RouteEventToInstance(
	ctx context.Context,
	m dogma.Message,
) (string, bool, error) {
	if h.RouteEventToInstanceFunc == nil {
		panic(dogma.UnexpectedMessage)
	}

	return h.RouteEventToInstanceFunc(ctx, m)
}

// HandleEvent handles an event message that has been routed to this
// handler.
//
// Handling an event message involves inspecting the state of the target
// process instance to determine what command messages, if any, should be
// produced.
//
// s provides access to the operations available within the scope of handling
// m, such as beginning or ending the targeted instance, accessing its state,
// sending command messages or scheduling timeouts.
//
// This method may manipulate the process's state directly.
//
// It panics with the UnexpectedMessage value if m is not one of the event
// types that is routed to this handler via Configure().
//
// If h.HandleEventFunc is non-nil it calls h.HandleEventFunc(ctx, s, m).
func (h *ProcessMessageHandler) HandleEvent(
	ctx context.Context,
	s dogma.ProcessEventScope,
	m dogma.Message,
) error {
	if h.HandleEventFunc != nil {
		return h.HandleEventFunc(ctx, s, m)
	}

	return nil
}

// HandleTimeout handles a timeout message that has been scheduled with
// ProcessScope.ScheduleTimeout().
//
// Timeouts can be used to model time within the domain. For example, an
// application might use a timeout to mark an invoice as overdue after some
// period of non-payment.
//
// Handling a timeout is much like handling an event in that much the same
// operations are available to the handler via s.
//
// This method may manipulate the process's state directly.
//
// If m was not expected by the handler the implementation must panic with an
// UnexpectedMessage value.
//
// If h.HandleTimeoutFunc is non-nil it calls h.HandleTimeoutFunc(ctx, s, m).
func (h *ProcessMessageHandler) HandleTimeout(
	ctx context.Context,
	s dogma.ProcessTimeoutScope,
	m dogma.Message,
) error {
	if h.HandleTimeoutFunc != nil {
		return h.HandleTimeoutFunc(ctx, s, m)
	}

	return nil
}

// TimeoutHint returns a duration that is suitable for computing a deadline
// for the handling of the given message by this handler.
//
// If h.TimeoutHintFunc is non-nil it calls h.TimeoutHintFunc(m).
func (h *ProcessMessageHandler) TimeoutHint(m dogma.Message) time.Duration {
	if h.TimeoutHintFunc != nil {
		return h.TimeoutHintFunc(m)
	}

	return 0
}
