package dogma

func (AggregateRegistration) isHandlerRegistration()   {}
func (ProcessRegistration) isHandlerRegistration()     {}
func (IntegrationRegistration) isHandlerRegistration() {}
func (ProjectionRegistration) isHandlerRegistration()  {}
