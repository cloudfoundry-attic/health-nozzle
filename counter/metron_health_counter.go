package counter

func NewMetronHealthCounter() *HealthCounter {
	return &HealthCounter{
		Ingress: newCounter(),
		Egress:  newCounter(),
		Dropped: newCounter(),
		metricMapping: map[string]string{
			"dropsondeAgentListener.receivedMessageCount": "in",
			"grpc.sentMessageCount":                       "out",
			"udp.sentMessageCount":                        "out",
			"udp.sendErrorCount":                          "dropped",
			"grpc.sendErrorCount":                         "dropped",
			// this metric has no way to determine if it was grpc or udp
			// "DopplerForwarder.sentMessages":               "out",
		},
		name:   "metron",
		origin: "MetronAgent",
	}
}
