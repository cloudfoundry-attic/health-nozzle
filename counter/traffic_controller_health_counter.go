package counter

func NewTCHealthCounter() *HealthCounter {
	return &HealthCounter{
		Ingress: newCounter(),
		Egress:  newCounter(),
		Dropped: newCounter(),
		metricMapping: map[string]string{
			"listeners.receivedEnvelopes": "in",
		},
		name:   "traffic_controller",
		origin: "LoggregatorTrafficController",
	}
}
