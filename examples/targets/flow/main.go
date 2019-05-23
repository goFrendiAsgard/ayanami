package main

import (
	"github.com/state-alchemists/ayanami/service"
)

var configs service.CommonService

func init() {
	configs = SrvcConfigs{
		"pre": SrvcNewFlowConfig(SrvcSingleFlowConfig{
			FlowName: "pre",
			Input: []SrvcServiceIO{
				SrvcServiceIO{
					EventName: "trig.request.get./pre.out.text",
					VarName:   "form",
				},
			},
			Output: []SrvcServiceIO{
				SrvcServiceIO{
					EventName: "trig.response.get./pre.in.code",
					VarName:   "code",
				},
				SrvcServiceIO{
					EventName: "trig.request.get./pre.in.content",
					VarName:   "preResult",
				},
			},
			Flows: []SrvcEventFlow{
				SrvcEventFlow{
					VarName: "code",
					Value:   200,
				},
			},
		}),
		"echo": SrvcSingleConfig{
			Input: []SrvcServiceIO{
				SrvcServiceIO{
					EventName: "trig.request.get./echo.out.form",
					VarName:   "form",
				},
			},
			Output: []SrvcServiceIO{
				SrvcServiceIO{
					EventName: "trig.response.get./echo.in.code",
					VarName:   "code",
				},
				SrvcServiceIO{
					EventName: "trig.response.get./echo.in.content",
					VarName:   "content",
				},
			},
			Function: WrappedEcho,
		},
	}
}

func main() {
	// consume and publish forever
	ch := make(chan bool)
	SrvcConsumeAndPublish(configs)
	<-ch
}
