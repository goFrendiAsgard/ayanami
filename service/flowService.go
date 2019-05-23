package service

import (
	"encoding/json"
	"fmt"
	nats "github.com/nats-io/nats.go"
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

// SrvcNewFlowConfig create new singleConfig for flow
func SrvcNewFlowConfig(flowConfig FlowService) CommonService {
	var singleConfig CommonService
	// populate inputs
	singleConfig.Input = flowConfig.Input
	inputVarNames := GetUniqueVarNames(flowConfig.Input)
	for _, varName := range inputVarNames {
		eventName := fmt.Sprintf("flow.%s.in.%s", flowConfig.FlowName, varName)
		singleConfig.Input = append(singleConfig.Input, IO{VarName: varName, EventName: eventName})
	}
	// populate outputs
	singleConfig.Output = flowConfig.Output
	outputVarNames := GetUniqueVarNames(flowConfig.Output)
	for _, varName := range outputVarNames {
		eventName := fmt.Sprintf("flow.%s.in.%s", flowConfig.FlowName, varName)
		singleConfig.Output = append(singleConfig.Output, IO{VarName: varName, EventName: eventName})
	}
	singleConfig.Function = createFlowWrapper(flowConfig.Flows, outputVarNames)
	return singleConfig
}

func createFlowWrapper(flows []FlowEvent, outputVarNames []string) WrappedFunction {
	natsURL := GetNatsURL()
	return func(inputs Dictionary) Dictionary {
		var outputs Dictionary
		// create ID
		ID, err := CreateID()
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
		completed := make(chan bool)
		for _, flow := range flows {
			if flow.InputEvent != "" {
				inputEventName := fmt.Sprintf("%s.%s", ID, flow.InputEvent)
				outputEventName := fmt.Sprintf("%s.%s", ID, flow.OutputEvent)
				varName := flow.VarName
				nc.Subscribe(inputEventName, func(m *nats.Msg) {
					// get the message and populate outputs based on received message
					var pkg Package
					JSONByte := m.Data
					log.Printf("[INFO] Get message from `%s`: %s", inputEventName, string(JSONByte))
					err := json.Unmarshal(JSONByte, &pkg)
					if err != nil {
						log.Printf("[ERROR] %s: %s", inputEventName, err)
						return
					}
					outputs[varName] = pkg.Data
					// publish the data
					if outputEventName != "" {
						pkg.Data = outputs[varName]
						JSONByte, err := json.Marshal(&pkg)
						if err != nil {
							log.Printf("[ERROR] %s: %s", ID, err)
							return
						}
						nc.Publish(outputEventName, JSONByte)
						log.Printf("[INFO] Publish into `%s`: `%#v`", outputEventName, pkg)
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
					pkg := Package{ID: ID, Data: value}
					JSONByte, err := json.Marshal(&pkg)
					if err != nil {
						log.Printf("[ERROR] %s: %s", ID, err)
					}
					nc.Publish(outputEventName, JSONByte)
					log.Printf("[INFO] Publish into `%s`: `%#v`", outputEventName, pkg)
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
