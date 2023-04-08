package fixtures

import "github.com/dogmatiq/dogma"

// Application is a test implementation of [dogma.Application].
type Application struct {
	ConfigureFunc func(dogma.ApplicationConfigurer)
}

var _ dogma.Application = &Application{}

// Configure describes the application's configuration to the engine.
func (a *Application) Configure(c dogma.ApplicationConfigurer) {
	if a.ConfigureFunc != nil {
		a.ConfigureFunc(c)
	}
}
