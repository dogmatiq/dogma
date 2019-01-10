package dogma

// App is a definition of a Dogma application.
type App struct {
	// Name is a unique name for the application.
	//
	// The engine may make use of the application name for message routing.
	Name string

	// Aggregates is a collection of the aggregate message handlers that the
	// application contains.
	Aggregates []AggregateMessageHandler

	// Aggregates is a collection of the process message handlers that the
	// application contains.
	Processes []ProcessMessageHandler

	// Aggregates is a collection of the integration message handlers that the
	// application contains.
	Integrations []IntegrationMessageHandler

	// Aggregates is a collection of the projection message handlers that the
	// application contains.
	Projections []ProjectionMessageHandler
}
