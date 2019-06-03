package integrationtest

import (
	"github.com/state-alchemists/ayanami/config"
	"github.com/state-alchemists/ayanami/msgbroker"
	"github.com/state-alchemists/ayanami/service"
	"log"
)

// MainServiceCmd emulating cmd's main function
func MainServiceCmd() {
	serviceName := "cmd"
	// define broker
	broker, err := msgbroker.NewNats(config.GetNatsURL())
	if err != nil {
		log.Fatal(err)
	}
	// define services
	services := service.Services{
		"cowsay": service.NewCmd(serviceName, "cowsay",
			[]string{"input"},
			[]string{"output"},
			[]string{"/bin/sh", "-c", "echo $input | cowsay -n"},
		),
		"figlet": service.NewCmd(serviceName, "figlet",
			[]string{"input"},
			[]string{"output"},
			[]string{"figlet", "$input"},
		),
	}
	// consume and publish forever
	ch := make(chan bool)
	services.ConsumeAndPublish(broker, serviceName)
	<-ch
}
