package fixtures

import (
	"context"

	"github.com/dogmatiq/dogma"
)

// ProjectionMessageHandler is a test implementation of
// [dogma.ProjectionMessageHandler].
type ProjectionMessageHandler struct {
	ConfigureFunc       func(dogma.ProjectionConfigurer)
	HandleEventFunc     func(context.Context, []byte, []byte, []byte, dogma.ProjectionEventScope, dogma.Event) (bool, error)
	ResourceVersionFunc func(context.Context, []byte) ([]byte, error)
	CloseResourceFunc   func(context.Context, []byte) error
	CompactFunc         func(context.Context, dogma.ProjectionCompactScope) error
}

var _ dogma.ProjectionMessageHandler = &ProjectionMessageHandler{}

// Configure describes the handler's configuration to the engine.
func (h *ProjectionMessageHandler) Configure(c dogma.ProjectionConfigurer) {
	if h.ConfigureFunc != nil {
		h.ConfigureFunc(c)
	}
}

// HandleEvent updates the projection to reflect the occurrence of an event.
func (h *ProjectionMessageHandler) HandleEvent(
	ctx context.Context,
	r, c, n []byte,
	s dogma.ProjectionEventScope,
	e dogma.Event,
) (bool, error) {
	if h.HandleEventFunc != nil {
		return h.HandleEventFunc(ctx, r, c, n, s, e)
	}
	return true, nil
}

// ResourceVersion returns the current version of a resource.
func (h *ProjectionMessageHandler) ResourceVersion(
	ctx context.Context,
	r []byte,
) ([]byte, error) {
	if h.ResourceVersionFunc != nil {
		return h.ResourceVersionFunc(ctx, r)
	}
	return nil, nil
}

// CloseResource informs the handler that the engine has no further use for
// a resource.
func (h *ProjectionMessageHandler) CloseResource(
	ctx context.Context,
	r []byte,
) error {
	if h.CloseResourceFunc != nil {
		return h.CloseResourceFunc(ctx, r)
	}
	return nil
}

// Compact attempts to reduce the size of the projection.
func (h *ProjectionMessageHandler) Compact(
	ctx context.Context,
	s dogma.ProjectionCompactScope,
) error {
	if h.CompactFunc != nil {
		return h.CompactFunc(ctx, s)
	}
	return nil
}
