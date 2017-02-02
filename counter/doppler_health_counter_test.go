package counter_test

import (
	"github.com/cloudfoundry/sonde-go/events"
	"github.com/gogo/protobuf/proto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-incubator/health-nozzle/counter"
)

var _ = Describe("Doppler Health Counter", func() {
	Context("When we receive a counter metric with ingress", func() {
		createEnvelope := func(index, metricName string, total uint64, protocol, eventType string) *events.Envelope {
			return &events.Envelope{
				Origin:     proto.String("DopplerServer"),
				Deployment: proto.String("loggregator"),
				Job:        proto.String("doppler"),
				Index:      proto.String(index),
				EventType:  events.Envelope_CounterEvent.Enum(),
				Tags: map[string]string{
					"protocol":   protocol,
					"event_type": eventType,
				},
				CounterEvent: &events.CounterEvent{
					Name:  proto.String(metricName),
					Delta: proto.Uint64(0),
					Total: proto.Uint64(total),
				},
			}
		}

		Context("with dropped metrics", func() {
			createDroppedGRPCEnvelope := func(index string, total uint64) *events.Envelope {
				return createEnvelope(index, "doppler.shedEnvelopes", total, "", "")
			}

			It("increases the Dropped value", func() {
				first := createDroppedGRPCEnvelope("31f39565-8707-4c8d-aa1c-88cf0b06ac06", 44)
				second := createDroppedGRPCEnvelope("31f39565-8707-4c8d-aa1c-88cf0b06ac06", 54)

				counter := counter.NewDopplerHealthCounter()
				counter.Process(first)

				Expect(counter.Ingress.Value()).To(Equal(uint64(0)))
				Expect(counter.Egress.Value()).To(Equal(uint64(0)))
				Expect(counter.Dropped.Value()).To(Equal(uint64(0)))

				counter.Process(second)

				Expect(counter.Ingress.Value()).To(Equal(uint64(0)))
				Expect(counter.Egress.Value()).To(Equal(uint64(0)))
				Expect(counter.Dropped.Value()).To(Equal(uint64(10)))
			})
		})

		Context("with egress metrics", func() {
			createEgressEnvelope := func(index string, total uint64) *events.Envelope {
				return createEnvelope(index, "httpServer.receivedMessages", total, "", "")
			}

			It("increases the Egress value", func() {
				first := createEgressEnvelope("31f39565-8707-4c8d-aa1c-88cf0b06ac06", 24)
				second := createEgressEnvelope("31f39565-8707-4c8d-aa1c-88cf0b06ac06", 34)

				counter := counter.NewDopplerHealthCounter()
				counter.Process(first)

				Expect(counter.Ingress.Value()).To(Equal(uint64(0)))
				Expect(counter.Egress.Value()).To(Equal(uint64(0)))
				Expect(counter.Dropped.Value()).To(Equal(uint64(0)))

				counter.Process(second)

				Expect(counter.Ingress.Value()).To(Equal(uint64(0)))
				Expect(counter.Egress.Value()).To(Equal(uint64(10)))
				Expect(counter.Dropped.Value()).To(Equal(uint64(0)))
			})
		})

		Context("with ingress metrics", func() {
			createIngressGRPCEnvelope := func(index string, total uint64, eventType string) *events.Envelope {
				return createEnvelope(index, "listeners.receivedEnvelopes", total, "grpc", eventType)
			}

			createIngressUDPEnvelope := func(index string, total uint64, eventType string) *events.Envelope {
				return createEnvelope(index, "udp.receivedMessageCount", total, "", eventType)
			}

			It("increases the Ingress value", func() {
				first := createIngressGRPCEnvelope("31f39565-8707-4c8d-aa1c-88cf0b06ac06", 24, "")
				second := createIngressGRPCEnvelope("31f39565-8707-4c8d-aa1c-88cf0b06ac06", 34, "")
				third := createIngressUDPEnvelope("dontcare", 54, "")
				fourth := createIngressUDPEnvelope("dontcare", 55, "")
				fifth := createIngressGRPCEnvelope("31f39565-8707-4c8d-aa1c-88cf0b06ac06", 44, "")

				counter := counter.NewDopplerHealthCounter()
				counter.Process(first)

				Expect(counter.Ingress.Value()).To(Equal(uint64(0)))
				Expect(counter.Egress.Value()).To(Equal(uint64(0)))
				Expect(counter.Dropped.Value()).To(Equal(uint64(0)))

				counter.Process(second)

				Expect(counter.Ingress.Value()).To(Equal(uint64(10)))
				Expect(counter.Egress.Value()).To(Equal(uint64(0)))
				Expect(counter.Dropped.Value()).To(Equal(uint64(0)))

				counter.Process(third)

				Expect(counter.Ingress.Value()).To(Equal(uint64(10)))
				Expect(counter.Egress.Value()).To(Equal(uint64(0)))
				Expect(counter.Dropped.Value()).To(Equal(uint64(0)))

				counter.Process(fourth)

				Expect(counter.Ingress.Value()).To(Equal(uint64(11)))
				Expect(counter.Egress.Value()).To(Equal(uint64(0)))
				Expect(counter.Dropped.Value()).To(Equal(uint64(0)))

				counter.Process(fifth)

				Expect(counter.Ingress.Value()).To(Equal(uint64(21)))
				Expect(counter.Egress.Value()).To(Equal(uint64(0)))
				Expect(counter.Dropped.Value()).To(Equal(uint64(0)))
			})

			It("increases the Ingress value when envelopes have different event_type tags", func() {
				first := createIngressGRPCEnvelope("31f39565-8707-4c8d-aa1c-88cf0b06ac06", 24, "HttpStartStop")
				second := createIngressGRPCEnvelope("31f39565-8707-4c8d-aa1c-88cf0b06ac06", 34, "ValueMetric")
				third := createIngressGRPCEnvelope("31f39565-8707-4c8d-aa1c-88cf0b06ac06", 44, "HttpStartStop")

				counter := counter.NewDopplerHealthCounter()
				counter.Process(first)

				Expect(counter.Ingress.Value()).To(Equal(uint64(0)))
				Expect(counter.Egress.Value()).To(Equal(uint64(0)))
				Expect(counter.Dropped.Value()).To(Equal(uint64(0)))

				counter.Process(second)

				Expect(counter.Ingress.Value()).To(Equal(uint64(0)))
				Expect(counter.Egress.Value()).To(Equal(uint64(0)))
				Expect(counter.Dropped.Value()).To(Equal(uint64(0)))

				counter.Process(third)

				Expect(counter.Ingress.Value()).To(Equal(uint64(20)))
				Expect(counter.Egress.Value()).To(Equal(uint64(0)))
				Expect(counter.Dropped.Value()).To(Equal(uint64(0)))
			})
		})
	})

	Context("When we receive a metric for something else", func() {
		It("does not increase any counters", func() {
			envelope := &events.Envelope{}

			counter := counter.NewDopplerHealthCounter()
			counter.Process(envelope)
			Expect(counter.Ingress.Value()).To(Equal(uint64(0)))
			Expect(counter.Egress.Value()).To(Equal(uint64(0)))
			Expect(counter.Dropped.Value()).To(Equal(uint64(0)))
		})
	})
})
