package counter_test

import (
	"github.com/cloudfoundry/sonde-go/events"
	"github.com/gogo/protobuf/proto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry-incubator/health-nozzle/counter"
)

var _ = Describe("Metron Health Counter", func() {
	Context("When we receive a counter metric with ingress", func() {
		createEnvelope := func(index, metricName string, total uint64, protocol, origin string) *events.Envelope {
			return &events.Envelope{
				Origin:     proto.String(origin),
				Deployment: proto.String("dontcare"),
				Job:        proto.String("dontcare"),
				Index:      proto.String(index),
				EventType:  events.Envelope_CounterEvent.Enum(),
				Tags: map[string]string{
					"protocol": protocol,
				},
				CounterEvent: &events.CounterEvent{
					Name:  proto.String(metricName),
					Delta: proto.Uint64(0),
					Total: proto.Uint64(total),
				},
			}
		}

		Context("with dropped metrics", func() {
			createDroppedUDPEnvelope := func(index string, total uint64) *events.Envelope {
				return createEnvelope(index, "udp.sendErrorCount", total, "", "MetronAgent")
			}
			createDroppedGRPCEnvelope := func(index string, total uint64) *events.Envelope {
				return createEnvelope(index, "grpc.sendErrorCount", total, "", "MetronAgent")
			}

			It("increases the Dropped count with UPD and GRPC", func() {
				first := createDroppedUDPEnvelope("31f39565-8707-4c8d-aa1c-88cf0b06ac06", 24)
				second := createDroppedUDPEnvelope("31f39565-8707-4c8d-aa1c-88cf0b06ac06", 34)
				third := createDroppedGRPCEnvelope("31f39565-8707-4c8d-aa1c-88cf0b06ac06", 44)
				forth := createDroppedGRPCEnvelope("31f39565-8707-4c8d-aa1c-88cf0b06ac06", 54)

				counter := counter.NewMetronHealthCounter()
				counter.Process(first)
				counter.Process(second)

				Expect(counter.Ingress.Value()).To(Equal(uint64(0)))
				Expect(counter.Egress.Value()).To(Equal(uint64(0)))
				Expect(counter.Dropped.Value()).To(Equal(uint64(10)))

				counter.Process(third)
				counter.Process(forth)

				Expect(counter.Ingress.Value()).To(Equal(uint64(0)))
				Expect(counter.Egress.Value()).To(Equal(uint64(0)))
				Expect(counter.Dropped.Value()).To(Equal(uint64(20)))
			})
		})

		Context("with egress metrics", func() {
			createEgressGRPCEnvelope := func(index string, total uint64) *events.Envelope {
				return createEnvelope(index, "grpc.sentMessageCount", total, "", "MetronAgent")
			}

			createEgressUDPEnvelope := func(index string, total uint64) *events.Envelope {
				return createEnvelope(index, "udp.sentMessageCount", total, "", "MetronAgent")
			}

			It("increases the Egress count with UDP", func() {
				first := createEgressUDPEnvelope("31f39565-8707-4c8d-aa1c-88cf0b06ac06", 24)
				second := createEgressUDPEnvelope("31f39565-8707-4c8d-aa1c-88cf0b06ac06", 34)

				counter := counter.NewMetronHealthCounter()
				counter.Process(first)
				counter.Process(second)

				Expect(counter.Ingress.Value()).To(Equal(uint64(0)))
				Expect(counter.Egress.Value()).To(Equal(uint64(10)))
				Expect(counter.Dropped.Value()).To(Equal(uint64(0)))
			})

			It("increases the Egress count with GRPC", func() {
				first := createEgressGRPCEnvelope("31f39565-8707-4c8d-aa1c-88cf0b06ac06", 24)
				second := createEgressGRPCEnvelope("31f39565-8707-4c8d-aa1c-88cf0b06ac06", 34)

				counter := counter.NewMetronHealthCounter()
				counter.Process(first)
				counter.Process(second)

				Expect(counter.Ingress.Value()).To(Equal(uint64(0)))
				Expect(counter.Egress.Value()).To(Equal(uint64(10)))
				Expect(counter.Dropped.Value()).To(Equal(uint64(0)))
			})

		})
		Context("with ingress metrics", func() {
			createIngressEnvelope := func(index string, total uint64) *events.Envelope {
				return createEnvelope(index, "dropsondeAgentListener.receivedMessageCount", total, "dontcare", "MetronAgent")
			}

			It("increases the Ingress count", func() {
				first := createIngressEnvelope("31f39565-8707-4c8d-aa1c-88cf0b06ac06", 24)
				second := createIngressEnvelope("31f39565-8707-4c8d-aa1c-88cf0b06ac06", 34)
				third := createIngressEnvelope("dontcare", 54)
				fourth := createIngressEnvelope("dontcare", 55)
				fifth := createIngressEnvelope("31f39565-8707-4c8d-aa1c-88cf0b06ac06", 44)

				counter := counter.NewMetronHealthCounter()
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
		})

		Context("when metric comes from different origin", func() {
			createIngressEnvelope := func(index string, total uint64) *events.Envelope {
				return createEnvelope(index, "dropsondeAgentListener.receivedMessageCount", total, "dontcare", "DopplerServer")
			}

			It("does not increase any counters", func() {
				first := createIngressEnvelope("dontcare", 1)
				second := createIngressEnvelope("dontcare", 11)

				counter := counter.NewMetronHealthCounter()

				counter.Process(first)
				counter.Process(second)

				Expect(counter.Ingress.Value()).To(Equal(uint64(0)))
				Expect(counter.Egress.Value()).To(Equal(uint64(0)))
				Expect(counter.Dropped.Value()).To(Equal(uint64(0)))
			})
		})
	})

	Context("When we receive a metric for something else", func() {
		It("does not increase any counters", func() {
			envelope := &events.Envelope{}

			counter := counter.NewMetronHealthCounter()
			counter.Process(envelope)
			Expect(counter.Ingress.Value()).To(Equal(uint64(0)))
			Expect(counter.Egress.Value()).To(Equal(uint64(0)))
			Expect(counter.Dropped.Value()).To(Equal(uint64(0)))
		})
	})

})
