package main

import (
	"github.com/state-alchemists/ayanami/msgbroker"
	"github.com/state-alchemists/ayanami/service"
	"log"
)

const serviceName = "cmd"

func main() {
	// define broker
	broker, err := msgbroker.NewNats()
	if err != nil {
		log.Fatal(err)
	}
	// define services
	services := service.Services{
		"cowsay": service.NewCmd(
			serviceName,
			"cowsay",
			[]string{"input"},
			[]string{"output"},
			[]string{"cowsay", "-n", "$input"},
		),
		"figlet": service.NewCmd(
			serviceName,
			"figlet",
			[]string{"input"},
			[]string{"output"},
			[]string{"figlet", "$input"},
		),
	}
	// consume and publish forever
	ch := make(chan bool)
	services.ConsumeAndPublish(broker)
	<-ch
}
