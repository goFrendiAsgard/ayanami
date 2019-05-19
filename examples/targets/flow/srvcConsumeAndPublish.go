package main

import (
	"encoding/json"
	"fmt"
	nats "github.com/nats-io/nats.go"
	"log"
)

// ConsumeAndPublish consume from queue and Publish
func ConsumeAndPublish(natsURL string, configs Configs) {
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

func consumeAndPublishSingle(nc *nats.Conn, inputConfig StringDictionary, outputConfig StringDictionary, wrappedFunction WrappedFunction) {
	// allInputs
	allInputs := make(map[string]Dictionary)
	for rawEventName, inputName := range inputConfig {
		eventName := fmt.Sprintf("*.%s", rawEventName)
		log.Printf("[INFO] Consume from `%s`", eventName)
		nc.Subscribe(eventName, func(m *nats.Msg) {
			var pkg Pkg
			JSONByte := m.Data
			err := json.Unmarshal(JSONByte, &pkg)
			if err != nil {
				log.Printf("[ERROR] %s: %s", eventName, err)
				return
			}
			// fill up allInputs
			ID := pkg.ID
			data := pkg.Data
			if _, exists := allInputs[ID]; !exists {
				allInputs[ID] = Dictionary{}
			}
			allInputs[ID][inputName] = data
			inputs := allInputs[ID]
			log.Printf("[INFO] Inputs for %s: %#v", ID, inputs)
			// execute wrapper
			if isInputComplete(inputConfig, inputs) {
				log.Printf("[INFO] Inputs for %s completed", ID)
				outputs := wrappedFunction(inputs)
				publish(nc, ID, outputConfig, outputs)
			}
		})
	}
}

func isInputComplete(inputConfig StringDictionary, inputs Dictionary) bool {
	for _, inputName := range inputConfig {
		if _, exists := inputs[inputName]; !exists {
			return false
		}
	}
	return true
}

func publish(nc *nats.Conn, ID string, outputConfig StringDictionary, outputs Dictionary) {
	for outputName, rawEventName := range outputConfig {
		data := outputs[outputName]
		pkg := Pkg{ID: ID, Data: data}
		publishPkg(nc, ID, rawEventName, pkg)
	}
}

func publishPkg(nc *nats.Conn, ID string, rawEventName string, pkg Pkg) {
	JSONByte, err := json.Marshal(&pkg)
	if err != nil {
		log.Printf("[ERROR] %s: %s", ID, err)
		return
	}
	eventName := fmt.Sprintf("%s.%s", ID, rawEventName)
	nc.Publish(eventName, JSONByte)
	log.Printf("[INFO] Publish into `%s`: `%#v`", eventName, pkg)
}
