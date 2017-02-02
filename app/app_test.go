package app_test

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"

	"github.com/cloudfoundry-incubator/health-nozzle/app"
	"github.com/cloudfoundry/sonde-go/events"
	"github.com/gogo/protobuf/proto"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("App", func() {
	createEnvelope := func(index, metricName string, total uint64, protocol string) *events.Envelope {
		return &events.Envelope{
			Origin:     proto.String("MetronAgent"),
			Deployment: proto.String("loggregator"),
			Job:        proto.String("metron"),
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

	createIngressEnvelope := func(index string, total uint64) *events.Envelope {
		return createEnvelope(index, "dropsondeAgentListener.receivedMessageCount", total, "dontcare")
	}

	getTCPPort := func() string {
		l, err := net.Listen("tcp", ":0")
		if err != nil {
			panic(err)
		}
		defer l.Close()
		return fmt.Sprintf("%d", l.Addr().(*net.TCPAddr).Port)
	}

	Context("starting the app", func() {
		It("writes incoming messages to the counters", func() {
			port := getTCPPort()
			msgChan := make(chan *events.Envelope)

			app := app.NewApp(msgChan)
			err := app.Start(port)
			Expect(err).ToNot(HaveOccurred())

			Expect(app.Counters[0].Ingress.Value()).To(Equal(uint64(0)))

			msgChan <- createIngressEnvelope("unique-index", 10)
			msgChan <- createIngressEnvelope("unique-index", 20)

			Eventually(app.Counters[0].Ingress.Value).Should(Equal(uint64(10)))
			app.Stop()
		})

		It("starts the HTTP endpoint", func() {
			port := getTCPPort()
			msgChan := make(chan *events.Envelope)

			app := app.NewApp(msgChan)
			err := app.Start(port)
			Expect(err).ToNot(HaveOccurred())

			response, err := http.Get(fmt.Sprintf("http://localhost:%s", port))
			Expect(err).ToNot(HaveOccurred())

			defer response.Body.Close()
			body, err := ioutil.ReadAll(response.Body)
			Expect(err).ToNot(HaveOccurred())

			Expect(body).To(MatchJSON(`
			{
				"doppler": {
				  "ingress": 0,
				  "egress": 0,
				  "dropped": 0
				},
				"metron": {
				  "ingress": 0,
				  "egress": 0,
				  "dropped": 0
				},
				"traffic_controller": {
				  "ingress": 0,
				  "egress": 0,
				  "dropped": 0
				}
			}
			`))
			app.Stop()
		})
	})
})
