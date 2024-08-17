package fixtures

import (
	"context"

	"github.com/dogmatiq/dogma"
)

// ProcessRoot is a test implementation of [dogma.ProcessRoot].
type ProcessRoot struct {
	Value any
}

var _ dogma.ProcessRoot = &ProcessRoot{}

// ProcessMessageHandler is a test implementation of
// [dogma.ProcessMessageHandler].
type ProcessMessageHandler struct {
	NewFunc                  func() dogma.ProcessRoot
	ConfigureFunc            func(dogma.ProcessConfigurer)
	RouteEventToInstanceFunc func(context.Context, dogma.Event) (string, bool, error)
	HandleEventFunc          func(context.Context, dogma.ProcessRoot, dogma.ProcessEventScope, dogma.Event) error
	HandleTimeoutFunc        func(context.Context, dogma.ProcessRoot, dogma.ProcessTimeoutScope, dogma.Timeout) error
}

var _ dogma.ProcessMessageHandler = &ProcessMessageHandler{}

// Configure describes the handler's configuration to the engine.
func (h *ProcessMessageHandler) Configure(c dogma.ProcessConfigurer) {
	if h.ConfigureFunc != nil {
		h.ConfigureFunc(c)
	}
}

// New returns a process root instance in its initial state.
func (h *ProcessMessageHandler) New() dogma.ProcessRoot {
	if h.NewFunc != nil {
		return h.NewFunc()
	}
	return &ProcessRoot{}
}

// RouteEventToInstance returns the ID of the instance that handles a specific
// event.
func (h *ProcessMessageHandler) RouteEventToInstance(
	ctx context.Context,
	e dogma.Event,
) (string, bool, error) {
	if h.RouteEventToInstanceFunc == nil {
		panic(dogma.UnexpectedMessage)
	}
	return h.RouteEventToInstanceFunc(ctx, e)
}

// HandleEvent begins or continues the process in response to an event.
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

// HandleTimeout continues the process in response to a timeout.
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
