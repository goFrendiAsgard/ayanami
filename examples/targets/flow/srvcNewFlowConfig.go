package main

import (
	"fmt"
)

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
