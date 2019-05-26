package main

import (
	"github.com/state-alchemists/ayanami/msgbroker"
	"github.com/state-alchemists/ayanami/service"
	"log"
)

const serviceName = "html"

func main() {
	// define broker
	broker, err := msgbroker.NewNats()
	if err != nil {
		log.Fatal(err)
	}
	// define services
	services := service.Services{
		"pre": service.NewService(serviceName, "pre",
			[]string{"input"},
			[]string{"output"},
			WrappedPre,
		),
	}
	// consume and publish forever
	ch := make(chan bool)
	services.ConsumeAndPublish(broker)
	<-ch
}
