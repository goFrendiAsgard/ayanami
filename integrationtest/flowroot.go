package integrationtest

import (
	"github.com/state-alchemists/ayanami/config"
	"github.com/state-alchemists/ayanami/msgbroker"
	"github.com/state-alchemists/ayanami/service"
	"log"
)

func helloWorld(input interface{}) interface{} {
	return "hello world"
}

// MainFlowRoot declaration
func MainFlowRoot() {
	// define broker
	broker, err := msgbroker.NewNats(config.GetNatsURL())
	if err != nil {
		log.Fatal(err)
	}
	// define services
	services := service.Services{
		service.NewFlow("root", broker,
			// inputs
			[]string{"content", "code"},
			// outputs
			[]string{"content", "code"},
			[]service.FlowEvent{
				service.FlowEvent{
					InputEvent:  "trig.request.get.out.req",
					OutputEvent: "trig.response.get.in.code",
					VarName:     "code",
					UseValue:    true,
					Value:       200,
				},
				service.FlowEvent{
					InputEvent:  "trig.request.get.out.req",
					OutputEvent: "trig.response.get.in.content",
					VarName:     "content",
					UseFunction: true,
					Function:    helloWorld,
				},
			},
		),
	}
	// consume and publish forever
	forever := make(chan bool)
	services.ConsumeAndPublish(broker, "flow")
	<-forever
}
