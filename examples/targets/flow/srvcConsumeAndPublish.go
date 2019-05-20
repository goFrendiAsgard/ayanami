package main

import (
	"encoding/json"
	"fmt"
	nats "github.com/nats-io/nats.go"
	"log"
)

// SrvcConsumeAndPublish consume from queue and Publish
func SrvcConsumeAndPublish(configs SrvcConfigs) {
	natsURL := SrvcGetNatsURL()
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Print(err)
		return
	}
	for _, singleConfig := range configs {
		inputConfig := singleConfig.Input
		outputConfig := singleConfig.Output
		wrappedFunction := singleConfig.Function
		consumeAndPublishSingle(nc, inputConfig, outputConfig, wrappedFunction)
	}
}

func consumeAndPublishSingle(nc *nats.Conn, inputConfig []SrvcServiceIO, outputConfig []SrvcServiceIO, wrappedFunction SrvcWrappedFunction) {
	// allInputs
	allInputs := make(map[string]SrvcDictionary)
	rawEventNames := SrvcGetUniqueEventNames(inputConfig)
	inputNames := SrvcGetUniqueVarNames(inputConfig)
	for _, rawEventName := range rawEventNames {
		eventName := fmt.Sprintf("*.%s", rawEventName)
		log.Printf("[INFO] Consume from `%s`", eventName)
		nc.Subscribe(eventName, func(m *nats.Msg) {
			var pkg SrvcPkg
			JSONByte := m.Data
			log.Printf("[INFO] Get message from `%s`: %s", eventName, string(JSONByte))
			err := json.Unmarshal(JSONByte, &pkg)
			if err != nil {
				log.Printf("[ERROR] %s: %s", eventName, err)
				return
			}
			// prepare allInputs
			ID := pkg.ID
			data := pkg.Data
			if _, exists := allInputs[ID]; !exists {
				allInputs[ID] = SrvcDictionary{}
			}
			// populate allInputs[ID] with eventInputNames and data
			eventInputNames := SrvcGetEventVarNames(inputConfig, rawEventName)
			for _, inputName := range eventInputNames {
				allInputs[ID][inputName] = data
			}
			inputs := allInputs[ID]
			log.Printf("[INFO] Inputs for %s: %#v", ID, inputs)
			// execute wrapper
			if isInputComplete(inputNames, inputs) {
				log.Printf("[INFO] Inputs for %s completed", ID)
				outputs := wrappedFunction(inputs)
				publish(nc, ID, outputConfig, outputs)
			}
		})
	}
}

func isInputComplete(inputNames []string, inputs SrvcDictionary) bool {
	for _, inputName := range inputNames {
		if _, exists := inputs[inputName]; !exists {
			return false
		}
	}
	return true
}

func publish(nc *nats.Conn, ID string, outputConfig []SrvcServiceIO, outputs SrvcDictionary) {
	outputNames := SrvcGetUniqueVarNames(outputConfig)
	for _, outputName := range outputNames {
		data := outputs[outputName]
		pkg := SrvcPkg{ID: ID, Data: data}
		rawEventNames := SrvcGetVarEventNames(outputConfig, outputName)
		for _, rawEventName := range rawEventNames {
			publishPkg(nc, ID, rawEventName, pkg)
		}
	}
}

func publishPkg(nc *nats.Conn, ID string, rawEventName string, pkg SrvcPkg) {
	JSONByte, err := json.Marshal(&pkg)
	if err != nil {
		log.Printf("[ERROR] %s: %s", ID, err)
		return
	}
	eventName := fmt.Sprintf("%s.%s", ID, rawEventName)
	nc.Publish(eventName, JSONByte)
	log.Printf("[INFO] Publish into `%s`: `%#v`", eventName, pkg)
}
