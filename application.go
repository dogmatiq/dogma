package dogma

// An Application is a collection of message handlers that model a single
// logical business domain.
type Application interface {
	// Configure describes the application's configuration to the engine.
	Configure(ApplicationConfigurer)
}

// An ApplicationConfigurer configures the engine for use with a specific
// application.
type ApplicationConfigurer interface {
	// Identity configures the application's identity.
	//
	// n is a short human-readable name. It MAY change over the application's
	// lifetime. It MUST contain solely printable, non-space UTF-8 characters.
	// It must be between 1 and 255 bytes (not characters) in length.
	//
	// k is a unique key used to associate engine state with the application.
	// The key SHOULD NOT change over the application's lifetime. k MUST be an
	// RFC 4122 UUID, such as "5195fe85-eb3f-4121-84b0-be72cbc5722f".
	//
	// Use of hard-coded literals for both values is RECOMMENDED.
	Identity(n string, k string)

	// Routes configures the application to route messages via specific message
	// handlers.
	Routes(...HandlerRoute)
}
