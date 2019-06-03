package integrationtest

import (
	"github.com/state-alchemists/ayanami/config"
	"github.com/state-alchemists/ayanami/msgbroker"
	"github.com/state-alchemists/ayanami/service"
	"log"
)

// MainFlow emulating flow's main function
func MainFlow() {
	serviceName := "flow"
	// define broker
	broker, err := msgbroker.NewNats(config.GetNatsURL())
	if err != nil {
		log.Fatal(err)
	}
	// define services
	services := service.Services{
		"main": service.NewFlow(serviceName, "main", broker,
			// inputs
			[]string{"request"},
			// outputs
			[]string{"content", "code"},
			[]service.FlowEvent{
				service.FlowEvent{
					InputEvent: "trig.request.get.out.req",
					VarName:    "request",
				},
				service.FlowEvent{
					VarName:     "request.form.text.0",
					OutputEvent: "srvc.cmd.figlet.in.input",
				},
				service.FlowEvent{
					InputEvent:  "srvc.cmd.figlet.out.output",
					VarName:     "figletOutput",
					OutputEvent: "srvc.cmd.cowsay.in.input",
				},
				service.FlowEvent{
					InputEvent:  "srvc.cmd.cowsay.out.output",
					VarName:     "cowsayOutput",
					OutputEvent: "srvc.html.pre.in.input",
				},
				// normal response
				service.FlowEvent{
					InputEvent:  "srvc.html.pre.out.output",
					VarName:     "content",
					OutputEvent: "trig.response.get.in.content",
				},
				service.FlowEvent{
					InputEvent:  "srvc.html.pre.out.output",
					VarName:     "code",
					UseValue:    true,
					Value:       200,
					OutputEvent: "trig.response.get.in.code",
				},
				// error response from figlet
				service.FlowEvent{
					InputEvent: "srvc.cmd.figlet.err.message",
					VarName:    "code",
					UseValue:   true,
					Value:      500,
				},
				service.FlowEvent{
					InputEvent: "srvc.cmd.figlet.err.message",
					VarName:    "content",
				},
				// error response from cowsay
				service.FlowEvent{
					InputEvent: "srvc.cmd.cowsay.err.message",
					VarName:    "code",
					UseValue:   true,
					Value:      500,
				},
				service.FlowEvent{
					InputEvent: "srvc.cmd.cowsay.err.message",
					VarName:    "content",
				},
				// error response from pre
				service.FlowEvent{
					InputEvent: "srvc.html.pre.err.message",
					VarName:    "code",
					UseValue:   true,
					Value:      500,
				},
				service.FlowEvent{
					InputEvent: "srvc.html.pre.err.message",
					VarName:    "content",
				},
			},
		),
	}
	// consume and publish forever
	ch := make(chan bool)
	services.ConsumeAndPublish(broker, "flow")
	<-ch
}
