package dogma

// An Application defines a collection of message handlers that work together to
// implement a specific business domain.
//
// An application in the general sense must provide at least one Dogma
// [Application] implementation.
type Application interface {
	// Configure declares the application's configuration by calling methods on c.
	//
	// The configuration includes the application's identity and handler routes.
	//
	// The engine calls this method at least once during startup. If called more
	// than once, it must produce the same configuration each time.
	Configure(c ApplicationConfigurer)
}

// ApplicationConfigurer is the interface an [Application] uses to declare its
// configuration.
//
// The engine provides the implementation to [Application].Configure during
// startup.
type ApplicationConfigurer interface {
	// Identity configures the application's unique identity.
	//
	// n is a short human-readable name used in logs, telemetry, and other
	// places where the application's identity appears. The value must be
	// between 1 and 255 bytes in length, and contain only printable, non-space
	// UTF-8 characters. Changing the application's name does not affect its
	// behavior.
	//
	// k is a key that uniquely identifies the application. The engine uses the
	// key to associate application state with the correct application instance,
	// so it must not change. The value must be a canonical RFC 4122 UUID
	// string, such as "5195fe85-eb3f-4121-84b0-be72cbc5722f", and is
	// case-insensitive.
	Identity(n, k string)

	// Routes configures the application to route messages to and from specific
	// message handlers.
	//
	// It accepts routes created by [ViaAggregate], [ViaProcess],
	// [ViaIntegration], and [ViaProjection].
	//
	// The application does not declare routes for specific message types
	// directly; it inherits routes from the handlers it contains.
	Routes(...HandlerRoute)
}
