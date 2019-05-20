package main

import (
	"fmt"
)

// SrvcSingleConfig single configuration
type SrvcSingleConfig struct {
	Input    []SrvcServiceIO
	Output   []SrvcServiceIO
	Function SrvcWrappedFunction
}

// SrvcConfigs configuration
type SrvcConfigs = map[string]SrvcSingleConfig

// SrvcEventFlow event flow
type SrvcEventFlow struct {
	InputEvent  string
	OutputEvent string
}

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

// SrvcNewFlowConfig create new singleConfig for flow
func SrvcNewFlowConfig(flowName string, methodName string, inputs []string, outputs []string, flows []SrvcEventFlow) SrvcSingleConfig {
	// get inputConfig
	var inputConfig []SrvcServiceIO
	for _, inputName := range inputs {
		eventName := fmt.Sprintf("flow.%s.%s.in.%s", flowName, methodName, inputName)
		inputConfig = append(inputConfig, SrvcServiceIO{VarName: inputName, EventName: eventName})
	}
	// get outputConfig
	var outputConfig []SrvcServiceIO
	var outputEventNames []string
	for _, outputName := range outputs {
		eventName := fmt.Sprintf("flow.%s.%s.in.%s", flowName, methodName, outputName)
		outputEventNames = append(outputEventNames, eventName)
		outputConfig = append(outputConfig, SrvcServiceIO{VarName: outputName, EventName: eventName})
	}
	// return config
	wrappedFunction := createFlowWrappedFunction(flows, outputEventNames)
	return SrvcSingleConfig{
		Input:    inputConfig,
		Output:   outputConfig,
		Function: wrappedFunction,
	}
}

func createFlowWrappedFunction(flow []SrvcEventFlow, outputEventNames []string) SrvcWrappedFunction {
	return func(inputs SrvcDictionary) SrvcDictionary {
		outputs := make(SrvcDictionary)
		return outputs
	}
}
