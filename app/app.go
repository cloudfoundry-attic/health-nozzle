package app

import (
	"fmt"
	"net"
	"net/http"

	"github.com/cloudfoundry-incubator/health-nozzle/counter"
	"github.com/cloudfoundry-incubator/health-nozzle/server"
	"github.com/cloudfoundry/sonde-go/events"
)

type app struct {
	messages <-chan *events.Envelope
	server   *server.Server
	listener net.Listener

	Counters []*counter.HealthCounter
}

func NewApp(messages <-chan *events.Envelope) *app {
	counters := []*counter.HealthCounter{
		counter.NewMetronHealthCounter(),
		counter.NewDopplerHealthCounter(),
		counter.NewTCHealthCounter(),
	}
	return &app{
		messages: messages,
		Counters: counters,
		server:   server.NewServer(counters),
	}
}

func (a *app) Start(port string) error {
	go func() {
		for {
			select {
			case msg := <-a.messages:
				for _, counter := range a.Counters {
					counter.Process(msg)
				}
			}
		}
	}()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return err
	}

	a.listener = listener
	go http.Serve(listener, a.server)

	return nil
}

func (a *app) Stop() {
	a.listener.Close()
}
