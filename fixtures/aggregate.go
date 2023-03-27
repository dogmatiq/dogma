package fixtures

import "github.com/dogmatiq/dogma"

// AggregateRoot is a test implementation of dogma.AggregateRoot.
type AggregateRoot struct {
	AppliedEvents  []dogma.Event
	ApplyEventFunc func(dogma.Event) `json:"-"`
}

var _ dogma.AggregateRoot = &AggregateRoot{}

// ApplyEvent updates the aggregate instance to reflect the fact that a
// particular domain event has occurred.
//
// It appends e to v.AppliedEvents.
//
// If v.ApplyEventFunc is non-nil, it calls v.ApplyEventFunc(e, v.Value).
func (v *AggregateRoot) ApplyEvent(e dogma.Event) {
	v.AppliedEvents = append(v.AppliedEvents, e)

	if v.ApplyEventFunc != nil {
		v.ApplyEventFunc(e)
	}
}

// AggregateMessageHandler is a test implementation of
// dogma.AggregateMessageHandler.
type AggregateMessageHandler struct {
	NewFunc                    func() dogma.AggregateRoot
	ConfigureFunc              func(dogma.AggregateConfigurer)
	RouteCommandToInstanceFunc func(dogma.Command) string
	HandleCommandFunc          func(dogma.AggregateRoot, dogma.AggregateCommandScope, dogma.Command)
}

var _ dogma.AggregateMessageHandler = &AggregateMessageHandler{}

// New constructs a new aggregate instance and returns its root.
//
// If h.NewFunc is non-nil, it returns h.NewFunc(), otherwise it returns a new
// empty new empty fixtures.AggregateRoot.
func (h *AggregateMessageHandler) New() dogma.AggregateRoot {
	if h.NewFunc != nil {
		return h.NewFunc()
	}

	return &AggregateRoot{}
}

// Configure configures the behavior of the engine as it relates to this
// handler.
//
// If h.ConfigureFunc is non-nil, it calls h.ConfigureFunc(c).
func (h *AggregateMessageHandler) Configure(c dogma.AggregateConfigurer) {
	if h.ConfigureFunc != nil {
		h.ConfigureFunc(c)
	}
}

// RouteCommandToInstance returns the ID of the aggregate instance that is
// targeted by c.
//
// If h.RouteCommandToInstanceFunc is non-nil it returns
// h.RouteCommandToInstanceFunc(c), otherwise it panics.
func (h *AggregateMessageHandler) RouteCommandToInstance(c dogma.Command) string {
	if h.RouteCommandToInstanceFunc == nil {
		panic(dogma.UnexpectedMessage)
	}

	return h.RouteCommandToInstanceFunc(c)
}

// HandleCommand handles a domain command message that has been routed to this
// handler.
//
// If h.HandleCommandFunc is non-nil it calls h.HandleCommandFunc(r, s, c).
func (h *AggregateMessageHandler) HandleCommand(
	r dogma.AggregateRoot,
	s dogma.AggregateCommandScope,
	c dogma.Command,
) {
	if h.HandleCommandFunc != nil {
		h.HandleCommandFunc(r, s, c)
	}
}
