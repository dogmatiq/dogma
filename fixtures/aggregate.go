package fixtures

import "github.com/dogmatiq/dogma"

// AggregateRoot is a test implementation of [dogma.AggregateRoot].
type AggregateRoot struct {
	AppliedEvents  []dogma.Event
	ApplyEventFunc func(dogma.Event) `json:"-"`
}

var _ dogma.AggregateRoot = &AggregateRoot{}

// ApplyEvent updates aggregate instance to reflect the occurrence of an event.
func (v *AggregateRoot) ApplyEvent(e dogma.Event) {
	v.AppliedEvents = append(v.AppliedEvents, e)

	if v.ApplyEventFunc != nil {
		v.ApplyEventFunc(e)
	}
}

// AggregateMessageHandler is a test implementation of
// [dogma.AggregateMessageHandler].
type AggregateMessageHandler struct {
	NewFunc                    func() dogma.AggregateRoot
	ConfigureFunc              func(dogma.AggregateConfigurer)
	RouteCommandToInstanceFunc func(dogma.Command) string
	HandleCommandFunc          func(dogma.AggregateRoot, dogma.AggregateCommandScope, dogma.Command)
}

var _ dogma.AggregateMessageHandler = &AggregateMessageHandler{}

// Configure describes the handler's configuration to the engine.
func (h *AggregateMessageHandler) Configure(c dogma.AggregateConfigurer) {
	if h.ConfigureFunc != nil {
		h.ConfigureFunc(c)
	}
}

// New returns an aggregate root instance in its initial state.
func (h *AggregateMessageHandler) New() dogma.AggregateRoot {
	if h.NewFunc != nil {
		return h.NewFunc()
	}
	return &AggregateRoot{}
}

// RouteCommandToInstance returns the ID of the instance that handles a specific
// command.
func (h *AggregateMessageHandler) RouteCommandToInstance(c dogma.Command) string {
	if h.RouteCommandToInstanceFunc == nil {
		panic(dogma.UnexpectedMessage)
	}

	return h.RouteCommandToInstanceFunc(c)
}

// HandleCommand executes business logic in response to a command.
func (h *AggregateMessageHandler) HandleCommand(
	r dogma.AggregateRoot,
	s dogma.AggregateCommandScope,
	c dogma.Command,
) {
	if h.HandleCommandFunc != nil {
		h.HandleCommandFunc(r, s, c)
	}
}
