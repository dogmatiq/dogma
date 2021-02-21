package fixtures

import "github.com/dogmatiq/dogma"

// AggregateRoot is a test implementation of dogma.AggregateRoot.
type AggregateRoot struct {
	AppliedEvents  []dogma.Message
	ApplyEventFunc func(dogma.Message) `json:"-"`
}

var _ dogma.AggregateRoot = &AggregateRoot{}

// ApplyEvent updates the aggregate instance to reflect the fact that a
// particular domain event has occurred.
//
// It appends m to v.AppliedEvents.
//
// If v.ApplyEventFunc is non-nil, it calls v.ApplyEventFunc(m, v.Value).
func (v *AggregateRoot) ApplyEvent(m dogma.Message) {
	v.AppliedEvents = append(v.AppliedEvents, m)

	if v.ApplyEventFunc != nil {
		v.ApplyEventFunc(m)
	}
}

// AggregateMessageHandler is a test implementation of
// dogma.AggregateMessageHandler.
type AggregateMessageHandler struct {
	NewFunc                    func() dogma.AggregateRoot
	ConfigureFunc              func(dogma.AggregateConfigurer)
	RouteCommandToInstanceFunc func(dogma.Message) string
	HandleCommandFunc          func(dogma.AggregateRoot, dogma.AggregateCommandScope, dogma.Message)
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
// targeted by m.
//
// If h.RouteCommandToInstanceFunc is non-nil it returns
// h.RouteCommandToInstanceFunc(m), otherwise it panics.
func (h *AggregateMessageHandler) RouteCommandToInstance(m dogma.Message) string {
	if h.RouteCommandToInstanceFunc == nil {
		panic(dogma.UnexpectedMessage)
	}

	return h.RouteCommandToInstanceFunc(m)
}

// HandleCommand handles a domain command message that has been routed to this
// handler.
//
// If h.HandleCommandFunc is non-nil it calls h.HandleCommandFunc(r, s, m).
func (h *AggregateMessageHandler) HandleCommand(
	r dogma.AggregateRoot,
	s dogma.AggregateCommandScope,
	m dogma.Message,
) {
	if h.HandleCommandFunc != nil {
		h.HandleCommandFunc(r, s, m)
	}
}
