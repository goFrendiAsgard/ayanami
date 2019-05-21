package main

import (
	"fmt"
	nats "github.com/nats-io/nats.go"
	"log"
)

// SrvcNewFlowConfig create new singleConfig for flow
func SrvcNewFlowConfig(flowConfig SrvcSingleFlowConfig) SrvcSingleConfig {
	var singleConfig SrvcSingleConfig
	// populate inputs
	singleConfig.Input = flowConfig.Input
	inputVarNames := SrvcGetUniqueVarNames(flowConfig.Input)
	for _, varName := range inputVarNames {
		eventName := fmt.Sprintf("flow.%s.in.%s", flowConfig.FlowName, varName)
		singleConfig.Input = append(singleConfig.Input, SrvcServiceIO{VarName: varName, EventName: eventName})
	}
	// populate outputs
	singleConfig.Output = flowConfig.Output
	outputVarNames := SrvcGetUniqueVarNames(flowConfig.Output)
	for _, varName := range outputVarNames {
		eventName := fmt.Sprintf("flow.%s.in.%s", flowConfig.FlowName, varName)
		singleConfig.Output = append(singleConfig.Output, SrvcServiceIO{VarName: varName, EventName: eventName})
	}
	singleConfig.Function = createFlowWrapper(flowConfig.Flows)
	return singleConfig
}

func createFlowWrapper(flows []SrvcEventFlow) SrvcWrappedFunction {
	natsURL := SrvcGetNatsURL()
	return func(inputs SrvcDictionary) SrvcDictionary {
		var outputs SrvcDictionary
		// create ID
		ID, err := SrvcCreateID()
		if err != nil {
			log.Print(err)
			return outputs
		}
		// connect to nats
		nc, err := nats.Connect(natsURL)
		if err != nil {
			log.Print(err)
			return outputs
		}
		// add values to every eventFlow without inputEvent
		for index, flow := range flows {
			if flow.InputEvent == "" {
				flows[index].Value = inputs[flow.VarName]
			}
		}
		// subscribe to every eventFlow that has inputEvent
		for _, flow := range flows {
			if flow.InputEvent != "" {
				eventName := fmt.Sprintf("%s.%s", ID, flow.InputEvent)
				nc.Subscribe(eventName, func(m *nats.Msg) {
					// TODO: extract, add to outputs
					// TODO: if outputs is completed, return
				})
			}
		}
		// TODO: publish every eventFlow without inputName but with outputEvent
		return outputs
	}
}
