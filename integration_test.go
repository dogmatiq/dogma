package dogma_test

import (
	"context"
	"testing"

	. "github.com/dogmatiq/dogma"
)

// testIntegrationHandler demonstrates the new HandleEvent capability
type testIntegrationHandler struct{}

func (h *testIntegrationHandler) Configure(c IntegrationConfigurer) {
	c.Identity("test-integration", "12345678-90ab-cdef-1234-567890abcdef")
	c.Routes(
		HandlesCommand[testCommand](),
		HandlesEvent[testEvent](), // New: integration handlers can now handle events
		RecordsEvent[testNotification](),
	)
}

func (h *testIntegrationHandler) HandleCommand(ctx context.Context, s IntegrationCommandScope, cmd Command) error {
	return nil
}

func (h *testIntegrationHandler) HandleEvent(ctx context.Context, s IntegrationEventScope, evt Event) error {
	// NEW: Handle events for things like notifications, external API calls, etc.
	_ = s.RecordedAt() // Use the new RecordedAt method
	s.Log("handling event")
	return nil
}

// Test message types
type testCommand struct{}

func (testCommand) MessageDescription() string            { return "test command" }
func (testCommand) Validate(CommandValidationScope) error { return nil }

type testEvent struct{}

func (testEvent) MessageDescription() string          { return "test event" }
func (testEvent) Validate(EventValidationScope) error { return nil }

type testNotification struct{}

func (testNotification) MessageDescription() string          { return "test notification" }
func (testNotification) Validate(EventValidationScope) error { return nil }

func TestIntegrationHandlerWithEvents(t *testing.T) {
	handler := &testIntegrationHandler{}

	// Test that the handler can be used with ViaIntegration
	route := ViaIntegration(handler)
	if route.Handler != handler {
		t.Fatal("unexpected handler")
	}

	// Test that HandlesEvent route can be used with integration handlers
	eventRoute := HandlesEvent[testEvent]()
	var _ IntegrationRoute = eventRoute // This should compile now - demonstrates the new capability
}
