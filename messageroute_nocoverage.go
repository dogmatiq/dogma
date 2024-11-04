package dogma

func (HandlesCommandRoute) isMessageRoute()   {}
func (ExecutesCommandRoute) isMessageRoute()  {}
func (HandlesEventRoute) isMessageRoute()     {}
func (RecordsEventRoute) isMessageRoute()     {}
func (SchedulesTimeoutRoute) isMessageRoute() {}
