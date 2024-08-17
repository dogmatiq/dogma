package fixtures

import (
	"context"

	"github.com/dogmatiq/dogma"
)

// IntegrationMessageHandler is a test implementation of
// [dogma.IntegrationMessageHandler].
type IntegrationMessageHandler struct {
	ConfigureFunc     func(dogma.IntegrationConfigurer)
	HandleCommandFunc func(context.Context, dogma.IntegrationCommandScope, dogma.Command) error
}

var _ dogma.IntegrationMessageHandler = &IntegrationMessageHandler{}

// Configure describes the handler's configuration to the engine.
func (h *IntegrationMessageHandler) Configure(c dogma.IntegrationConfigurer) {
	if h.ConfigureFunc != nil {
		h.ConfigureFunc(c)
	}
}

// HandleCommand handles a command, typically by invoking some external API.
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
