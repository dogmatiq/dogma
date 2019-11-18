package fixtures

import "github.com/dogmatiq/dogma"

// AggregateRoot is a test implementation of dogma.AggregateRoot.
type AggregateRoot struct {
	Value          interface{}
	ApplyEventFunc func(dogma.Message, interface{})
}

var _ dogma.AggregateRoot = &AggregateRoot{}

// ApplyEvent updates the aggregate instance to reflect the fact that a
// particular domain event has occurred.
//
// It calls v.ApplyEventFunc(m, v.Value)
func (v *AggregateRoot) ApplyEvent(m dogma.Message) {
	if v.ApplyEventFunc != nil {
		v.ApplyEventFunc(m, v.Value)
	}
}

// AggregateMessageHandler is a test implementation of dogma.AggregateMessageHandler.
type AggregateMessageHandler struct {
	NewFunc                    func() dogma.AggregateRoot
	ConfigureFunc              func(dogma.AggregateConfigurer)
	RouteCommandToInstanceFunc func(dogma.Message) string
	HandleCommandFunc          func(dogma.AggregateCommandScope, dogma.Message)
}

var _ dogma.AggregateMessageHandler = &AggregateMessageHandler{}

// New constructs a new aggregate instance and returns its root.
//
// If h.NewFunc is nil, it returns a new empty fixtures.AggregateRoot, otherwise
// it calls h.NewFunc().
func (h *AggregateMessageHandler) New() dogma.AggregateRoot {
	if h.NewFunc != nil {
		return h.NewFunc()
	}

	return &AggregateRoot{}
}

// Configure configures the behavior of the engine as it relates to this
// handler.
//
// c provides access to the various configuration options, such as specifying
// which types of domain command messages are routed to this handler.
//
// If h.ConfigureFunc is non-nil, it calls h.ConfigureFunc(c)
func (h *AggregateMessageHandler) Configure(c dogma.AggregateConfigurer) {
	if h.ConfigureFunc != nil {
		h.ConfigureFunc(c)
	}
}

// RouteCommandToInstance returns the ID of the aggregate instance that is
// targetted by m.
//
// It panics with the UnexpectedMessage value if m is not one of the domain
// command types that is routed to this handler via Configure().
//
// If h.RouteCommandToInstanceFunc is non-nil it returns the result of
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
// Handling a domain command message involves inspecting the state of the
// target aggregate instance to determine what changes, if any, should occur.
// Each change is indicated by recording a domain event message.
//
// s provides access to the operations available within the scope of handling
// m, such as creating or destroying the targeted instance, accessing its
// state, and recording domain event messages.
//
// This method must not modify the targeted instance directly. All
// modifications must be applied by the instance's ApplyEvent() method, which
// is called for each domain event message that is recorded via s.
//
// It panics with the UnexpectedMessage value if m is not one of the domain
// command types that is routed to this handler via Configure().
//
// If h.HandleCommandFunc is non-nil it calls h.HandleCommandFunc(s, m).
func (h *AggregateMessageHandler) HandleCommand(s dogma.AggregateCommandScope, m dogma.Message) {
	if h.HandleCommandFunc != nil {
		h.HandleCommandFunc(s, m)
	}
}
