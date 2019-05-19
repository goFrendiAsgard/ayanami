package main

import (
	"encoding/json"
	"fmt"
	nats "github.com/nats-io/nats.go"
	"log"
)

// Call call other service
func Call(serviceType, serviceName, methodName string, inputNames, outputNames []string, inputs Dictionary) (Dictionary, error) {
	var err error
	outputs := make(Dictionary)
	natsURL := GetNatsURL()
	nc, err := nats.Connect(natsURL)
	if err != nil {
		return outputs, err
	}
	ID, err := CreateID()
	if err != nil {
		return outputs, err
	}
	ch := make(chan bool)
	// create consumer for each output
	for _, outputName := range outputNames {
		eventName := fmt.Sprintf("%s.%s.%s.%s.out.%s", ID, serviceType, serviceName, methodName, outputName)
		log.Printf("[INFO] Prepare to consume %s", eventName)
		nc.Subscribe(eventName, func(m *nats.Msg) {
			log.Printf("[INFO] Get message from `%s`: `%s`", eventName, string(m.Data))
			var pkg Pkg
			JSONByte := m.Data
			err = json.Unmarshal(JSONByte, &pkg)
			if err != nil {
				ch <- false
			}
			outputs[outputName] = pkg.Data
			log.Printf("[INFO] Output for %s: %#v", ID, inputs)
			// done
			if isOutputComplete(outputNames, outputs) {
				log.Printf("[INFO] Outputs for %s completed", ID)
				ch <- true
			}
		})
	}
	// creat publisher for each input
	for _, inputName := range inputNames {
		eventName := fmt.Sprintf("%s.%s.%s.%s.in.%s", ID, serviceType, serviceName, methodName, inputName)
		pkg := Pkg{ID: ID, Data: inputs[inputName]}
		JSONByte, err := json.Marshal(&pkg)
		if err != nil {
			ch <- false
		}
		nc.Publish(eventName, JSONByte)
	}
	<-ch
	return outputs, err
}

func isOutputComplete(outputNames []string, outputs Dictionary) bool {
	for _, outputName := range outputNames {
		if _, exists := outputs[outputName]; !exists {
			return false
		}
	}
	return true
}
