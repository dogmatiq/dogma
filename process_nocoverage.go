package dogma

func (HandlesEventRoute) isProcessRoute()     {}
func (ExecutesCommandRoute) isProcessRoute()  {}
func (SchedulesTimeoutRoute) isProcessRoute() {}
