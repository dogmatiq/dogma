package dogma

func (ViaAggregateRoute) isHandlerRoute()   {}
func (ViaProcessRoute) isHandlerRoute()     {}
func (ViaIntegrationRoute) isHandlerRoute() {}
func (ViaProjectionRoute) isHandlerRoute()  {}

func (WithAggregateSpec) isHandlerSpec()   {}
func (WithProcessSpec) isHandlerSpec()     {}
func (WithIntegrationSpec) isHandlerSpec() {}
func (WithProjectionSpec) isHandlerSpec()  {}
