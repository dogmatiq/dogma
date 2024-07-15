package dogma

func (HandlesCommandRoute) isRoute()   {}
func (ExecutesCommandRoute) isRoute()  {}
func (HandlesEventRoute) isRoute()     {}
func (RecordsEventRoute) isRoute()     {}
func (SchedulesTimeoutRoute) isRoute() {}
