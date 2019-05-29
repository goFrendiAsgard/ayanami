package service

import (
	"fmt"
	"github.com/state-alchemists/ayanami/msgbroker"
	"github.com/state-alchemists/ayanami/servicedata"
	"log"
)

// WrappedFunction any function wrapped for ConsumeAndPublish
type WrappedFunction = func(inputs Dictionary) (Dictionary, error)

// CommonService single configuration
type CommonService struct {
	Input          IOList
	Output         IOList
	ErrorEventName string
	Function       WrappedFunction
}

// Services configuration
type Services map[string]CommonService

// ConsumeAndPublish consume from queue and Publish
func (services Services) ConsumeAndPublish(broker msgbroker.CommonBroker) {
	for _, service := range services {
		consumeAndPublishSingle(broker, service.Input, service.Output, service.ErrorEventName, service.Function)
	}
}

func consumeAndPublishSingle(broker msgbroker.CommonBroker, inputIOList, outputIOList IOList, rawErrorEventName string, wrappedFunction WrappedFunction) {
	// allInputs
	allInputs := make(map[string]*Dictionary)
	rawInputEventNames := inputIOList.GetUniqueEventNames()
	inputVarNames := inputIOList.GetUniqueVarNames()
	for _, rawInputEventName := range rawInputEventNames {
		inputEventName := fmt.Sprintf("*.%s", rawInputEventName)
		varNames := inputIOList.GetEventVarNames(rawInputEventName)
		log.Printf("[INFO] Consume from `%s`", inputEventName)
		broker.Consume(inputEventName,
			// success callback
			func(pkg servicedata.Package) {
				// prepare allInputs
				ID := pkg.ID
				data := pkg.Data
				if _, exists := allInputs[ID]; !exists {
					allInputs[ID] = &Dictionary{}
				}
				// populate allInputs[ID] with varNames and servicedata
				for _, varName := range varNames {
					allInputs[ID].Set(varName, data)
				}
				inputs := allInputs[ID]
				log.Printf("[INFO] Inputs for %s: %#v", ID, inputs)
				// execute wrapper
				if inputs.HasAll(inputVarNames) {
					log.Printf("[INFO] Inputs for %s completed", ID)
					outputs, err := wrappedFunction(*inputs)
					if err != nil {
						log.Printf("[ERROR] Error while consuming from %s: %s", inputEventName, err)
						publishError(broker, rawErrorEventName, ID, err)
						return
					}
					publish(broker, rawErrorEventName, ID, outputIOList, outputs)
				}
			},
			// error callback
			func(err error) {
				log.Printf("[ERROR] Error while consuming from %s: %s", inputEventName, err)
				publishError(broker, rawErrorEventName, "No-ID", err)
			},
		)
	}
}

func publish(msgBroker msgbroker.CommonBroker, rawErrorEventName, ID string, outputIOList IOList, outputs Dictionary) error {
	outputVarNames := outputIOList.GetUniqueVarNames()
	for _, outputVarName := range outputVarNames {
		data := outputs.Get(outputVarName)
		pkg := servicedata.Package{ID: ID, Data: data}
		rawOutputEventNames := outputIOList.GetVarEventNames(outputVarName)
		for _, rawOutputEventName := range rawOutputEventNames {
			eventName := fmt.Sprintf("%s.%s", ID, rawOutputEventName)
			log.Printf("[INFO] Publish into `%s`: `%#v`", eventName, pkg)
			err := msgBroker.Publish(eventName, pkg)
			if err != nil {
				log.Printf("[INFO] Error while publishing into `%s`: `%s`", eventName, err)
				return publishError(msgBroker, rawErrorEventName, ID, err)
			}
		}
	}
	return nil
}

func publishError(msgBroker msgbroker.CommonBroker, rawErrorEventName, ID string, err error) error {
	errorPkg := servicedata.Package{ID: ID, Data: err}
	errorEventName := fmt.Sprintf("%s.%s", ID, rawErrorEventName)
	return msgBroker.Publish(errorEventName, errorPkg)
}
