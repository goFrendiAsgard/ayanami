package service

import (
	"fmt"
	"github.com/state-alchemists/ayanami/msgbroker"
	"github.com/state-alchemists/ayanami/servicedata"
	"log"
)

// FlowEvent event flow
type FlowEvent struct {
	InputEvent  string
	OutputEvent string
	VarName     string      // read from inputEvent, put into var, publish into outputEvent
	Value       interface{} // if InputEvent == "", the Value will be published instead
}

// FlowService single flow config
type FlowService struct {
	FlowName string
	Input    []IO
	Output   []IO
	Flows    []FlowEvent
}

// NewFlow create new service for flow
func NewFlow(broker msgbroker.CommonBroker, flowConfig FlowService) CommonService {
	var service CommonService
	// populate inputs
	service.Input = flowConfig.Input
	inputVarNames := GetUniqueVarNames(flowConfig.Input)
	for _, varName := range inputVarNames {
		eventName := fmt.Sprintf("flow.%s.in.%s", flowConfig.FlowName, varName)
		service.Input = append(service.Input, IO{VarName: varName, EventName: eventName})
	}
	// populate outputs
	service.Output = flowConfig.Output
	outputVarNames := GetUniqueVarNames(flowConfig.Output)
	for _, varName := range outputVarNames {
		eventName := fmt.Sprintf("flow.%s.in.%s", flowConfig.FlowName, varName)
		service.Output = append(service.Output, IO{VarName: varName, EventName: eventName})
	}
	service.Function = createFlowWrapper(broker, flowConfig.Flows, outputVarNames)
	return service
}

func createFlowWrapper(broker msgbroker.CommonBroker, flows []FlowEvent, outputVarNames []string) WrappedFunction {
	return func(inputs Dictionary) Dictionary {
		var outputs Dictionary
		// create ID
		ID, err := CreateID()
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
		completed := make(chan bool)
		for _, flow := range flows {
			if flow.InputEvent != "" {
				inputEventName := fmt.Sprintf("%s.%s", ID, flow.InputEvent)
				outputEventName := fmt.Sprintf("%s.%s", ID, flow.OutputEvent)
				varName := flow.VarName
				broker.Consume(inputEventName, func(pkg servicedata.Package) {
					// get the message and populate outputs based on received message
					log.Printf("[INFO] Get message from `%s`: %#v", inputEventName, pkg)
					outputs[varName] = pkg.Data
					// publish the servicedata
					if outputEventName != "" {
						log.Printf("[INFO] Publish into `%s`: `%#v`", outputEventName, pkg)
						broker.Publish(outputEventName, pkg)
					}
					if isOutputComplete(outputVarNames, outputs) {
						completed <- true
					}
				})
			}
		}
		// set predefined variables
		for _, flow := range flows {
			if flow.InputEvent == "" {
				value := flow.Value
				varName := flow.VarName
				outputs[varName] = value
				log.Printf("[INFO] Set `%s` into `%#v`", varName, value)
				if flow.OutputEvent != "" {
					outputEventName := fmt.Sprintf("%s.%s", ID, flow.OutputEvent)
					pkg := servicedata.Package{ID: ID, Data: value}
					log.Printf("[INFO] Publish into `%s`: `%#v`", outputEventName, pkg)
					broker.Publish(outputEventName, pkg)
				}
			}
		}
		<-completed
		return outputs
	}
}

func isOutputComplete(outputVarNames []string, outputs Dictionary) bool {
	for _, outputVarName := range outputVarNames {
		if _, exists := outputs[outputVarName]; !exists {
			return false
		}
	}
	return true
}
