package main

import (
	"fmt"
)

// SrvcNewServiceConfig create new singleConfig for service
func SrvcNewServiceConfig(serviceName string, methodName string, inputs []string, outputs []string, wrappedFunction SrvcWrappedFunction) SrvcSingleConfig {
	// get inputConfig
	var inputConfig []SrvcServiceIO
	for _, inputName := range inputs {
		eventName := fmt.Sprintf("srvc.%s.%s.in.%s", serviceName, methodName, inputName)
		inputConfig = append(inputConfig, SrvcServiceIO{VarName: inputName, EventName: eventName})
	}
	// get outputConfig
	var outputConfig []SrvcServiceIO
	for _, outputName := range outputs {
		eventName := fmt.Sprintf("srvc.%s.%s.in.%s", serviceName, methodName, outputName)
		outputConfig = append(outputConfig, SrvcServiceIO{VarName: outputName, EventName: eventName})
	}
	// return config
	return SrvcSingleConfig{
		Input:    inputConfig,
		Output:   outputConfig,
		Function: wrappedFunction,
	}
}
