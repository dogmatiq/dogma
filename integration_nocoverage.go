package dogma

func (HandlesCommandRoute) isIntegrationRoute() {}
func (HandlesEventRoute) isIntegrationRoute()   {}
func (RecordsEventRoute) isIntegrationRoute()   {}
