package dogma

func (ViaAggregateRoute) isHandlerRoute()   {}
func (ViaProcessRoute) isHandlerRoute()     {}
func (ViaIntegrationRoute) isHandlerRoute() {}
func (ViaProjectionRoute) isHandlerRoute()  {}