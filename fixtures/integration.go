package fixtures

import (
	"context"
	"time"

	"github.com/dogmatiq/dogma"
)

// IntegrationMessageHandler is a test implementation of
// dogma.IntegrationMessageHandler.
type IntegrationMessageHandler struct {
	ConfigureFunc     func(dogma.IntegrationConfigurer)
	HandleCommandFunc func(context.Context, dogma.IntegrationCommandScope, dogma.Command) error
	TimeoutHintFunc   func(dogma.XMessage) time.Duration
}

var _ dogma.IntegrationMessageHandler = &IntegrationMessageHandler{}

// Configure configures the behavior of the engine as it relates to this
// handler.
//
// If h.ConfigureFunc is non-nil, it calls h.ConfigureFunc(c).
func (h *IntegrationMessageHandler) Configure(c dogma.IntegrationConfigurer) {
	if h.ConfigureFunc != nil {
		h.ConfigureFunc(c)
	}
}

// HandleCommand handles an integration command message that has been routed to
// this handler.
//
// If h.HandleCommandFunc is non-nil it returns h.HandleCommandFunc(s, c),
// otherwise it returns nil.
func (h *IntegrationMessageHandler) HandleCommand(
	ctx context.Context,
	s dogma.IntegrationCommandScope,
	c dogma.Command,
) error {
	if h.HandleCommandFunc != nil {
		return h.HandleCommandFunc(ctx, s, c)
	}

	return nil
}

// TimeoutHint returns a duration that is suitable for computing a deadline for
// the handling of the given message by this handler.
//
// If h.TimeoutHintFunc is non-nil it returns h.TimeoutHintFunc(m), otherwise it
// returns 0.
func (h *IntegrationMessageHandler) TimeoutHint(m dogma.XMessage) time.Duration {
	if h.TimeoutHintFunc != nil {
		return h.TimeoutHintFunc(m)
	}

	return 0
}
