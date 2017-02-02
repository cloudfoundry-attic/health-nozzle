package main

import (
	"crypto/tls"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/cloudfoundry-incubator/health-nozzle/app"
	"github.com/cloudfoundry/noaa/consumer"
)

const firehoseSubscriptionId = "super-awesome-healthkit"

var (
	dopplerAddress = os.Getenv("DOPPLER_ADDR")
	authToken      = os.Getenv("CF_ACCESS_TOKEN")
	port           = os.Getenv("PORT")
)

func waitForKill() {
	waitFor := make(chan os.Signal, 2)
	signal.Notify(waitFor, os.Interrupt, syscall.SIGTERM)

	waiter := sync.WaitGroup{}
	waiter.Add(1)

	go func() {
		<-waitFor
		waiter.Done()
	}()
	waiter.Wait()
}

func main() {
	consumer := consumer.New(dopplerAddress, &tls.Config{InsecureSkipVerify: true}, nil)
	consumer.SetDebugPrinter(ConsoleDebugPrinter{})

	msgChan, errorChan := consumer.Firehose(firehoseSubscriptionId, authToken)
	go func() {
		for err := range errorChan {
			fmt.Fprintf(os.Stderr, "%v\n", err.Error())
		}
	}()
	app := app.NewApp(msgChan)
	err := app.Start(port)
	if err != nil {
		panic(err)
	}
	waitForKill()
}

type ConsoleDebugPrinter struct{}

func (c ConsoleDebugPrinter) Print(title, dump string) {
	println(title)
	println(dump)
}
