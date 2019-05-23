package main

import (
	"github.com/state-alchemists/ayanami/msgbroker"
	"github.com/state-alchemists/ayanami/service"
	"log"
)

const serviceName = "textDecorator"

func main() {
	// define broker
	broker, err := msgbroker.NewNats()
	if err != nil {
		log.Fatal(err)
	}
	// define services
	services := service.Services{
		"pre": service.NewService(
			serviceName,
			"pre",
			[]string{"text"},
			[]string{"text"},
			WrappedPre,
		),
		"cowsay": service.NewService(
			serviceName,
			"cowsay",
			[]string{"text"},
			[]string{"text"},
			WrappedCowsay,
		),
		"figlet": service.NewService(
			serviceName,
			"figlet",
			[]string{"text"},
			[]string{"text"},
			WrappedFiglet,
		),
	}
	// consume and publish forever
	ch := make(chan bool)
	services.ConsumeAndPublish(broker)
	<-ch
}
