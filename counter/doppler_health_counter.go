package counter

func NewDopplerHealthCounter() *HealthCounter {
	return &HealthCounter{
		Ingress: newCounter(),
		Egress:  newCounter(),
		Dropped: newCounter(),
		metricMapping: map[string]string{
			"doppler.shedEnvelopes":       "dropped",
			"httpServer.receivedMessages": "out",
			"listeners.receivedEnvelopes": "in",
			"udp.receivedMessageCount":    "in",
		},
		name:   "doppler",
		origin: "DopplerServer",
	}
}
