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

// ProcessMessageHandler is a test implementation of
// dogma.ProcessMessageHandler.
type ProcessMessageHandler struct {
	NewFunc                  func() dogma.ProcessRoot
	ConfigureFunc            func(dogma.ProcessConfigurer)
	RouteEventToInstanceFunc func(context.Context, dogma.Event) (string, bool, error)
	HandleEventFunc          func(context.Context, dogma.ProcessRoot, dogma.ProcessEventScope, dogma.Event) error
	HandleTimeoutFunc        func(context.Context, dogma.ProcessRoot, dogma.ProcessTimeoutScope, dogma.Timeout) error
	TimeoutHintFunc          func(dogma.Message) time.Duration
}

var _ dogma.ProcessMessageHandler = &ProcessMessageHandler{}

// New constructs a new process instance and returns its root.
//
// If h.NewFunc is non-nil, it returns h.NewFunc(), otherwise it returns a new
// empty fixtures.ProcessRoot.
func (h *ProcessMessageHandler) New() dogma.ProcessRoot {
	if h.NewFunc != nil {
		return h.NewFunc()
	}

	return &ProcessRoot{}
}

// Configure configures the behavior of the engine as it relates to this
// handler.
//
// If h.ConfigureFunc is non-nil, it calls h.ConfigureFunc(c).
func (h *ProcessMessageHandler) Configure(c dogma.ProcessConfigurer) {
	if h.ConfigureFunc != nil {
		h.ConfigureFunc(c)
	}
}

// RouteEventToInstance returns the ID of the process instance that is targeted
// by e.
//
// If h.RouteEventToInstance is non-nil it returns h.RouteEventToInstance(ctx,
// e), otherwise it panics.
func (h *ProcessMessageHandler) RouteEventToInstance(
	ctx context.Context,
	e dogma.Event,
) (string, bool, error) {
	if h.RouteEventToInstanceFunc == nil {
		panic(dogma.UnexpectedMessage)
	}

	return h.RouteEventToInstanceFunc(ctx, e)
}

// HandleEvent handles an event message that has been routed to this handler.
//
// If h.HandleEventFunc is non-nil it returns h.HandleEventFunc(ctx, r, s, e),
// otherwise it returns nil.
func (h *ProcessMessageHandler) HandleEvent(
	ctx context.Context,
	r dogma.ProcessRoot,
	s dogma.ProcessEventScope,
	e dogma.Event,
) error {
	if h.HandleEventFunc != nil {
		return h.HandleEventFunc(ctx, r, s, e)
	}

	return nil
}

// HandleTimeout handles a timeout message that has been scheduled with
// ProcessScope.ScheduleTimeout().
//
// If h.HandleTimeoutFunc is non-nil it returns h.HandleTimeoutFunc(ctx, r, s, t),
// otherwise it returns nil.
func (h *ProcessMessageHandler) HandleTimeout(
	ctx context.Context,
	r dogma.ProcessRoot,
	s dogma.ProcessTimeoutScope,
	t dogma.Timeout,
) error {
	if h.HandleTimeoutFunc != nil {
		return h.HandleTimeoutFunc(ctx, r, s, t)
	}

	return nil
}

// TimeoutHint returns a duration that is suitable for computing a deadline for
// the handling of the given message by this handler.
//
// If h.TimeoutHintFunc is non-nil it returns h.TimeoutHintFunc(m), otherwise it
// returns 0.
func (h *ProcessMessageHandler) TimeoutHint(m dogma.Message) time.Duration {
	if h.TimeoutHintFunc != nil {
		return h.TimeoutHintFunc(m)
	}

	return 0
}
