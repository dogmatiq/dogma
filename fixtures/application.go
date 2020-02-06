package fixtures

import "github.com/dogmatiq/dogma"

// Application is a test implementation of dogma.Application.
type Application struct {
	ConfigureFunc func(dogma.ApplicationConfigurer)
}

var _ dogma.Application = &Application{}

// Configure configures the behavior of the engine as it relates to this
// application.
//
// c provides access to the various configuration options, such as specifying
// which message handlers the application contains.
//
// If a.ConfigureFunc is non-nil, it calls a.ConfigureFunc(c).
func (a *Application) Configure(c dogma.ApplicationConfigurer) {
	if a.ConfigureFunc != nil {
		a.ConfigureFunc(c)
	}
}
