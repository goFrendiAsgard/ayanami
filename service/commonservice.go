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
func (services Services) ConsumeAndPublish(broker msgbroker.CommonBroker, serviceName string) {
	for methodName, service := range services {
		consumeAndPublishSingle(broker, serviceName, methodName, service.Input, service.Output, service.ErrorEventName, service.Function)
	}
}

func consumeAndPublishSingle(broker msgbroker.CommonBroker, serviceName, methodName string, inputIOList, outputIOList IOList, rawErrorEventName string, wrappedFunction WrappedFunction) {
	// allInputs
	allInputs := NewDictionaryRW()
	rawInputEventNames := inputIOList.GetUniqueEventNames()
	inputVarNames := inputIOList.GetUniqueVarNames()
	for _, rawInputEventName := range rawInputEventNames {
		inputEventName := fmt.Sprintf("*.%s", rawInputEventName)
		varNames := inputIOList.GetEventVarNames(rawInputEventName)
		log.Printf("[INFO: %s.%s] Consume from `%s`", serviceName, methodName, inputEventName)
		broker.Consume(inputEventName,
			// success callback
			func(pkg servicedata.Package) {
				// prepare allInputs
				ID := pkg.ID
				data := pkg.Data
				if !allInputs.Has(ID) {
					allInputs.Set(ID, NewDictionaryRW())
				}
				// populate allInputs[ID] with varNames and servicedata
				for _, varName := range varNames {
					allInputs.Get(ID).(*DictionaryRW).Set(varName, data)
				}
				// inputs := allInputs.Get(ID).(*DictionaryRW)
				log.Printf("[INFO: %s.%s] Inputs for %s: %#v", serviceName, methodName, ID, allInputs.Get(ID).(*DictionaryRW).GetDictionary())
				// execute wrapper
				if allInputs.Get(ID).(*DictionaryRW).HasAll(inputVarNames) {
					log.Printf("[INFO: %s.%s] Inputs for %s completed", serviceName, methodName, ID)
					outputs, err := wrappedFunction(allInputs.Get(ID).(*DictionaryRW).GetDictionary())
					defer allInputs.Delete(ID)
					if err != nil {
						log.Printf("[ERROR: %s.%s] Error while consuming from %s: %s", serviceName, methodName, inputEventName, err)
						publishError(broker, serviceName, methodName, rawErrorEventName, ID, err)
						return
					}
					log.Printf("[INFO: %s.%s] Outputs of %s are: %#v", serviceName, methodName, ID, outputs)
					publish(broker, serviceName, methodName, rawErrorEventName, ID, outputIOList, outputs)
				}
			},
			// error callback
			func(err error) {
				log.Printf("[ERROR: %s.%s] Error while consuming from %s: %s", serviceName, methodName, inputEventName, err)
				publishError(broker, serviceName, methodName, rawErrorEventName, "No-ID", err)
			},
		)
	}
}

func publish(msgBroker msgbroker.CommonBroker, serviceName, methodName, rawErrorEventName, ID string, outputIOList IOList, outputs Dictionary) error {
	outputVarNames := outputIOList.GetUniqueVarNames()
	for _, outputVarName := range outputVarNames {
		data := outputs.Get(outputVarName)
		pkg := servicedata.Package{ID: ID, Data: data}
		rawOutputEventNames := outputIOList.GetVarEventNames(outputVarName)
		for _, rawOutputEventName := range rawOutputEventNames {
			eventName := fmt.Sprintf("%s.%s", ID, rawOutputEventName)
			log.Printf("[INFO: %s.%s] Publish into `%s`: `%#v`", serviceName, methodName, eventName, pkg)
			err := msgBroker.Publish(eventName, pkg)
			if err != nil {
				log.Printf("[ERROR: %s.%s] Error while publishing into `%s`: `%s`", serviceName, methodName, eventName, err)
				return publishError(msgBroker, serviceName, methodName, rawErrorEventName, ID, err)
			}
		}
	}
	return nil
}

func publishError(msgBroker msgbroker.CommonBroker, serviceName, methodName, rawErrorEventName, ID string, err error) error {
	errorPkg := servicedata.Package{ID: ID, Data: fmt.Sprintf("%s", err)}
	errorEventName := fmt.Sprintf("%s.%s", ID, rawErrorEventName)
	log.Printf("[INFO: %s.%s] Publish error to `%s`: `%s`", serviceName, methodName, errorEventName, err)
	return msgBroker.Publish(errorEventName, errorPkg)
}
