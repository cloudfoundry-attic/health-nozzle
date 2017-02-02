package counter

import "github.com/cloudfoundry/sonde-go/events"

type HealthCounter struct {
	Ingress counter
	Egress  counter
	Dropped counter

	metricMapping map[string]string
	name          string
	origin        string
}

func (c *HealthCounter) GetName() string {
	return c.name
}

func (c *HealthCounter) Process(envelope *events.Envelope) {
	if envelope.GetEventType() != events.Envelope_CounterEvent {
		return
	}

	if envelope.GetOrigin() != c.origin {
		return
	}

	switch c.metricMapping[envelope.CounterEvent.GetName()] {
	case "in":
		c.Ingress.Add(envelope)
	case "out":
		c.Egress.Add(envelope)
	case "dropped":
		c.Dropped.Add(envelope)
	}
}
