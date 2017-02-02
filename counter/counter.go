package counter

import (
	"strings"
	"sync"

	"github.com/cloudfoundry/sonde-go/events"
)

type mappingKey struct {
	protocol  string
	index     string
	eventType string
}

type mappingTotal map[mappingKey]uint64

type counter struct {
	sync.RWMutex

	value uint64
	base  mappingTotal
	total mappingTotal
}

func (c *counter) Value() uint64 {
	c.RLock()
	defer c.RUnlock()
	return c.value
}

func (c *counter) Add(envelope *events.Envelope) {

	index := envelope.GetIndex()
	protocol := envelope.GetTags()["protocol"]
	eventType := envelope.GetTags()["event_type"]
	total := envelope.CounterEvent.GetTotal()
	name := envelope.CounterEvent.GetName()

	if protocol == "" {
		seperated := strings.Split(name, ".")
		protocol = seperated[0]
	}

	key := mappingKey{protocol, index, eventType}

	c.Lock()
	if c.base[key] == 0 {
		c.base[key] = total
	} else {
		oldTotal := c.total[key]
		c.total[key] = total - c.base[key]
		c.value = c.value + c.total[key] - oldTotal
	}
	c.Unlock()
}

func newCounter() counter {
	return counter{
		base:  make(map[mappingKey]uint64),
		total: make(map[mappingKey]uint64),
	}
}
