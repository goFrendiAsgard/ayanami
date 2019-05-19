package main

import (
	"fmt"
)

// SingleConfig single configuration
type SingleConfig struct {
	Input    StringDictionary
	Output   StringDictionary
	Function WrappedFunction
}

// Configs configuration
type Configs = map[string]SingleConfig

// NewServiceConfig create new singleConfig for service
func NewServiceConfig(serviceName string, methodName string, wrappedFunction WrappedFunction, inputs []string, outputs []string) SingleConfig {
	// get inputConfig
	inputConfig := make(StringDictionary)
	for _, inputName := range inputs {
		eventName := fmt.Sprintf("srvc.%s.%s.in.%s", serviceName, methodName, inputName)
		inputConfig[eventName] = inputName
	}
	// get outputConfig
	outputConfig := make(StringDictionary)
	for _, outputName := range outputs {
		eventName := fmt.Sprintf("srvc.%s.%s.in.%s", serviceName, methodName, outputName)
		outputConfig[outputName] = eventName
	}
	// return config
	return SingleConfig{
		Input:    inputConfig,
		Output:   outputConfig,
		Function: wrappedFunction,
	}
}
