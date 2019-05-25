package service

import (
	"fmt"
	"github.com/state-alchemists/ayanami/msgbroker"
	"github.com/state-alchemists/ayanami/servicedata"
	"log"
)

// WrappedFunction any function wrapped for ConsumeAndPublish
type WrappedFunction = func(inputs Dictionary) Dictionary

// CommonService single configuration
type CommonService struct {
	Input    IOList
	Output   IOList
	Function WrappedFunction
}

// NewService create new singleConfig for service
func NewService(serviceName string, methodName string, inputs []string, outputs []string, wrappedFunction WrappedFunction) CommonService {
	// get inputConfig
	var inputConfig []IO
	for _, inputName := range inputs {
		eventName := fmt.Sprintf("srvc.%s.%s.in.%s", serviceName, methodName, inputName)
		inputConfig = append(inputConfig, IO{VarName: inputName, EventName: eventName})
	}
	// get outputConfig
	var outputConfig []IO
	for _, outputName := range outputs {
		eventName := fmt.Sprintf("srvc.%s.%s.in.%s", serviceName, methodName, outputName)
		outputConfig = append(outputConfig, IO{VarName: outputName, EventName: eventName})
	}
	// return config
	return CommonService{
		Input:    inputConfig,
		Output:   outputConfig,
		Function: wrappedFunction,
	}
}

// Services configuration
type Services map[string]CommonService

// ConsumeAndPublish consume from queue and Publish
func (services Services) ConsumeAndPublish(broker msgbroker.CommonBroker) {
	for _, service := range services {
		inputIOList := service.Input
		outputIOList := service.Output
		wrappedFunction := service.Function
		consumeAndPublishSingle(broker, inputIOList, outputIOList, wrappedFunction)
	}
}

func consumeAndPublishSingle(broker msgbroker.CommonBroker, inputIOList IOList, outputIOList IOList, wrappedFunction WrappedFunction) {
	// allInputs
	allInputs := make(map[string]Dictionary)
	rawInputEventNames := inputIOList.GetUniqueEventNames()
	inputVarNames := inputIOList.GetUniqueVarNames()
	for _, rawInputEventName := range rawInputEventNames {
		inputEventName := fmt.Sprintf("*.%s", rawInputEventName)
		log.Printf("[INFO] Consume from `%s`", inputEventName)
		broker.Consume(inputEventName,
			// success callback
			func(pkg servicedata.Package) {
				// prepare allInputs
				ID := pkg.ID
				data := pkg.Data
				if _, exists := allInputs[ID]; !exists {
					allInputs[ID] = Dictionary{}
				}
				// populate allInputs[ID] with eventInputNames and servicedata
				eventInputNames := inputIOList.GetEventVarNames(rawInputEventName)
				for _, inputVarName := range eventInputNames {
					allInputs[ID][inputVarName] = data
				}
				inputs := allInputs[ID]
				log.Printf("[INFO] Inputs for %s: %#v", ID, inputs)
				// execute wrapper
				if isInputComplete(inputVarNames, inputs) {
					log.Printf("[INFO] Inputs for %s completed", ID)
					outputs := wrappedFunction(inputs)
					publish(broker, ID, outputIOList, outputs)
				}
			},
			// error callback
			func(err error) {
				log.Printf("[ERROR] Error while consuming from %s: %s", inputEventName, err)
			},
		)
	}
}

func isInputComplete(inputVarNames []string, inputs Dictionary) bool {
	for _, inputVarName := range inputVarNames {
		if _, exists := inputs[inputVarName]; !exists {
			return false
		}
	}
	return true
}

func publish(msgBroker msgbroker.CommonBroker, ID string, outputIOList IOList, outputs Dictionary) {
	outputVarNames := outputIOList.GetUniqueVarNames()
	for _, outputVarName := range outputVarNames {
		data := outputs[outputVarName]
		pkg := servicedata.Package{ID: ID, Data: data}
		rawOutputEventNames := outputIOList.GetVarEventNames(outputVarName)
		for _, rawOutputEventName := range rawOutputEventNames {
			eventName := fmt.Sprintf("%s.%s", ID, rawOutputEventName)
			log.Printf("[INFO] Publish into `%s`: `%#v`", eventName, pkg)
			err := msgBroker.Publish(eventName, pkg)
			if err != nil {
				log.Printf("[INFO] Error while publishing into `%s`: `%s`", eventName, err)
			}
		}
	}
}
