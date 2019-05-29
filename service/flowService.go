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
	Value       interface{} // if InputEvent == "", we will use the value instead
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
		for _, flow := range flows {
			if flow.VarName == inputName && flow.InputEvent != "" {
				eventName := flow.InputEvent
				serviceInputs = append(serviceInputs, IO{VarName: inputName, EventName: eventName})
			}
		}
	}
	// populate outputs
	var serviceOutputs []IO
	for _, outputName := range outputs {
		eventName := fmt.Sprintf("flow.%s.out.%s", flowName, outputName)
		serviceOutputs = append(serviceOutputs, IO{VarName: outputName, EventName: eventName})
		for _, flow := range flows {
			if flow.VarName == outputName && flow.OutputEvent != "" {
				eventName := flow.OutputEvent
				serviceOutputs = append(serviceOutputs, IO{VarName: outputName, EventName: eventName})
			}
		}
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
		outputs, vars, ID, err := initFlowWrapper(inputs)
		if err != nil {
			return outputs, err
		}
		// subscribe to every eventFlow's inputEvent and publish to it's output event
		completed := make(chan bool)
		for _, flow := range flows {
			inputEventName := fmt.Sprintf("%s.%s", ID, flow.InputEvent)
			// varName already exists in vars
			if vars.Has(flow.VarName) && flow.OutputEvent != "" {
				publishFlowPkg(broker, flow.OutputEvent, ID, vars.Get(flow.VarName), completed)
			}
			// flow has inputEvent
			if flow.InputEvent != "" {
				log.Printf("[INFO] consuming from %s", inputEventName)
				broker.Consume(inputEventName,
					func(pkg servicedata.Package) {
						// get the message and populate vars based on received message
						log.Printf("[INFO] Get message from `%s`: %#v", inputEventName, pkg)
						vars.Set(flow.VarName, pkg.Data)
						// publish the servicedata
						if flow.OutputEvent != "" {
							publishFlowPkg(broker, flow.OutputEvent, pkg.ID, pkg.Data, completed)
						}
						processCompletedOutput(outputVarNames, vars, completed)
					},
					// error callback
					func(flowErr error) {
						log.Printf("[ERROR] Error while consuming from %s: %s", inputEventName, flowErr)
						completed <- true
					},
				)
			}
		}
		// set all predefined variables
		vars = getFlowDefaultVars(vars, broker, ID, flows, completed)
		// just in case all default value is already provided
		processCompletedOutput(outputVarNames, vars, completed)
		<-completed
		for _, outputName := range outputVarNames {
			outputs.Set(outputName, vars.Get(outputName))
		}
		return outputs, err
	}
}

func processCompletedOutput(outputVarNames []string, vars Dictionary, completed chan bool) {
	if vars.HasAll(outputVarNames) {
		completed <- true
	}
}

func getFlowDefaultVars(vars Dictionary, broker msgbroker.CommonBroker, ID string, flows []FlowEvent, completed chan bool) Dictionary {
	for _, flow := range flows {
		// flow doesn't have inputEvent, use it's value to populate vars
		if flow.InputEvent == "" {
			vars.Set(flow.VarName, flow.Value)
			log.Printf("[INFO] Set `%s` into `%#v`", flow.VarName, flow.Value)
			// flow also has outputEvent, publish the value
			if flow.OutputEvent != "" {
				publishFlowPkg(broker, flow.OutputEvent, ID, flow.Value, completed)
			}
		}
	}
	return vars
}

func initFlowWrapper(inputs Dictionary) (Dictionary, Dictionary, string, error) {
	vars := make(Dictionary)
	outputs := make(Dictionary)
	// create ID
	ID, err := CreateID()
	// preset vars
	for inputName, inputVal := range inputs {
		vars.Set(inputName, inputVal)
	}
	return outputs, vars, ID, err
}

func publishFlowPkg(broker msgbroker.CommonBroker, rawOutputEventName, ID string, data interface{}, completed chan bool) {
	pkg := servicedata.Package{ID: ID, Data: data}
	outputEventName := fmt.Sprintf("%s.%s", ID, rawOutputEventName)
	log.Printf("[INFO] Publish into `%s`: `%#v`", outputEventName, pkg)
	err := broker.Publish(outputEventName, pkg)
	if err != nil {
		log.Printf("[ERROR] Error: %s", err)
		completed <- true
	}
}
