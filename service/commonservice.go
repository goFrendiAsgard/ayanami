package service

import (
	"fmt"
	"github.com/state-alchemists/ayanami/msgbroker"
	"github.com/state-alchemists/ayanami/servicedata"
	"log"
	"sync"
)

// WrappedFunction any function wrapped for ConsumeAndPublish
type WrappedFunction = func(inputs Dictionary) (Dictionary, error)

// CommonService single configuration
type CommonService struct {
	ServiceName    string
	MethodName     string
	Input          IOList
	Output         IOList
	ErrorEventName string
	Function       WrappedFunction
}

// Services configuration
type Services []CommonService

// ConsumeAndPublish consume from queue and Publish
func (services Services) ConsumeAndPublish(broker msgbroker.CommonBroker, serviceName string) {
	for _, service := range services {
		consumeAndPublishService(broker, service.ServiceName, service.MethodName, service.Input, service.Output, service.ErrorEventName, service.Function)
	}
}

func consumeAndPublishService(broker msgbroker.CommonBroker, serviceName, methodName string, inputIOList, outputIOList IOList, rawErrorEventName string, wrappedFunction WrappedFunction) {
	var lock sync.RWMutex
	allInputs := make(map[string]Dictionary)
	rawInputEventNames := inputIOList.GetUniqueEventNames()
	inputVarNames := inputIOList.GetUniqueVarNames()
	for _, rawInputEventName := range rawInputEventNames {
		inputEventName := fmt.Sprintf("*.%s", rawInputEventName)
		varNames := inputIOList.GetEventVarNames(rawInputEventName)
		log.Printf("[INFO: %s.%s] Subscribe from `%s`", serviceName, methodName, inputEventName)
		broker.Subscribe(inputEventName,
			// success callback
			createServiceConsumerSuccessHandler(broker, serviceName, methodName, inputEventName, rawErrorEventName, outputIOList, inputVarNames, varNames, wrappedFunction, allInputs, &lock),
			// error callback
			createServiceConsumerErrorHandler(broker, serviceName, methodName, inputEventName, rawErrorEventName),
		)
	}
}

func createServiceConsumerSuccessHandler(broker msgbroker.CommonBroker, serviceName, methodName, inputEventName, rawErrorEventName string,
	outputIOList IOList, inputVarNames, varNames []string, wrappedFunction WrappedFunction, allInputs map[string]Dictionary, lock *sync.RWMutex) func(servicedata.Package) {
	return func(pkg servicedata.Package) {
		// prepare allInputs
		ID := pkg.ID
		data := pkg.Data
		lock.Lock()
		if _, exists := allInputs[ID]; !exists {
			allInputs[ID] = make(Dictionary)
		}
		// populate allInputs[ID] with varNames and servicedata
		for _, varName := range varNames {
			allInputs[ID][varName] = data
		}
		lock.Unlock()
		lock.RLock()
		inputs := allInputs[ID]
		log.Printf("[INFO: %s.%s] Inputs for %s: %#v", serviceName, methodName, ID, inputs)
		inputCompleted := inputs.HasAll(inputVarNames)
		lock.RUnlock()
		// execute wrapper
		if inputCompleted {
			log.Printf("[INFO: %s.%s] Inputs for %s completed", serviceName, methodName, ID)
			lock.RLock()
			outputs, err := wrappedFunction(inputs)
			lock.RUnlock()
			// defer allInputs.Delete(ID)
			if err != nil {
				log.Printf("[ERROR: %s.%s] Error while consuming from %s: %s", serviceName, methodName, inputEventName, err)
				err = publishServiceError(broker, serviceName, methodName, rawErrorEventName, ID, err)
				if err != nil {
					log.Printf("[ERROR: %s.%s] Error while publishing error: %s", serviceName, methodName, err)
				}
				return
			}
			log.Printf("[INFO: %s.%s] Outputs of %s are: %#v", serviceName, methodName, ID, outputs)
			err = publishServiceOutput(broker, serviceName, methodName, rawErrorEventName, ID, outputIOList, outputs)
			if err != nil {
				log.Printf("[ERROR: %s.%s] Error while publishing error: %s", serviceName, methodName, err)
			}
		}
	}
}

func createServiceConsumerErrorHandler(broker msgbroker.CommonBroker, serviceName, methodName, inputEventName, rawErrorEventName string) func(error) {
	return func(err error) {
		log.Printf("[ERROR: %s.%s] Error while consuming from %s: %s", serviceName, methodName, inputEventName, err)
		err = publishServiceError(broker, serviceName, methodName, rawErrorEventName, "No-ID", err)
		if err != nil {
			log.Printf("[ERROR: %s.%s] Error while publishing error: %s", serviceName, methodName, err)
		}
	}
}

func publishServiceOutput(broker msgbroker.CommonBroker, serviceName, methodName, rawErrorEventName, ID string, outputIOList IOList, outputs Dictionary) error {
	outputVarNames := outputIOList.GetUniqueVarNames()
	for _, outputVarName := range outputVarNames {
		// if wrapped function doesn't produce current outputVarName, ignore it
		if !outputs.Has(outputVarName) {
			continue
		}
		data := outputs.Get(outputVarName)
		rawOutputEventNames := outputIOList.GetVarEventNames(outputVarName)
		for _, rawOutputEventName := range rawOutputEventNames {
			eventName := fmt.Sprintf("%s.%s", ID, rawOutputEventName)
			err := Publish(serviceName, methodName, broker, ID, eventName, data)
			if err != nil {
				log.Printf("[ERROR: %s.%s] Error while publishing into `%s`: `%s`", serviceName, methodName, eventName, err)
				return publishServiceError(broker, serviceName, methodName, rawErrorEventName, ID, err)
			}
		}
	}
	return nil
}

func publishServiceError(broker msgbroker.CommonBroker, serviceName, methodName, rawErrorEventName, ID string, err error) error {
	errorMessage := fmt.Sprintf("%s", err)
	errorEventName := fmt.Sprintf("%s.%s", ID, rawErrorEventName)
	return Publish(serviceName, methodName, broker, ID, errorEventName, errorMessage)
}
