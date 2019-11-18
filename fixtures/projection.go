package fixtures

import (
	"context"
	"time"

	"github.com/dogmatiq/dogma"
)

// ProjectionMessageHandler is a test implementation of dogma.ProjectionMessageHandler.
type ProjectionMessageHandler struct {
	ConfigureFunc       func(dogma.ProjectionConfigurer)
	HandleEventFunc     func(context.Context, []byte, []byte, []byte, dogma.ProjectionEventScope, dogma.Message) (bool, error)
	ResourceVersionFunc func(context.Context, []byte) ([]byte, error)
	CloseResourceFunc   func(context.Context, []byte) error
	TimeoutHintFunc     func(m dogma.Message) time.Duration
}

var _ dogma.ProjectionMessageHandler = &ProjectionMessageHandler{}

// Configure configures the behavior of the engine as it relates to this
// handler.
//
// c provides access to the various configuration options, such as specifying
// which types of event messages are routed to this handler.
//
// If h.ConfigureFunc is non-nil, it calls h.ConfigureFunc(c)
func (h *ProjectionMessageHandler) Configure(c dogma.ProjectionConfigurer) {
	if h.ConfigureFunc != nil {
		h.ConfigureFunc(c)
	}
}

// HandleEvent handles a domain event message that has been routed to this
// handler.
//
// s provides access to the operations available within the scope of handling m.
//
// It panics with the UnexpectedMessage value if m is not one of the event
// types that is routed to this handler via Configure().
//
// If h.HandleEventFunc is non-nil it calls h.HandleEventFunc(ctx, r,c,n, s, m).
func (h *ProjectionMessageHandler) HandleEvent(
	ctx context.Context,
	r, c, n []byte,
	s dogma.ProjectionEventScope,
	m dogma.Message,
) (bool, error) {
	if h.HandleEventFunc != nil {
		return h.HandleEventFunc(ctx, r, c, n, s, m)
	}

	return true, nil
}

// ResourceVersion returns the version of the resource r.
//
// If h.ResourceVersionFunc is non-nil it calls h.ResourceVersionFunc(ctx, k).
func (h *ProjectionMessageHandler) ResourceVersion(ctx context.Context, k []byte) ([]byte, error) {
	if h.ResourceVersionFunc != nil {
		return h.ResourceVersionFunc(ctx, k)
	}

	return nil, nil
}

// CloseResource informs the projection that the resource r will not be used in
// any future calls to HandleEvent().
//
// If h.CloseResourceFunc is non-nil it calls h.CloseResourceFunc(ctx, k).
func (h *ProjectionMessageHandler) CloseResource(ctx context.Context, k []byte) error {
	if h.CloseResourceFunc != nil {
		return h.CloseResourceFunc(ctx, k)
	}

	return nil
}

// TimeoutHint returns a duration that is suitable for computing a deadline
// for the handling of the given message by this handler.
//
// If h.TimeoutHintFunc is non-nil it calls h.TimeoutHintFunc(m).
func (h *ProjectionMessageHandler) TimeoutHint(m dogma.Message) time.Duration {
	if h.TimeoutHintFunc != nil {
		return h.TimeoutHintFunc(m)
	}

	return 0
}
