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
	Input    IOList
	Output   IOList
	Flows    []FlowEvent
}

// NewFlow create new service for flow
func NewFlow(broker msgbroker.CommonBroker, flowName string, inputs, outputs []string, flows []FlowEvent) CommonService {
	// populate inputs
	var serviceInputs []IO
	for _, inputName := range inputs {
		eventName := fmt.Sprintf("flow.%s.in.%s", flowName, inputName)
		serviceInputs = append(serviceInputs, IO{VarName: inputName, EventName: eventName})
	}
	// populate outputs
	var serviceOutputs []IO
	for _, outputName := range outputs {
		eventName := fmt.Sprintf("flow.%s.out.%s", flowName, outputName)
		serviceOutputs = append(serviceOutputs, IO{VarName: outputName, EventName: eventName})
	}
	// get errorEventName
	errorEventName := fmt.Sprintf("flow.%s.err", flowName)
	// get flowWrappedFunction
	wrappedFunction := createFlowWrapper(broker, flows, outputs)
	return CommonService{
		Input:          serviceInputs,
		Output:         serviceOutputs,
		ErrorEventName: errorEventName,
		Function:       wrappedFunction,
	}
}

func createFlowWrapper(broker msgbroker.CommonBroker, flows []FlowEvent, outputVarNames []string) WrappedFunction {
	return func(inputs Dictionary) (Dictionary, error) {
		var outputs Dictionary
		// create ID
		ID, err := CreateID()
		if err != nil {
			log.Print(err)
			return outputs, err
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
				broker.Consume(inputEventName,
					func(pkg servicedata.Package) {
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
					},
					// error callback
					func(flowErr error) {
						err = flowErr
						log.Printf("[ERROR] Error: %s", err)
					},
				)
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
		return outputs, err
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
