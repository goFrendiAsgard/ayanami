package main

import (
	"github.com/state-alchemists/ayanami/msgbroker"
	"github.com/state-alchemists/ayanami/service"
	"log"
)

func main() {
	// define broker
	broker, err := msgbroker.NewNats()
	if err != nil {
		log.Fatal(err)
	}
	// define services
	services := service.Services{
		"pre": service.NewFlow(broker, service.FlowService{
			FlowName: "pre",
			Input: []service.IO{
				service.IO{
					EventName: "trig.request.get /pre.out.text",
					VarName:   "form",
				},
			},
			Output: []service.IO{
				service.IO{
					EventName: "trig.response.get /pre.in.code",
					VarName:   "code",
				},
				service.IO{
					EventName: "trig.request.get /pre.in.content",
					VarName:   "preResult",
				},
			},
			Flows: []service.FlowEvent{
				service.FlowEvent{
					VarName: "code",
					Value:   200,
				},
			},
		}),
		"echo": service.CommonService{
			Input: []service.IO{
				service.IO{
					EventName: "trig.request.get /echo.out.form",
					VarName:   "form",
				},
			},
			Output: []service.IO{
				service.IO{
					EventName: "trig.response.get /echo.in.code",
					VarName:   "code",
				},
				service.IO{
					EventName: "trig.response.get /echo.in.content",
					VarName:   "content",
				},
			},
			Function: WrappedEcho,
		},
	}
	// consume and publish forever
	ch := make(chan bool)
	services.ConsumeAndPublish(broker)
	<-ch
}
