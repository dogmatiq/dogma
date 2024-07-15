package dogma

func (UnicastProjectionDeliveryPolicy) isProjectionDeliveryPolicy()   {}
func (BroadcastProjectionDeliveryPolicy) isProjectionDeliveryPolicy() {}

func (HandlesEventRoute) isProjectionRoute() {}
